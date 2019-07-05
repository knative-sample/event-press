package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"time"

	cloudevents "github.com/cloudevents/sdk-go"
	"github.com/knative-sample/event-press/pkg/kncloudevents"
)

/*
Example Output:

☁  cloudevents.Event:
Validation: valid
Context Attributes,
  SpecVersion: 0.2
  Type: dev.knative.eventing.samples.heartbeat
  Source: https://github.com/knative/eventing-sources/cmd/heartbeats/#local/demo
  ID: 3d2b5a1f-10ca-437b-a374-9c49e43c02fb
  Time: 2019-03-14T21:21:29.366002Z
  ContentType: application/json
  Extensions:
    the: 42
    beats: true
    heart: yes
Transport Context,
  URI: /
  Host: localhost:8080
  Method: POST
Data,
  {
    "id":162,
    "label":""
  }
*/

const (
	CHANNEL_CACHE = 50
)

var cacheChannel = make(chan cloudevents.Event, CHANNEL_CACHE)

func dispatch(ctx context.Context, event cloudevents.Event) {
	//tctx := cloudevents.HTTPTransportContextFrom(ctx)
	//h, _ := json.Marshal(tctx)
	//fmt.Printf("event: %s, header: %s \n", event.ID(), h)
	if direct {
		fmt.Printf("start: cloudevents.Event: %s \n", event.ID())
		time.Sleep(time.Duration(timewait) * time.Second)
		fmt.Printf("end: cloudevents.Event: %s \n", event.ID())
	} else {
		cacheChannel <- event
	}

}

var (
	concurrency int
	timewait    int
	cache       int
	direct      bool
)

func init() {
	flag.IntVar(&concurrency, "concurrency", 5, "concurrency workers.")
	flag.IntVar(&timewait, "timewait", 20, "time wait.")
	flag.IntVar(&cache, "cache", 5, "cache event.")
	flag.BoolVar(&direct, "direct", false, "direct process.")
}
func main() {
	flag.Parse()
	fmt.Println(fmt.Sprintf("concurrency: %v ", concurrency))
	fmt.Println(fmt.Sprintf("timewait: %v ", timewait))
	fmt.Println(fmt.Sprintf("cache: %v ", cache))

	cacheChannel = make(chan cloudevents.Event, cache)

	c, err := kncloudevents.NewDefaultClient()
	if err != nil {
		log.Fatal("Failed to create client, ", err)
	}
	go process()
	log.Fatal(c.StartReceiver(context.Background(), dispatch))
}

func process() {
	for i := 0; i < concurrency; i++ {
		go worker()
	}
}
func worker() {
	for {
		select {
		case event := <-cacheChannel:
			fmt.Printf("start: cloudevents.Event: %s \n", event.ID())
			time.Sleep(time.Duration(timewait) * time.Second)
			fmt.Printf("end: cloudevents.Event: %s \n", event.ID())
		}
	}
}
