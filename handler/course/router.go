package course

import "github.com/gin-gonic/gin"

func RegisterRouter(r *gin.Engine) {
	rg := r.Group("/course")
	rg.GET("", GetCourse)
}
