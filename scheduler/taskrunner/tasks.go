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

// 从数据库读取要删除的文件名
func VideoClearDispatcher(dc dataChan) error {
	res, err := dbops.ReadVideoDeletionRecord(3) // 批量读取减少数据库压力
	if err != nil {
		log.Printf("Video clear dispatcher error: %v", err)
	}
	if len(res) == 0 {
		return errors.New("All tasks finished")
	}
	for _, id := range res {
		dc <- id // 把读取出来的放入channel
	}
	return nil
}

func VideoClearExecutor(dc dataChan) error {
	errMap := &sync.Map{}
	var err error
forloop:
	for {
		select {
		case vid := <-dc:
			// 开一个新协程去删除，异步处理存在对应数据还没有删除，VideoClearDispatcher由从数据库读取出来
			go func(id interface{}) {
				if err := deleteVideo(id.(string)); err != nil {
					errMap.Store(id, err)
					return
				}
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
		if err != nil {
			return false
		}
		return true
	})

	return err
}
