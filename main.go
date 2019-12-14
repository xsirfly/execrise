package main

import (
	"exercise/conf"
	"exercise/database"
	"exercise/docker"
	"exercise/handler/course"
	"exercise/handler/submission"
	"exercise/log"
	"flag"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"exercise/handler/chapter"
	"exercise/cache"
	"net/http"
	"exercise/push"
)

func main() {
	env := flag.String("env", "dev", "runtime env")

	if err := conf.Init(*env); err != nil {
		logrus.WithError(err).Error("config init failed")
		return
	}

	if err := log.Init(); err != nil {
		logrus.WithError(err).Error("log init failed")
		return
	}

	if err := docker.Init(); err != nil {
		logrus.WithError(err).Error("docker client init failed")
		return
	}

	if err := database.Init(); err != nil {
		logrus.WithError(err).Error("database init failed")
		return
	}

	if err := cache.Init(); err != nil {
		logrus.WithError(err).Error("redis init failed")
		return
	}

	push.Init()

	r := gin.Default()
	submission.RegisterRouter(r)
	course.RegisterRouter(r)
	chapter.RegisterRouter(r)

	r.GET("/ws", func(c *gin.Context) {
		connId := c.Query("connId")
		if connId == "" {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
		push.ServeWs(connId, c.Writer, c.Request)
	})

	if err := r.Run(conf.GetConf().Port); err != nil {
		logrus.WithError(err).Error("gin launch failed")
		return
	}
}
