package main

import (
	"context"
	"fmt"
	"sync"

	"github.com/yeeaiclub/a2a-go/sdk/server/event"
)

type QueueManager struct {
	queues map[string]*event.Queue
	mutex  sync.RWMutex
}

func NewQueueManager() *QueueManager {
	return &QueueManager{
		queues: make(map[string]*event.Queue),
	}
}

func (q *QueueManager) Add(ctx context.Context, taskId string, queue *event.Queue) error {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	
	q.queues[taskId] = queue
	return nil
}

func (q *QueueManager) Get(ctx context.Context, taskId string) (*event.Queue, error) {
	q.mutex.RLock()
	defer q.mutex.RUnlock()
	
	queue, exists := q.queues[taskId]
	if !exists {
		return nil, fmt.Errorf("queue not found for task %s", taskId)
	}
	return queue, nil
}

func (q *QueueManager) Tap(ctx context.Context, taskId string) (*event.Queue, error) {
	q.mutex.RLock()
	defer q.mutex.RUnlock()
	
	queue, exists := q.queues[taskId]
	if !exists {
		return nil, fmt.Errorf("queue not found for task %s", taskId)
	}
	return queue, nil
}

func (q *QueueManager) Close(ctx context.Context, taskId string) error {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	
	delete(q.queues, taskId)
	return nil
}

func (q *QueueManager) CreateOrTap(ctx context.Context, taskId string) (*event.Queue, error) {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	
	queue, exists := q.queues[taskId]
	if !exists {
		queue = event.NewQueue(10)
		q.queues[taskId] = queue
	}
	return queue, nil
}
