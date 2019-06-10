package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"time"

	"github.com/go-redis/redis"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/regeda/copydb"
	"github.com/regeda/copydb/examples/drones"
	"github.com/regeda/copydb/indexes/spatial"
)

var (
	redisHost  = flag.String("redis-host", "127.0.0.1:6379", "redis host")
	dbCapacity = flag.Int("cap", 2000, "database capacity")
	dbTTL      = flag.Duration("ttl", 5*time.Minute, "item TTL")
	listenHost = flag.String("listen-host", ":8080", "listen host for incoming events")
)

var handlerDurationSummary = prometheus.NewSummaryVec(prometheus.SummaryOpts{
	Name:       "handler_duration",
	Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001, 1.0: 0.001},
	MaxAge:     time.Minute,
}, []string{"handler"})

func init() {
	prometheus.MustRegister(handlerDurationSummary)
}

var dronesCtx = context.Background()

func main() {
	flag.Parse()

	log := log.New(os.Stdout, "dronesdb: ", log.Lshortfile|log.LstdFlags)

	redis.SetLogger(log)

	redisOpts := redis.ClusterOptions{
		Addrs:       []string{*redisHost},
		DialTimeout: time.Second,
	}

	log.Println("connect redis...")

	cluster := redis.NewClusterClient(&redisOpts)
	if res := cluster.Ping(); res.Err() != nil {
		log.Fatalf("failed to ping the server %s: %v", *redisHost, res.Err())
	}

	db, err := copydb.New(cluster,
		copydb.WithCapacity(*dbCapacity),
		copydb.WithTTL(*dbTTL),
		copydb.WithLogger(log),
		copydb.WithBufferedQueries(1024),
		copydb.WithPubSubOpts(copydb.PubSubOpts{
			ReceiveTimeout: time.Second,
			ChannelSize:    1024,
		}),
		copydb.WithPool(spatial.NewIndex(13, func() copydb.Item {
			return make(copydb.SimpleItem)
		})),
	)
	if err != nil {
		log.Fatalf("failed to create a database: %v", err)
	}

	prometheus.MustRegister(newDBStatsCollector(db, log))

	http.HandleFunc("/drone", observeHandler(droneHandler(db), handlerDurationSummary.WithLabelValues("drone")))
	http.HandleFunc("/search", observeHandler(searchHandler(db), handlerDurationSummary.WithLabelValues("search")))
	http.Handle("/metrics", promhttp.Handler())
	go func() {
		log.Println("run http...")

		log.Fatal(http.ListenAndServe(*listenHost, nil))
	}()

	log.Println("run db...")

	log.Fatal(db.Serve())
}

func observeHandler(h http.Handler, s prometheus.Observer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer summaryObserve(s, time.Now())
		h.ServeHTTP(w, r)
	}
}

func summaryObserve(s prometheus.Observer, now time.Time) {
	s.Observe(time.Since(now).Seconds())
}

func droneHandler(db *copydb.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var in drones.Request
		decoder := json.NewDecoder(r.Body)
		defer func() {
			_ = r.Body.Close()
		}()
		if err := decoder.Decode(&in); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		ctx, cancel := context.WithTimeout(dronesCtx, time.Second)
		defer cancel()

		if err := applyRequest(ctx, db, &in); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func searchHandler(db *copydb.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req spatial.SearchRequest

		// read HTTP query
		query := r.URL.Query()
		_, err := fmt.Sscanf(query.Get("coord"), "%f,%f", &req.Lon, &req.Lat)
		if err != nil {
			http.Error(w, fmt.Sprintf("coord parameter is wrong: %v", err), http.StatusBadRequest)
			return
		}
		_, _ = fmt.Sscan(query.Get("radius"), &req.Radius)
		_, _ = fmt.Sscan(query.Get("limit"), &req.Limit)

		ctx, cancel := context.WithTimeout(dronesCtx, time.Second)
		defer cancel()

		// run CopyDB query
		err = db.Query(ctx, copydb.QueryPool(spatial.Search(
			req,
			func(item copydb.Item) {
				it := item.(copydb.SimpleItem)
				_, _ = fmt.Fprintf(w, "%s - %s\n", it["$id"], it["geom"])
			},
			func(err error) {
			},
		)))
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to query db: %v", err), http.StatusInternalServerError)
			return
		}
	}
}

func applyRequest(ctx context.Context, db *copydb.DB, in *drones.Request) error {
	stmt := copydb.NewStatement(in.ID)
	if in.Remove {
		stmt.Remove()
	} else {
		for name, data := range in.Set {
			stmt.Set(name, data)
		}
		for _, name := range in.Unset {
			stmt.Unset(name)
		}
	}
	return stmt.Exec(ctx, db, in.Currtime)
}

type dbStatsCollector struct {
	db     *copydb.DB
	logger *log.Logger

	itemsApplied            *prometheus.Desc
	itemsFailed             *prometheus.Desc
	itemsEvicted            *prometheus.Desc
	itemsReplicated         *prometheus.Desc
	versionConflictDetected *prometheus.Desc
	dbScanned               *prometheus.Desc
}

func newDBStatsCollector(db *copydb.DB, logger *log.Logger) *dbStatsCollector {
	return &dbStatsCollector{
		db:     db,
		logger: logger,
		itemsApplied: prometheus.NewDesc(
			"copydb_items_applied",
			"count of applied items",
			nil, nil),
		itemsFailed: prometheus.NewDesc(
			"copydb_items_failed",
			"count of failed items",
			nil, nil),
		itemsEvicted: prometheus.NewDesc(
			"copydb_items_evicted",
			"count of evicted items",
			nil, nil),
		itemsReplicated: prometheus.NewDesc(
			"copydb_items_replicated",
			"count of replicated items",
			nil, nil),
		versionConflictDetected: prometheus.NewDesc(
			"copydb_version_conflict_detected",
			"count of version conflict detected",
			nil, nil),
		dbScanned: prometheus.NewDesc(
			"copydb_db_scanned",
			"count of db scanned",
			nil, nil),
	}
}

func (c *dbStatsCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.itemsApplied
	ch <- c.itemsFailed
	ch <- c.itemsEvicted
	ch <- c.itemsReplicated
	ch <- c.versionConflictDetected
	ch <- c.dbScanned
}

func (c *dbStatsCollector) Collect(ch chan<- prometheus.Metric) {
	ctx, cancel := context.WithTimeout(dronesCtx, time.Second)
	defer cancel()

	err := c.db.Query(ctx, copydb.QueryStats(func(stats copydb.Stats) {
		ch <- prometheus.MustNewConstMetric(c.itemsApplied, prometheus.GaugeValue, float64(stats.ItemsApplied))
		ch <- prometheus.MustNewConstMetric(c.itemsFailed, prometheus.GaugeValue, float64(stats.ItemsFailed))
		ch <- prometheus.MustNewConstMetric(c.itemsEvicted, prometheus.GaugeValue, float64(stats.ItemsEvicted))
		ch <- prometheus.MustNewConstMetric(c.itemsReplicated, prometheus.GaugeValue, float64(stats.ItemsReplicated))
		ch <- prometheus.MustNewConstMetric(c.versionConflictDetected, prometheus.GaugeValue, float64(stats.VersionConfictDetected))
		ch <- prometheus.MustNewConstMetric(c.dbScanned, prometheus.GaugeValue, float64(stats.DBScanned))
	}))
	if err != nil {
		c.logger.Printf("failed to query stats: %v", err)
	}
}
