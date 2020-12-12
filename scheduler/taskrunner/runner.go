package taskrunner

type Runner struct {
	Controller controlChan
	Error      controlChan
	Data       dataChan
	dataSize   int
	longLived  bool // 是否长期存活
	Dispatcher fn
	Executor   fn
}

func NewRunner(size int, longlived bool, d fn, e fn) *Runner {
	return &Runner{
		Controller: make(chan string, 1), // 带buffer的非阻塞channel
		Error:      make(chan string, 1),
		Data:       make(chan interface{}, size),
		longLived:  longlived,
		dataSize:   size,
		Dispatcher: d,
		Executor:   e,
	}
}

func (r *Runner) startDispatch() {
	defer func() {
		if !r.longLived {
			close(r.Controller)
			close(r.Data)
			close(r.Error)
		}
	}()

	for {
		select {
		case c := <-r.Controller:
			if c == READY_TO_DISPATCH {
				// 生产者
				err := r.Dispatcher(r.Data) // 读取任务，实质是通过VideoClearDispatcher函数从数据库读取
				if err != nil {
					r.Error <- CLOSE
				} else {
					r.Controller <- READY_TO_EXECUTE // 通知执行任务
				}
			}

			if c == READY_TO_EXECUTE {
				// 消费者
				err := r.Executor(r.Data) // 执行任务，实质是通过VideoClearExecutor 删除视频，并从数据库删除
				if err != nil {
					r.Error <- CLOSE
				} else {
					r.Controller <- READY_TO_DISPATCH // 通知继续读取任务
				}
			}
		case e := <-r.Error: // 和r.Controller独立
			if e == CLOSE {
				return
			}
		}
	}
}

func (r *Runner) StartAll() {
	r.Controller <- READY_TO_DISPATCH
	r.startDispatch()
}
