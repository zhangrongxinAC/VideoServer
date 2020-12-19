package main

import (
	"log"
	"net/http"
	"video_server/scheduler/taskrunner"

	"github.com/julienschmidt/httprouter"
)

func RegisterHandlers() *httprouter.Router {
	router := httprouter.New()
	router.GET("/video-delete-record/:vid-id", vidDelRecHandler)
	return router
}

func main() {
	go taskrunner.Start()
	r := RegisterHandlers()
	err := http.ListenAndServe(":10001", r)
	if err != nil {
		log.Println(err)
	}
	log.Println("main fnish")
}
