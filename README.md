# 一、API部分

## session


1、服务启动从DB拉取session到cache

2、用户登录产生session

3、判断session是否过期


这里session用sync.map保存



### api前端部分

main->middleware -> defs(message, err)->handlers->dbops->response


# 二、stream部分

## limiter（流控机制

1、新建一个ConnLimiter使用chan来控制连接数量

```go
// 当bucker满的情况下
if len(cl.bucket) >= cl.concurrentConn {
   log.Printf("Reached the rate limitation.")
   return false
}
```

2、在主函数中添加一个NewMiddleWareHandle，把limiter校验放到函数中。

```go
type midelWareHandler struct {
   r *httprouter.Router
   l *ConnLimiter
}


func NewMiddleWareHandle(r *httprouter.Router, cc int) http.Handler {
   m := midelWareHandler{}
   m.r = r
   m.l = NewConnLimiter(cc)
   return m
}

func (m  midelWareHandler)ServeHTTP (w http.ResponseWriter, r *http.Request)  {
   // 判断如果超过流控值
   if !m.l.GetConn() {
      sendErrorResponse(w, http.StatusTooManyRequests, "Too Many Requests")
      return
   }

   m.r.ServeHTTP(w, r)
   // 释放token
   defer m.l.ReleaseConn()
}
```

midelWareHandler变成http.Handler需先实现ServeHTTP，这里将判断放到ServeHTTP函数中。

## Handler

返回的头文件中加入视频强制格式

> StreamHandler (读取文件产生到页面

```go
// 加入header视频文件强制提醒
w.Header().Set("Content-type", "video/mp4")
// 传输二进制流 播放视频
	http.ServeContent(w, r, "", time.Now(), video)
```

> uploadHandler（页面上传文件到服务端

```go
// 限定上传文件的大小
r.Body = http.MaxBytesReader(w, r.Body, MAX_UPLOAN_SIZE)
if err := r.ParseMultipartForm(MAX_UPLOAN_SIZE); err != nil{
   log.Printf("Error when try upload file: %v\n", err)
   sendErrorResponse(w, http.StatusBadRequest, "File is too big!")
   return
}
```

> testPageHandler(文件上传页面

# 三、schedule(调度程序

作用
- 处理延时操作
- 处理异步任务

结构
- ReSTful 的 http server(接收任务写道schedule
- Timer(定时器
- 生产者/消费者模型下的task runner



```flow
flow
st=>start: Producer/Dispatcher
op=>operation: channel
ops=>operation: Consumer/Executor
st->op->ops
```

## 1、runner

runner.go

startDispatcher:

- control	channel（信息交换
- data        channel（数据

## 2、task

- 延时删除

1. api-> cideoid -> mysql // api拿到id存到数据库
2. dispatcher -> mysql -> videoid ->  datachannel // 拿到id放到channel
3. executor -> datachannel -> videoid -> delete videos // 删除

## 3、trmain



## 4、Api

1. user -> api -> delete video
2. api -> scheduler -> write
3. timer
4. timer -> runner -> read -> exec -> delete video from folder



# 五、数据库表
```mysql
CREATE DATABASE video_server;
```

## users用户表:

```mysql
CREATE TABLE `video_server`.`users`  (
  `id` int unsigned primary key auto_increment,
  `login_name` varchar(64) unique key,
  `pwd` text
);
```

## video_info视频表:

```mysql
CREATE TABLE `video_server`.`video_info`  (
  `id` varchar(64) NOT NULL,
  `author_id` int(10) NULL,
  `name` text NULL,
  `display_ctime` text NULL,
  `create_time` datetime NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
);
```

## comments评论表:

```mysql
CREATE TABLE `video_server`.`comments`  (
  `id` varchar(64) NOT NULL,
  `video_id` varchar(64) NULL,
  `author_id` int(10) NULL,
  `content` text NULL,
  `time` datetime(0) NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
);
```

## session会话表:

```mysql
CREATE TABLE `video_server`.`sessions`  (
  `session_id` varchar(244)  NOT NULL,
  `TTL` tinytext NULL,
  `login_name` text NULL,
  PRIMARY KEY (`session_id`)
);
```

## video_del_rec待删除视频表:

```mysql
CREATE TABLE `video_server`.`video_del_rec`  (
  `video_id` varchar(64) NOT NULL,
  PRIMARY KEY (`video_id`)
);
```

# 六、编译
(1)在video_server目录执行go mod init
go mod init video_server
(2)分别进入 api、scheduler、streamserver、web进行go build生成对应的执行文件










