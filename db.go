package copydb

import (
	"context"
	"strconv"

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

// WithPool configurates items pool.
func WithPool(pool Pool) DBOpt {
	return func(db *DB) {
		db.pool = pool
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

	for {
		select {
		case msg := <-ch:
			if err := db.apply([]byte(msg.Payload), &buf); err != nil {
				db.monitor.ApplyFailed(err)
			}
		case <-ctx.Done():
			return nil
		}
	}
}

func (db *DB) init() error {
	ids, err := db.loadList()
	if err != nil {
		return errors.Wrap(err, "list failed")
	}
	for _, id := range ids {
		if err := db.loadItem(id); err != nil {
			return errors.Wrapf(err, "%s failed", id)
		}
	}
	return nil
}

func (db *DB) apply(payload []byte, buf *Update) error {
	if err := proto.Unmarshal(payload, buf); err != nil {
		return errors.Wrap(err, "unmarshal failed")
	}

	if err := db.items.item(buf.Id, db.pool).apply(buf); err != nil {
		if errors.Cause(err) != ErrVersionConflict {
			return errors.Wrapf(err, "update failed for %s", buf.Id)
		}
		db.monitor.VersionConflictDetected(err)
		return errors.Wrapf(db.loadItem(buf.Id), "refresh failed for %s", buf.Id)
	}

	return nil
}

func (db *DB) replicate(u *Update) error {
	itemKey := db.keys.item(u.Id)

	var res *redis.Cmd
	if u.Remove {
		res = db.processRemove(itemKey, u)
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

func (db *DB) processRemove(itemKey string, u *Update) *redis.Cmd {
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

func (db *DB) loadList() ([]string, error) {
	return db.r.ZRangeByScore(db.keys.list, redis.ZRangeBy{Min: "-inf", Max: "+inf"}).Result()
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

	db.items.item(id, db.pool).init(ver, data)

	return nil
}
