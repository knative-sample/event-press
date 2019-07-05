package main

import (
	"flag"

	"golang.org/x/net/context"

	"fmt"

	"time"

	cloudevents "github.com/cloudevents/sdk-go"
)

var (
	sink string
	num  int
)

func init() {
	flag.StringVar(&sink, "sink", "", "send the event sink.")
	flag.IntVar(&num, "num", 20, "send the event num.")
}
func main() {
	flag.Parse()
	fmt.Println(fmt.Sprintf("sink: %s ", sink))
	fmt.Println(fmt.Sprintf("event num: %v ", num))

	ctx := context.Background()
	t, err := cloudevents.NewHTTPTransport(
		cloudevents.WithTarget(sink),
		cloudevents.WithEncoding(cloudevents.HTTPBinaryV02),
	)
	if err != nil {
		panic("failed to create transport, " + err.Error())
	}

	c, err := cloudevents.NewClient(t)
	if err != nil {
		fmt.Println(fmt.Errorf("unable to create cloudevent client: %s", err.Error()))
		return
	}
	for i := 0; i < num; i++ {
		event := cloudevents.NewEvent(cloudevents.VersionV02)
		event.SetID(fmt.Sprintf("%v", i))
		event.SetType("com.cloudevents.press.sent")
		event.SetSource("http://localhost:8080/")
		event.SetData("{\"message\": \"Hello world!\"}")
		start := time.Now().Unix()
		fmt.Println(start)
		if _, err := c.Send(ctx, event); err != nil {
			fmt.Println(fmt.Errorf("%s, error: %s", event.ID(), err.Error()))
			continue
		}
		fmt.Println("send success: " + event.ID())
		fmt.Println(time.Now().Unix() - start)
	}

}
