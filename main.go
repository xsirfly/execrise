package main

import (
	"exercise/docker"
	"exercise/handler/submission"
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	err := docker.Init()
	if err != nil {
		log.Fatal(err)
		return
	}

	r := gin.Default()
	submission.RegisterRouter(r)

	if err := r.Run(); err != nil {
		log.Fatal(err)
	}
}
