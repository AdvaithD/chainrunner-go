package util

import (
	"fmt"
	"sync"
)

type CustomQueue struct {
	Queue []string
	lock  sync.RWMutex
}

func (c *CustomQueue) Enqueue(name string) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.Queue = append(c.Queue, name)
}

func (c *CustomQueue) Dequeue() error {
	if len(c.Queue) > 0 {
		c.lock.Lock()
		defer c.lock.Unlock()
		c.Queue = c.Queue[1:]
		return nil
	}
	return fmt.Errorf("Pop Error: Queue is empty")
}

func (c *CustomQueue) Front() (string, error) {
	if len(c.Queue) > 0 {
		c.lock.Lock()
		defer c.lock.Unlock()
		return c.Queue[0], nil
	}
	return "", fmt.Errorf("Peep Error: Queue is empty")
}

func (c *CustomQueue) Size() int {
	return len(c.Queue)
}

func (c *CustomQueue) Empty() bool {
	return len(c.Queue) == 0
}

func (c *CustomQueue) Contains(str string) bool {
	for _, v := range c.Queue {
		if v == str {
			return true
		}
	}
	return false
}

// func main() {
// 	CustomQueue := &CustomQueue{
// 		Queue: make([]string, 0),
// 	}

// 	fmt.Printf("Enqueue: A\n")
// 	CustomQueue.Enqueue("A")
// 	fmt.Printf("Enqueue: B\n")
// 	CustomQueue.Enqueue("B")
// 	fmt.Printf("Len: %d\n", CustomQueue.Size())

// 	for CustomQueue.Size() > 0 {
// 		frontVal, _ := CustomQueue.Front()
// 		fmt.Printf("Front: %s\n", frontVal)
// 		fmt.Printf("Dequeue: %s\n", frontVal)
// 		CustomQueue.Dequeue()
// 	}
// 	fmt.Printf("Len: %d\n", CustomQueue.Size())
// }