package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/paramusio/go-zetka"
)

var (
	Token = flag.String("token", "", "discord bot token")
)

func main() {
	flag.Parse()

	results := make(chan *zetka.GatewayEvent, 32)
	c, err := zetka.New(*Token)
	if err != nil {
		log.Fatal(err)
	}

	// Spawn a routine to listen for messages.
	go func() {
		err = c.Start(context.Background(), results)
		if err != nil {
			log.Fatal(err)
		}
	}()

	for {
		select {
		case msg := <-results:
			err := ioutil.WriteFile(fmt.Sprintf("dump/%d-%s.json", time.Now().Unix(), msg.Type), msg.Data, 0666)
			if err != nil {
				panic(err)
			}
		}
	}
}
