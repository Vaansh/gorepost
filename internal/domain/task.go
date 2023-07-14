package domain

import (
	"github.com/Vaansh/gore/internal/database"
	"github.com/Vaansh/gore/internal/model"
	"github.com/Vaansh/gore/internal/platform"
	"github.com/Vaansh/gore/internal/publisher"
	"github.com/Vaansh/gore/internal/subscriber"
	"sync"
)

type Task struct {
	Id         string
	Publishers []publisher.Publisher
	Subscriber subscriber.Subscriber
	Quit       chan struct{}
}

func NewTask(publisherIds []string, sources []platform.Name, subscriberId string, destination platform.Name, metadata model.MetaData, repository database.UserRepository) *Task {
	id := string(destination) + subscriberId
	if len(publisherIds) != len(sources) {
		return nil
	}

	publishers := make([]publisher.Publisher, len(publisherIds))
	for i, id := range publisherIds {
		switch sources[i] {
		case platform.YOUTUBE:
			publishers[i] = publisher.NewYoutubePublisher(id)
		default:
			return nil
		}
	}

	var consumer subscriber.Subscriber = nil
	if destination == platform.INSTAGRAM {
		consumer = subscriber.NewInstagramSubscriber(subscriberId, metadata, repository)
	}

	if consumer == nil {
		return nil
	}

	return &Task{
		Id:         id,
		Publishers: publishers,
		Subscriber: consumer,
		Quit:       make(chan struct{}),
	}
}

func (t *Task) Run(stop chan struct{}) {
	c := make(chan model.Post)

	var wg sync.WaitGroup
	wg.Add(len(t.Publishers))

	for _, p := range t.Publishers {
		go func(publisher publisher.Publisher) {
			defer wg.Done()
			publisher.PublishTo(c, t.Quit)
		}(p)
	}

	// Stop the goroutine execution
	go func() {
		select {
		case <-stop:
			close(t.Quit)
			return
		}
	}()

	go func() {
		t.Subscriber.SubscribeTo(c)
		close(t.Quit)
	}()

	wg.Wait()
}