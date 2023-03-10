package main

import (
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/micro"
)

func main() {
	nc, err := nats.Connect(
		nats.DefaultURL, // probably not ideal here, would be config based
		nats.UserCredentials("./nsc/keys/creds/local/APP/user.creds"),
		nats.MaxReconnects(-1),           // reconnect forever
		nats.RetryOnFailedConnect(false), // would set this to true normally in a long running application that wanted to start without nats running yet
		nats.Name("playground"),          // set to the unique instance name of the application
	)

	if err != nil {
		panic(err)
	}

	defer nc.Close()

	svc, err := mkService(nc)
	if err != nil {
		panic(err)
	}
	defer svc.Stop()

	resp, err := nc.Request("playground.v1.echo", []byte("echo this"), time.Second)
	if err != nil {
		fmt.Printf("ERR: %s\n", err.Error())
	} else {
		fmt.Printf("RES: %s\n", string(resp.Data))
	}

	sigCh := make(chan os.Signal)
	signal.Notify(sigCh, os.Interrupt, os.Kill)
	defer signal.Reset()

	<-sigCh
	fmt.Printf("Exiting Application\n")
}

func mkService(nc *nats.Conn) (micro.Service, error) {
	return micro.AddService(nc, micro.Config{
		Name:    "playground",
		Version: "1.0.0",
		Endpoint: &micro.EndpointConfig{
			Subject: "playground.v1.echo",
			Handler: micro.HandlerFunc(func(req micro.Request) {
				req.Respond(req.Data())
			}),
			Schema: &micro.Schema{
				Request:  "string",
				Response: "string",
			},
		},
		APIURL:      "https://mammatus.cloud",
		Description: "playground for things",
		StatsHandler: func(e *micro.Endpoint) interface{} {
			return map[string]interface{}{
				"count": 1,
			}
		},
		DoneHandler: func(s micro.Service) {
			fmt.Printf("Recess is over!\n")
		},
		ErrorHandler: func(s micro.Service, n *micro.NATSError) {
			fmt.Printf("ERR: %v", n.Error())
		},
	})
}
