package course

import (
	"github.com/gin-gonic/gin"
	"strconv"
	"net/http"
)

func RegisterRouter(r *gin.Engine) {
	g := r.Group("/course")
	g.GET("", func(c *gin.Context) {
		code, resp := GetCourses()
		c.JSON(code, resp)
	})

	g.GET("/:course/chapters", func(c *gin.Context) {
		courseIdStr := c.Param("course")
		courseId, err := strconv.Atoi(courseIdStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"message": "课程不存在",
			})
			return
		}
		c.JSON(GetChaptersByCourse(int64(courseId)))
	})
}
