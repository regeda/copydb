package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	_ "net/http/pprof"
	"time"

	"github.com/go-redis/redis"
	"github.com/mailru/easyjson"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/regeda/copydb"
	"github.com/regeda/copydb/examples/providers"
)

var (
	redisHost  = flag.String("redis-host", "127.0.0.1:6379", "redis host")
	dbCapacity = flag.Int("cap", 2000, "database capacity")
	listenHost = flag.String("listen-host", ":8080", "listen host for incoming events")
)

var defSummaryObjectives = map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001, 1.0: 0.001}

const defMaxAge = time.Minute

var (
	providerHandlerDurationSummary = prometheus.NewSummary(prometheus.SummaryOpts{
		Name:       "provider_handler_duration",
		Objectives: defSummaryObjectives,
		MaxAge:     defMaxAge,
	})
	versionConflictCounter = prometheus.NewCounter(prometheus.CounterOpts{Name: "version_conflict"})
)

func init() {
	prometheus.MustRegister(providerHandlerDurationSummary)
	prometheus.MustRegister(versionConflictCounter)
}

func main() {
	flag.Parse()

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
		copydb.WithMonitor(new(monitor)),
	)
	if err != nil {
		log.Fatalf("failed to create a database: %v", err)
	}

	http.HandleFunc("/provider", observeHandler(providerHandler(db), providerHandlerDurationSummary))
	http.Handle("/metrics", promhttp.Handler())
	go func() {
		log.Println("run http...")

		log.Fatal(http.ListenAndServe(*listenHost, nil))
	}()

	log.Println("run db...")

	log.Fatal(db.Serve(context.Background()))
}

type monitor struct{}

func (m *monitor) VersionConflictDetected(err error) {
	versionConflictCounter.Inc()
	log.Printf("ERROR: %v", err)
}

func (m *monitor) ApplyFailed(err error) {
	log.Printf("ERROR: %v", err)
}

func (m *monitor) PurgeFailed(err error) {
	log.Printf("ERROR: %v", err)
}

func observeHandler(h http.Handler, s prometheus.Summary) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer summaryObserve(s, time.Now())
		h.ServeHTTP(w, r)
	}
}

func summaryObserve(s prometheus.Summary, now time.Time) {
	s.Observe(time.Since(now).Seconds())
}

func providerHandler(db *copydb.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var in providers.Request
		if err := easyjson.UnmarshalFromReader(r.Body, &in); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err := applyRequest(db, &in); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func applyRequest(db *copydb.DB, in *providers.Request) error {
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
	return stmt.Exec(db, in.Currtime)
}
