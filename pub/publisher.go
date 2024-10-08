package pub

import (
	"log"
	"time"

	"github.com/nats-io/nats.go"
)

func GoPublish(natsUrl string, forSecs int, ch chan<- int) {
	pub := Publish(natsUrl, forSecs)
	log.Printf("Published %d messages", pub)
	ch <- pub
}

func Publish(natsUrl string, forSecs int) (published int) {
	con, err := nats.Connect(natsUrl)
	if err != nil {
		log.Fatal(err)
	}
	defer con.Close()

	log.Printf("Publishing for %d seconds", forSecs)
	published = 0

	start := time.Now()
	for time.Since(start) < time.Duration(forSecs)*time.Second {
		msg := time.Now().Format("15:04:05")
		if err := con.Publish("topic", []byte(msg)); err != nil {
			log.Fatal(err)
		}

		published++
		log.Printf("Sent: %s", msg)
		time.Sleep(1 * time.Second)
	}

	log.Printf("Stoped publishing after %d seconds", forSecs)
	return
}
