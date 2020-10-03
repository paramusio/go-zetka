package zetka_test

import (
	"context"
	"log"

	"github.com/paramusio/go-zetka"
)

func Example() {
	client, err := zetka.New("foobartoken")
	if err != nil {
		log.Fatal(err)
	}

	results := make(chan *zetka.GatewayEvent)

	go func(results chan *zetka.GatewayEvent) {
		for event := range results {
			log.Println(event.Type)
		}
	}(results)

	if err := client.Start(context.Background(), results); err != nil {
		log.Fatal(err)
	}
}
