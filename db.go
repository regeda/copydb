package copydb

import (
	"container/list"
	"context"
	"strconv"
	"time"

	"github.com/go-redis/redis"
	"github.com/gogo/protobuf/proto"
	"github.com/pkg/errors"
)

// DB implements all subscriptions and updates.
type DB struct {
	r Redis

	items
	keys

	pool Pool

	monitor Monitor

	queries chan Query

	ttl time.Duration
	lru *list.List
}

// DBOpt configures DB.
type DBOpt func(*DB)

// WithItemKeyPattern configures item key pattern.
func WithItemKeyPattern(s string) DBOpt {
	return func(db *DB) {
		db.keys.itemPattern = s
	}
}

// WithListKey configures list key.
func WithListKey(s string) DBOpt {
	return func(db *DB) {
		db.keys.list = s
	}
}

// WithChannelKey configures channel key.
func WithChannelKey(s string) DBOpt {
	return func(db *DB) {
		db.keys.channel = s
	}
}

// WithCapacity configures items capacity.
func WithCapacity(c int) DBOpt {
	return func(db *DB) {
		db.items = make(items, c)
	}
}

// WithMonitor configures a monitor.
func WithMonitor(m Monitor) DBOpt {
	return func(db *DB) {
		db.monitor = m
	}
}

// WithPool configures items pool.
func WithPool(pool Pool) DBOpt {
	return func(db *DB) {
		db.pool = pool
	}
}

// WithBufferedQueries configures a capacity of queries queue.
func WithBufferedQueries(c int) DBOpt {
	return func(db *DB) {
		db.queries = make(chan Query, c)
	}
}

// WithTTL configures a cache TTL.
func WithTTL(d time.Duration) DBOpt {
	return func(db *DB) {
		db.ttl = d
	}
}

// New creates a new DB.
func New(r Redis, opts ...DBOpt) (*DB, error) {
	if err := setupScripts(r); err != nil {
		return nil, err
	}

	db := DB{
		r:       r,
		items:   make(items),
		keys:    defaultKeys,
		pool:    new(defaultPool),
		monitor: defaultMonitor,
		queries: make(chan Query),
		lru:     list.New(),
	}

	for _, opt := range opts {
		opt(&db)
	}

	return &db, nil
}

// Serve subscribes to the channel updates.
func (db *DB) Serve(ctx context.Context) error {
	if err := db.init(); err != nil {
		return errors.Wrap(err, "init failed")
	}

	pubsub := db.r.Subscribe(db.keys.channel)
	if _, err := pubsub.Receive(); err != nil {
		return errors.Wrap(err, "subscribe failed")
	}
	defer func() {
		_ = pubsub.Close()
	}()

	var buf Update
	ch := pubsub.Channel()

	var evictChan <-chan time.Time
	if db.ttl > 0 {
		ticker := time.NewTicker(db.ttl)
		defer ticker.Stop()

		evictChan = ticker.C
	}

	for {
		select {
		case msg := <-ch:
			if err := db.apply([]byte(msg.Payload), &buf); err != nil {
				db.monitor.ApplyFailed(err)
			}
		case query := <-db.queries:
			query.scan(db.items)
		case now := <-evictChan:
			db.evictExpired(now)
		case <-ctx.Done():
			return nil
		}
	}
}

// Queries returns a channel to accept queries.
func (db *DB) Queries() chan<- Query {
	return db.queries
}

func (db *DB) evictExpired(t time.Time) {
	deadline := t.Add(-db.ttl).Unix()
	for db.items.evict(deadline, db.pool, db.lru) {
		// the loop stops when all expired items removed
	}
}

func (db *DB) init() error {
	min := "-inf"
	if db.ttl > 0 {
		min = strconv.FormatInt(time.Now().Add(-db.ttl).Unix(), 10)
	}
	ids, err := db.loadList(redis.ZRangeBy{Min: min, Max: "+inf"})
	if err != nil {
		return errors.Wrap(err, "list failed")
	}
	for _, id := range ids {
		if err := db.loadItem(id); err != nil {
			db.monitor.ApplyFailed(errors.Wrapf(err, "load failed for %s", id))
			continue
		}
	}
	return nil
}

func (db *DB) apply(payload []byte, buf *Update) error {
	if err := proto.Unmarshal(payload, buf); err != nil {
		return errors.Wrap(err, "unmarshal failed")
	}

	item := db.items.item(buf.Id, db.pool, db.lru)

	if err := item.apply(buf); err != nil {
		if errors.Cause(err) != ErrVersionConflict {
			return errors.Wrapf(err, "update failed for %s", buf.Id)
		}
		db.monitor.VersionConflictDetected(err)
		if err := db.loadItem(buf.Id); err != nil {
			return errors.Wrapf(err, "refresh failed for %s", buf.Id)
		}
	}

	return nil
}

func (db *DB) replicate(u *Update) error {
	itemKey := db.keys.item(u.Id)

	var res *redis.Cmd
	if u.Remove {
		res = db.processRemove(itemKey)
	} else {
		res = db.processUpdate(itemKey, u)
	}

	version, err := res.Int64()
	if err != nil {
		return errors.Wrap(err, "script failed")
	}

	if version == 0 {
		return nil
	}

	u.Version = version

	data, err := proto.Marshal(u)
	if err != nil {
		return errors.Wrap(err, "marshal failed")
	}

	pipe := db.r.Pipeline()
	pipe.ZAdd(db.keys.list, redis.Z{
		Score:  float64(u.Unix),
		Member: u.Id,
	})
	pipe.Publish(db.keys.channel, data)
	_, err = pipe.Exec()
	return errors.Wrap(err, "pipeline failed")
}

func (db *DB) processRemove(itemKey string) *redis.Cmd {
	return removeItemScript.Run(db.r, []string{itemKey})
}

func (db *DB) processUpdate(itemKey string, u *Update) *redis.Cmd {
	params := []interface{}{
		len(u.Set),
		len(u.Unset),
	}
	for _, f := range u.Set {
		params = append(params, f.Name, f.Data)
	}
	for _, f := range u.Unset {
		params = append(params, f.Name)
	}

	return updateItemScript.Run(db.r, []string{itemKey}, params...)
}

func (db *DB) loadList(zRangeBy redis.ZRangeBy) ([]string, error) {
	return db.r.ZRangeByScore(db.keys.list, zRangeBy).Result()
}

func (db *DB) loadItem(id string) error {
	itemKey := db.keys.item(id)

	cmd := db.r.HGetAll(itemKey)
	if err := cmd.Err(); err != nil {
		return errors.Wrap(err, "read failed")
	}

	data := cmd.Val()

	rawver, ok := data["__ver"]
	if !ok {
		return ErrVersionNotFound
	}
	ver, err := strconv.ParseInt(rawver, 10, 64)
	if err != nil {
		return errors.Wrap(err, "wrong version format")
	}
	delete(data, "__ver") // version field is treated separately

	db.items.item(id, db.pool, db.lru).init(ver, data)

	return nil
}
