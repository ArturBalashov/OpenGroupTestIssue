package model

import (
	"math/rand"
	"os"
	"sync"
	"time"
)

type File struct {
	Name    string
	Content []byte
}

type Queue struct {
	Push   chan File
	Buffer []File
	Mu     sync.RWMutex
	Wg     sync.WaitGroup
	Status map[string]string
}

func NewQueue() *Queue {
	return &Queue{
		Push:   make(chan File),
		Status: map[string]string{},
	}
}

func (q *Queue) PushVal(filename string, content []byte) {
	q.Mu.RLock()
	q.Status[filename] = "new"
	q.Mu.RUnlock()

	time.Sleep(RandomNumber(1, 10) * time.Second)

	q.Push <- File{
		Name:    filename,
		Content: content,
	}
}

func (q *Queue) Run() {
	for {
		select {

		case file := <-q.Push:
			q.Buffer = append(q.Buffer, file)
			q.Wg.Add(1)
		}

		for len(q.Buffer) > 0 {
			file := q.Buffer[0]
			q.Buffer = q.Buffer[1:]

			go func() {
				defer q.Wg.Done()

				q.WriteFile(file.Name, file.Content)
			}()
		}
	}
}

func (q *Queue) WriteFile(filename string, content []byte) {
	q.Mu.RLock()
	q.Status[filename] = "in_progress"
	q.Mu.RUnlock()

	time.Sleep(RandomNumber(1, 20) * time.Second)

	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		q.Mu.RLock()
		q.Status[filename] = "failed"
		q.Mu.RUnlock()
		return
	}
	defer f.Close()
	_, err = f.Write(content)
	if err != nil {
		q.Mu.RLock()
		q.Status[filename] = "failed"
		q.Mu.RUnlock()
		return
	}

	q.Mu.RLock()
	q.Status[filename] = "done"
	q.Mu.RUnlock()
}

func RandomNumber(min, max int) time.Duration {
	return time.Duration(min + rand.Intn(max-min))
}
