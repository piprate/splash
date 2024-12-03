package main

import (
	"context"
	"log"

	"github.com/piprate/splash"
)

func main() {

	g, err := splash.NewConnectorTestNet()
	if err != nil {
		log.Fatal(err)
		return
	}

	eventsFetcher := g.EventFetcher().
		Last(1000).
		Event("A.0b2a3299cc857e29.TopShot.Withdraw")
	//EventIgnoringFields("A.0b2a3299cc857e29.TopShot.Withdraw", []string{"field1", "field"})

	events, err := eventsFetcher.Run(context.Background())
	if err != nil {
		panic(err)
	}

	log.Printf("%v", events)

	//to send events to a discord eventhook use
	//	message, err := g.NewDiscordWebhook("http://your-webhook-url").SendEventsToWebhook(events)

}
