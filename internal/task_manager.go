package internal

import (
	"fmt"
	"github.com/Vaansh/gore/internal/platform"
	"github.com/Vaansh/gore/internal/platform/instagram"
	"github.com/Vaansh/gore/internal/platform/youtube"
	"sync"
)

type TaskManager struct {
	Tasks map[string]*Task
}

func NewTaskManager() *TaskManager {
	return &TaskManager{
		Tasks: make(map[string]*Task),
	}
}

func (tm *TaskManager) RunAll() {
	var wg sync.WaitGroup
	for _, task := range tm.Tasks {
		wg.Add(1)
		go func(t *Task) {
			defer wg.Done()
			t.Run()
			select {}
		}(task)
	}
	wg.Wait()
}

func (tm *TaskManager) AddTask(producerIDs []string, sources []PlatformName, consumerID string, destination PlatformName) error {
	taskID := string(destination) + consumerID
	if _, exists := tm.Tasks[taskID]; exists {
		return fmt.Errorf("task with ID %s already exists", taskID)
	}

	if len(producerIDs) != len(sources) {
		return fmt.Errorf("received %d producerIds and %d platforms", len(producerIDs), len(sources))
	}

	prods := make([]platform.Publisher, len(producerIDs))
	for i, id := range producerIDs {
		switch sources[i] {
		case PLATFORM:
			prods[i] = youtube.NewPublisher(id)
		default:
			return fmt.Errorf("platform not found %s for %s", sources[i], id)
		}
	}

	consumer := instagram.NewSubscriber(consumerID)
	task := NewTask(taskID, prods, consumer)
	tm.Tasks[task.ID] = task
	return nil
}

func (tm *TaskManager) EditTask(taskID string, producers []platform.Publisher, consumer platform.Subscriber) error {
	task, ok := tm.Tasks[taskID]
	if !ok {
		return fmt.Errorf("task %s not found", taskID)
	}
	task.Producers = producers
	task.Subscriber = consumer
	return nil
}

func (tm *TaskManager) DeleteTask(taskID string) {
	delete(tm.Tasks, taskID)
}
