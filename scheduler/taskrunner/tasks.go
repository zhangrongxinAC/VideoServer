package taskrunner

import (
	"errors"
	"log"
	"os"
	"path/filepath"
	"sync"
	"video_server/scheduler/dbops"
)

//删除文件
func deleteVideo(vid string) error {
	path, _ := filepath.Abs(VIDEO_PATH + vid)
	log.Println(path)
	err := os.Remove(VIDEO_PATH + vid)
	if err != nil && !os.IsNotExist(err) {
		log.Printf("Deleting video error: %v", err)
		return err
	}
	return nil
}

// 从数据库读取要删除的文件名 任务生产者
func VideoClearDispatcher(dc dataChan) error {
	res, err := dbops.ReadVideoDeletionRecord(3) // 批量读取减少数据库压力
	if err != nil {
		log.Printf("Video clear dispatcher error: %v", err)
	}
	if len(res) == 0 {
		return errors.New("All tasks finished")
	}
	for _, id := range res {
		dc <- id // 把读取出来的放入data channel
	}
	return nil
}

// 任务消费者
func VideoClearExecutor(dc dataChan) error {
	errMap := &sync.Map{} // 把出错存储到map
	var err error
forloop:
	for { // 退出这一个大循环
		select {
		case vid := <-dc:
			// 开一个新协程去删除，异步处理存在对应数据还没有删除，VideoClearDispatcher由从数据库读取出来
			go func(id interface{}) { // go xxx就是开一个协程
				// 删除文件
				if err := deleteVideo(id.(string)); err != nil {
					errMap.Store(id, err)
					return
				}
				// 从数据删除记录
				if err := dbops.DelVideoDeletionRecord(id.(string)); err != nil {
					errMap.Store(id, err)
					return
				}
			}(vid)
		default:
			break forloop //break label跳出循环不再执行for，不管有多少从for
		}
	}
	errMap.Range(func(k, v interface{}) bool {
		err = v.(error)
		if err != nil { // 遍历处理有没有问题
			return false
		}
		return true
	})

	return err
}
