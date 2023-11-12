package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	cloudevents "github.com/cloudevents/sdk-go/v2"
)

func main() {
	ctx := context.Background()
	p, err := cloudevents.NewHTTP(
		cloudevents.WithPort(8080),
		cloudevents.WithGetHandlerFunc(healthz),
	)
	if err != nil {
		log.Fatal("failed to create protocol: ", err)
	}

	c, err := cloudevents.NewClient(p)
	if err != nil {
		log.Fatal("failed to create client: ", err)
	}

	log.Fatal(c.StartReceiver(ctx, receive))
}

func receive(ctx context.Context, event cloudevents.Event) {
	t := event.Type()
	fmt.Printf("T: %+v\n", t)
	sub := event.Subject()
	fmt.Printf("SUB: %+v\n", sub)
}

func healthz(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "ok")
}
