package snowflake

import (
	"errors"
	"sync"
	"time"
)

const (
	workerBits  uint8 = 10
	numberBits  uint8 = 12
	workerMax   int64 = -1 ^ (-1 << workerBits)
	numberMax   int64 = -1 ^ (-1 << numberBits)
	timeShift   uint8 = workerBits + numberBits
	workerShift uint8 = numberBits
	startTime   int64 = 1609459200000
)

type Worker struct {
	mu        sync.Mutex
	timeStamp int64
	workerId  int64
	number    int64
}

func NewWorker(workerId int64) (*Worker, error) {
	if workerId < 0 || workerId > workerMax {
		return nil, errors.New("worker id out of range")
	}
	return &Worker{
		timeStamp: 0,
		workerId:  workerId,
		number:    0,
	}, nil
}

func (w *Worker) Id() int64 {
	w.mu.Lock()
	defer w.mu.Unlock()

	now := time.Now().UnixNano() / 1e6
	if w.timeStamp == now {
		w.number++
		if w.number > numberMax {
			for now <= w.timeStamp {
				now = time.Now().UnixNano() / 1e6
			}
		}
	} else {
		w.number = 0
		w.timeStamp = now
	}
	return ((now - startTime) << timeShift) | (w.workerId << workerShift) | w.number
}
