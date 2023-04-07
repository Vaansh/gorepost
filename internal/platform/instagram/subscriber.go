package instagram

import (
	"fmt"
	"time"
)

type Subscriber struct {
	InstagramID string
}

func NewInstagramSubscriber(InstagramID string) *Subscriber {
	return &Subscriber{InstagramID: InstagramID}
}

func (p *Subscriber) SubscribeTo(c <-chan string) {
	fmt.Println("Consuming...")
	for link := range c {
		fmt.Printf("New video uploaded: %s\n", link)
		time.Sleep(4 * time.Second)
	}
}

func (p *Subscriber) GetSubscriberID() string {
	return p.InstagramID
}