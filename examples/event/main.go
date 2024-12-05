package main

import (
	"context"
	"log"

	"github.com/onflow/flowkit/v2/config"
	"github.com/onflow/flowkit/v2/output"
	"github.com/piprate/splash"
)

func main() {

	stdoutLogger := output.NewStdoutLogger(output.InfoLog)
	g, err := splash.NewNetworkConnector(config.DefaultPaths(), splash.NewFileSystemLoader("examples"), "testnet", stdoutLogger)
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

}
