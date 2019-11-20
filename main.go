package main

import (
	"exercise/docker"
	"exercise/handler/submission"
	"github.com/gin-gonic/gin"
	"flag"
	"exercise/conf"
	"github.com/sirupsen/logrus"
	"exercise/log"
	"exercise/database"
	"exercise/handler/course"
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

	r := gin.Default()
	submission.RegisterRouter(r)
	course.RegisterRouter(r)

	if err := r.Run(conf.GetConf().Port); err != nil {
		logrus.WithError(err).Error("gin launch failed")
		return
	}
}
