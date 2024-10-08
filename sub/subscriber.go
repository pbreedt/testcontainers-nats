package sub

import (
	"log"
	"time"

	"github.com/nats-io/nats.go"
)

func GoSubscribe(natsUrl string, forSecs int, ch chan<- int) {
	rec := Subscribe(natsUrl, forSecs)
	log.Printf("Received %d messages", rec)
	ch <- rec
}

func Subscribe(natsUrl string, forSecs int) (received int) {
	con, err := nats.Connect(natsUrl)
	if err != nil {
		log.Fatal(err)
	}
	defer con.Close()

	received = 0
	sub, err := con.Subscribe("topic", func(msg *nats.Msg) {
		received++
		log.Printf("Received: %s", string(msg.Data))
	})
	if err != nil {
		log.Fatal(err)
	}
	defer sub.Unsubscribe()

	log.Printf("Subscribing for %d seconds", forSecs)

	time.Sleep(time.Duration(forSecs) * time.Second)
	// select {}

	log.Printf("Stoped publishing after %d seconds", forSecs)

	return
}
