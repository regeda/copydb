package copydb

import (
	"container/list"
	"context"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/go-redis/redis"
	"github.com/gogo/protobuf/proto"
	"github.com/pkg/errors"
	"github.com/regeda/copydb/internal/model"
)

// DB implements all subscriptions and updates.
type DB struct {
	r Redis

	pubsubOpts PubSubOpts

	items
	keys

	pool Pool

	replica chan model.Item
	queries chan Query

	ttl time.Duration
	lru *list.List

	stopCh chan struct{}

	stats  Stats
	logger *log.Logger
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
		// @NOTE: maybe configure replica capacity of the channel separately
		db.replica = make(chan model.Item, c)
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

// WithLogger configures a logger.
func WithLogger(logger *log.Logger) DBOpt {
	return func(db *DB) {
		db.logger = logger
	}
}

// PubSubOpts contains parameters for pubcub configuration.
type PubSubOpts struct {
	ReceiveTimeout time.Duration
	ChannelSize    int
}

// WithPubSubOpts configures pubsub.
func WithPubSubOpts(opts PubSubOpts) DBOpt {
	return func(db *DB) {
		db.pubsubOpts = opts
	}
}

// New creates a new DB.
func New(r Redis, opts ...DBOpt) (*DB, error) {
	if err := setupScripts(r); err != nil {
		return nil, err
	}

	db := DB{
		r:     r,
		items: make(items),
		keys:  defaultKeys,
		pool: &SimplePool{
			New: func() Item {
				return make(SimpleItem)
			},
		},
		replica: make(chan model.Item),
		queries: make(chan Query),
		lru:     list.New(),
		stopCh:  make(chan struct{}),
		logger:  log.New(os.Stderr, "copydb: ", log.LstdFlags|log.Lshortfile),
	}

	for _, opt := range opts {
		opt(&db)
	}

	if err := db.init(); err != nil {
		return nil, errors.Wrap(err, "init failed")
	}

	// @NOTE: make configurable
	go db.replicator(time.Second)

	return &db, nil
}

// MustNew creates a new DB. It panics if db setup failed.
func MustNew(r Redis, opts ...DBOpt) *DB {
	db, err := New(r, opts...)
	if err != nil {
		panic(err.Error())
	}
	return db
}

// Serve subscribes to the channel updates.
func (db *DB) Serve() error {
	pubsub := db.r.Subscribe(db.keys.channel)
	if _, err := pubsub.ReceiveTimeout(db.pubsubOpts.ReceiveTimeout); err != nil {
		return errors.Wrap(err, "subscribe failed")
	}
	defer func() {
		_ = pubsub.Close()
	}()

	ch := pubsub.ChannelSize(db.pubsubOpts.ChannelSize)

	var evictChan <-chan time.Time
	if db.ttl > 0 {
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()

		evictChan = ticker.C
	}

	var update model.Update

	for {
		select {
		case msg := <-ch:
			if err := proto.Unmarshal([]byte(msg.Payload), &update); err != nil {
				db.logger.Printf("unmarshal failed: %v", err)
				continue
			}
			for _, item := range update.Items {
				if err := db.apply(item); err != nil {
					db.logger.Println(err)
					db.stats.ItemsFailed++
				} else {
					db.stats.ItemsApplied++
				}
			}
		case query := <-db.queries:
			query.scan(db)
			db.stats.DBScanned++
		case now := <-evictChan:
			db.evictExpired(now)
		case <-db.stopCh:
			// @NOTE: drain replica channel
			return nil
		}
	}
}

func (db *DB) replicator(interval time.Duration) {
	var update model.Update

	ticker := time.NewTicker(interval)

	for {
		select {
		case item := <-db.replica:
			update.Items = append(update.Items, item)
		case <-ticker.C:
			if len(update.Items) == 0 {
				// nothing to update
				continue
			}
			if err := db.replicate(&update); err != nil {
				db.logger.Printf("replication failed: %v", err)
			} else {
				db.stats.ItemsReplicated += len(update.Items)
				update.Items = update.Items[:0]
			}
		}
	}
}

func (db *DB) replicate(update *model.Update) error {
	data, err := proto.Marshal(update)
	if err != nil {
		return errors.Wrap(err, "marshal failed")
	}

	// @NOTE; make memory pool for "zz"
	zz := make([]redis.Z, len(update.Items))
	for i, item := range update.Items {
		zz[i] = redis.Z{
			Score:  float64(item.Unix),
			Member: item.ID,
		}
	}

	pipe := db.r.Pipeline()
	pipe.ZAdd(db.keys.list, zz...)
	pipe.Publish(db.keys.channel, data)
	_, err = pipe.Exec()
	return errors.Wrap(err, "pipeline failed")
}

// MustServe starts the database. It panics if serve failed.
func (db *DB) MustServe() {
	if err := db.Serve(); err != nil {
		panic(err.Error())
	}
}

// Stop terminates serve.
func (db *DB) Stop(ctx context.Context) error {
	select {
	case db.stopCh <- struct{}{}:
	case <-ctx.Done():
		return ctx.Err()
	}
	return nil
}

// MustStop terminates serve. It panics if stopping failed.
func (db *DB) MustStop(ctx context.Context) {
	if err := db.Stop(ctx); err != nil {
		panic(err.Error())
	}
}

// Query runs a query and waits for scan completion.
func (db *DB) Query(ctx context.Context, q Query) error {
	done := make(chan struct{})

	select {
	case db.queries <- queryFunc(func(db *DB) {
		defer close(done)
		q.scan(db)
	}):
	case <-ctx.Done():
		return ctx.Err()
	}

	select {
	case <-done:
	case <-ctx.Done():
		return ctx.Err()
	}
	return nil
}

// QueriesIn returns a channel to accept queries.
func (db *DB) QueriesIn() chan<- Query {
	return db.queries
}

func (db *DB) evictExpired(t time.Time) {
	deadline := t.Add(-db.ttl).Unix()
	// the loop stops when all expired items removed
	for db.items.evict(deadline, db.pool, db.lru) {
		db.stats.ItemsEvicted++
	}
}

func (db *DB) init() error {
	min := "-inf"
	if db.ttl > 0 {
		min = strconv.FormatInt(time.Now().Add(-db.ttl).Unix(), 10)
	}
	zz, err := db.loadList(redis.ZRangeBy{Min: min, Max: "+inf"})
	if err != nil {
		return errors.Wrap(err, "list failed")
	}
	for _, z := range zz {
		id := z.Member.(string)
		unix := int64(z.Score)
		if err := db.loadItem(id, unix); err != nil {
			db.logger.Println(errors.Wrapf(err, "load failed for %s", id))
			db.stats.ItemsFailed++
		} else {
			db.stats.ItemsApplied++
		}
	}
	return nil
}

func (db *DB) apply(it model.Item) error {
	item := db.items.item(it.ID, db.pool, db.lru)

	err := item.apply(it)
	if err == nil {
		return nil
	}
	if errors.Cause(err) != ErrVersionConflict {
		return errors.Wrapf(err, "update failed for %s", it.ID)
	}

	db.logger.Println(err)
	db.stats.VersionConfictDetected++

	return errors.Wrapf(db.loadItem(it.ID, it.Unix), "refresh failed for %s", it.ID)
}

func (db *DB) exec(ctx context.Context, item model.Item) error {
	itemKey := db.keys.item(item.ID)

	var res *redis.Cmd
	if item.Remove {
		res = db.processRemove(itemKey)
	} else {
		res = db.processUpdate(itemKey, item)
	}

	version, err := res.Int64()
	if err != nil {
		return errors.Wrap(err, "script failed")
	}

	if version == 0 {
		return ErrZeroVersion
	}

	item.Version = version

	select {
	case db.replica <- item:
	case <-ctx.Done():
		return ctx.Err()
	}

	return nil
}

func (db *DB) processRemove(itemKey string) *redis.Cmd {
	return removeItemScript.Run(db.r, []string{itemKey})
}

func (db *DB) processUpdate(itemKey string, item model.Item) *redis.Cmd {
	// @NOTE: make memory pool for "params"
	params := []interface{}{
		len(item.Set),
		len(item.Unset),
	}
	for _, f := range item.Set {
		params = append(params, f.Name, f.Data)
	}
	for _, f := range item.Unset {
		params = append(params, f.Name)
	}

	return updateItemScript.Run(db.r, []string{itemKey}, params...)
}

func (db *DB) loadList(zRangeBy redis.ZRangeBy) ([]redis.Z, error) {
	return db.r.ZRangeByScoreWithScores(db.keys.list, zRangeBy).Result()
}

func (db *DB) loadItem(id string, unix int64) error {
	itemKey := db.keys.item(id)

	cmd := db.r.HGetAll(itemKey)
	if err := cmd.Err(); err != nil {
		return errors.Wrap(err, "read failed")
	}

	data := cmd.Val()

	rawver, ok := data[keyVer]
	if !ok {
		return ErrVersionNotFound
	}
	ver, err := strconv.ParseInt(rawver, 10, 64)
	if err != nil {
		return errors.Wrap(err, "wrong version format")
	}
	delete(data, keyVer) // version field is treated separately

	db.items.item(id, db.pool, db.lru).init(id, unix, ver, data)

	return nil
}
