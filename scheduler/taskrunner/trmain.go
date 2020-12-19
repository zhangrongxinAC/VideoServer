package taskrunner

import (
	"fmt"
	"time"
)

type Worker struct {
	ticker *time.Ticker
	runner *Runner
}

func NewWorker(interval time.Duration, r *Runner) *Worker {
	return &Worker{
		ticker: time.NewTicker(interval * time.Second),
		runner: r,
	}
}

func (w *Worker) startWorker() {
	for {
		select {
		case <-w.ticker.C: // 时间到了
			fmt.Println("startWorker")
			go w.runner.StartAll() // 每次定时时间到了就触发生产者消费模型
		}
	}
}

func Start() {
	fmt.Println("NewWorker Start")
	r := NewRunner(3, true, VideoClearDispatcher, VideoClearExecutor)
	w := NewWorker(3, r)
	go w.startWorker()
}
