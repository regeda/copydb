package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/google/uuid"
	geojson "github.com/paulmach/go.geojson"
	"github.com/regeda/copydb/examples/drones"
)

var (
	dronesCount = flag.Int("n", 1000, "drones count")
	updateSleep = flag.Duration("sleep", time.Millisecond, "time to sleep between update")
	droneURL    = flag.String("drone-host", "http://localhost:8080/drone", "drone http host")
)

func main() {
	flag.Parse()

	for req := range requestsGenerator(*dronesCount, *updateSleep) {
		data, _ := json.Marshal(req)
		resp, err := http.Post(*droneURL, "application/json", bytes.NewReader(data))
		if err != nil {
			log.Printf("request failed: %v", err)
			continue
		}
		if err := printResponse(resp); err != nil {
			log.Printf("read response failed: %v", err)
		}
	}
}

func printResponse(resp *http.Response) error {
	defer func() {
		_ = resp.Body.Close()
	}()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if len(data) > 0 {
		log.Println(string(data))
	}
	return nil
}

type requestUpdater func(*drones.Request)

func requestsGenerator(n int, sleep time.Duration) <-chan *drones.Request {
	ids := make([]string, n)
	for i := 0; i < n; i++ {
		ids[i] = uuid.New().String()
	}
	requestUpdaters := []requestUpdater{
		coordRandomizer(-73, 40), // New York
		statusRandomizer("x", "y", "z"),
	}
	outCh := make(chan *drones.Request, n)
	go func() {
		for {
			req := drones.Request{
				ID:       ids[rand.Intn(n)],
				Currtime: time.Now(),
				Set:      make(map[string][]byte),
			}

			fi := rand.Intn(len(requestUpdaters))
			requestUpdaters[fi](&req)

			outCh <- &req
			time.Sleep(sleep)
		}
	}()
	return outCh
}

func coordRandomizer(x, y float64) requestUpdater {
	geom := geojson.Geometry{
		Type: geojson.GeometryPoint,
	}
	return func(r *drones.Request) {
		geom.Point = []float64{x + rand.Float64(), y + rand.Float64()}

		r.Set["geom"], _ = geom.MarshalJSON()
	}
}

func statusRandomizer(statuses ...string) requestUpdater {
	data := make([][]byte, len(statuses))
	for i, s := range statuses {
		data[i] = []byte(s)
	}
	return func(r *drones.Request) {
		i := rand.Intn(len(data))
		r.Set["status"] = data[i]
	}
}
