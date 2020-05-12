package concurrent

import (
	"errors"
	"time"
)

type Operation struct {
	ch  chan dataExchange
	err error
}

type dataExchange struct {
	x      int
	result chan int
}

func New() *Operation {
	op := &Operation{
		ch: make(chan dataExchange),
	}

	for i := 0; i < 10; i++ {
		go op.worker()
	}

	return op
}

func (op *Operation) worker() {
	for data := range op.ch {
		y := data.x * 10

		if y > 100 {
			time.Sleep(100 * time.Millisecond)
			op.err = errors.New("can't be larger then 100")
		}

		data.result <- y
	}
}

func (op *Operation) Do(value int) (int, error) {
	ech := dataExchange{x: value, result: make(chan int)}
	op.ch <- ech
	result := <-ech.result

	return result, op.err
}
