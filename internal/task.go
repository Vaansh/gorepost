package internal

import (
	"github.com/Vaansh/gore/internal/model"
	"github.com/Vaansh/gore/internal/publisher"
	"github.com/Vaansh/gore/internal/subscriber"
)

type Task struct {
	ID         string
	Producers  []publisher.Publisher
	Subscriber subscriber.Subscriber
}

func NewTask(Id string, producers []publisher.Publisher, subscriber subscriber.Subscriber) *Task {
	return &Task{
		ID:         Id,
		Producers:  producers,
		Subscriber: subscriber,
	}
}

func (t *Task) Run() {
	c := make(chan model.Post)
	for _, p := range t.Producers {
		go p.PublishTo(c)
	}
	go t.Subscriber.SubscribeTo(c)
}