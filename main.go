package main

import (
	"fmt"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nkeys"
)

var (
	// these would idealy come from vault
	seed = "SUAFYO7YVH22YB2TPJSV2AGRW5ODLR7TWGTFSS2MH7OLXTQAMOGVIBCENY"
	jwt  = "eyJ0eXAiOiJKV1QiLCJhbGciOiJlZDI1NTE5LW5rZXkifQ.eyJqdGkiOiI2NzZBVFlCNFI3VUpXWlNIVDNIRE0yT003UFlGUEdCVTdGTUpXUVRKVEtET1VQNjJYVVRRIiwiaWF0IjoxNjYyMjI3NTM5LCJpc3MiOiJBRFVUSkRDVjNPRVZTQkFZM0FVRTdKVUxRVTVZNlpMUkZLQUFTWVlJRk9VWVpPM0pIV1ZGU1daRiIsIm5hbWUiOiJ1c2VyIiwic3ViIjoiVUFNWkxaR0sySEs1NEpZUVdPVVBDRDM0REdTQ0Q0T0hMRVFEVlZUUTdKTkFERk5HVkE2VlFWR0UiLCJuYXRzIjp7InB1YiI6e30sInN1YiI6e30sInN1YnMiOi0xLCJkYXRhIjotMSwicGF5bG9hZCI6LTEsImlzc3Vlcl9hY2NvdW50IjoiQUE0R0M0NTUySVlRNE1HVE1OUUhXM0tUVEVOTjZDWE5HTFRIWUhOT1RXS05ZWk5ZSE9aS1lSWk0iLCJ0eXBlIjoidXNlciIsInZlcnNpb24iOjJ9fQ.SExq9rSLLvfrlj-N8S8zQTWWAxb5Z8nj2rwX2dDqLjB9pf9_uHEBJNX9sOYsdJONtkOCq412DEtZJCIUKItmDQ"
)

func main() {
	nc, err := nats.Connect(
		nats.DefaultURL, // probably not ideal here, would be config based
		nats.UserJWT(func() (string, error) { return jwt, nil },
			func(nonce []byte) ([]byte, error) {
				kp, err := nkeys.FromSeed([]byte(seed))
				if err != nil {
					return nil, err
				}

				return kp.Sign(nonce)
			}),

		nats.MaxReconnects(-1),           // reconnect forever
		nats.RetryOnFailedConnect(false), // would set this to true normally in a long running application that wanted to start without nats running yet
		nats.Name("playground"),          // set to the unique instance name of the application
	)

	if err != nil {
		panic(err)
	}

	defer nc.Close()

	sub, err := nc.Subscribe("playground.v1.echo", func(msg *nats.Msg) {
		fmt.Printf("REQ: %s; %s\n", msg.Subject, string(msg.Data))
		_ = msg.Respond(msg.Data)
	})

	if err != nil {
		panic(err)
	}
	defer func() { _ = sub.Drain() }()

	resp, err := nc.Request("playground.v1.echo", []byte("echo this"), time.Second)
	if err != nil {
		panic(err)
	}

	fmt.Printf("RES: %s\n", string(resp.Data))
}
