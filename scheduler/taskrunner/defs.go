package taskrunner

const (
	READY_TO_DISPATCH = "d" // 开始生产数据(任务)
	READY_TO_EXECUTE  = "e" // 开始消费数据(任务)
	CLOSE             = "c" // 结束任务

	VIDEO_PATH = "./videos/"
)

type controlChan chan string

type dataChan chan interface{}

type fn func(dc dataChan) error // 重新定义类型 fn -> func(dc dataChan) error
