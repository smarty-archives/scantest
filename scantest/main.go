package main

import (
	"log"
	"time"

	"github.com/smartystreets/scantest/scantest/contract"
)

func main() {
	config := parseConfiguration()
	handlers := buildHandlers(config)

	for {
		context := &contract.Context{}

		for x := 0; x < len(handlers) && context.Error == nil; x++ {
			handlers[x].Handle(context)
		}

		if context.Error != contract.ContextComplete {
			log.Fatal(context.Error)
		}

		// TODO: select, looking for input command to run all tests (reset checksummer and continue), else sleep
		time.Sleep(time.Millisecond * 250)
	}
}
