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

const (
	defaultTTL = 4 * time.Hour
)

// DB implements all subscriptions and updates.
type DB struct {
	r Redis

	items
	keys

	ttl time.Duration
	lru *list.List

	pool Pool

	monitor Monitor
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

// WithTTL configures item's ttl.
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
		keys:    defaultKeys,
		ttl:     defaultTTL,
		lru:     list.New(),
		items:   make(items),
		monitor: defaultMonitor,
		pool:    new(defaultPool),
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

	purgeTicker := time.NewTicker(db.ttl)
	defer purgeTicker.Stop()

	for {
		select {
		case msg := <-ch:
			if err := db.apply([]byte(msg.Payload), &buf); err != nil {
				db.monitor.ApplyFailed(err)
			}
		case now := <-purgeTicker.C:
			db.purgeExpired(now)
		case <-ctx.Done():
			return nil
		}
	}
}

func (db *DB) init() error {
	ids, err := db.loadList(redis.ZRangeBy{
		Min: strconv.FormatInt(time.Now().Add(-db.ttl).Unix(), 10),
		Max: "+inf",
	})
	if err != nil {
		return errors.Wrap(err, "list failed")
	}
	for _, id := range ids {
		item, err := db.loadItem(id)
		if err != nil {
			db.monitor.ApplyFailed(errors.Wrapf(err, "load failed for %s", id))
			continue
		}
		db.lru.MoveToBack(item.elem)
	}
	return nil
}

func (db *DB) purgeExpired(t time.Time) {
	item := db.lru.Front()
	deadline := t.Add(-db.ttl).Unix()
	for item != nil {
		id := item.Value.(string)
		purged, err := db.processPurge(id, deadline)
		if err != nil {
			db.monitor.PurgeFailed(errors.Wrapf(err, "purge failed for %s", id))
			continue
		}
		if !purged {
			return
		}
		next := item.Next()
		db.lru.Remove(item)
		db.items.destroy(id, db.pool)
		item = next
	}
}

func (db *DB) processPurge(id string, deadline int64) (bool, error) {
	res := purgeItemScript.Run(db.r, []string{
		db.keys.list,
		db.keys.item(id),
	}, id, deadline)
	if err := res.Err(); err != nil {
		return false, errors.Wrap(err, "script failed")
	}
	affected, err := res.Int64()
	if err != nil {
		return false, errors.Wrap(err, "malformed response")
	}
	return affected == 1, nil
}

func (db *DB) apply(payload []byte, buf *Update) error {
	if err := proto.Unmarshal(payload, buf); err != nil {
		return errors.Wrap(err, "unmarshal failed")
	}

	item := db.items.item(buf.Id, db.pool)

	if err := item.apply(buf); err != nil {
		if errors.Cause(err) != ErrVersionConflict {
			return errors.Wrapf(err, "update failed for %s", buf.Id)
		}
		db.monitor.VersionConflictDetected(err)
		item, err = db.loadItem(buf.Id)
		if err != nil {
			return errors.Wrapf(err, "refresh failed for %s", buf.Id)
		}
	}

	db.lru.MoveToBack(item.elem)

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

	if err := res.Err(); err != nil {
		return errors.Wrap(err, "script failed")
	}

	version, err := res.Int64()
	if err != nil {
		return errors.Wrap(err, "read version failed")
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

func (db *DB) loadItem(id string) (*item, error) {
	itemKey := db.keys.item(id)

	cmd := db.r.HGetAll(itemKey)
	if err := cmd.Err(); err != nil {
		return nil, errors.Wrap(err, "read failed")
	}

	data := cmd.Val()

	rawver, ok := data["__ver"]
	if !ok {
		return nil, ErrVersionNotFound
	}
	ver, err := strconv.ParseInt(rawver, 10, 64)
	if err != nil {
		return nil, errors.Wrap(err, "wrong version format")
	}
	delete(data, "__ver") // version field is treated separately

	item := db.items.item(id, db.pool)
	item.init(ver, data)

	return item, nil
}
