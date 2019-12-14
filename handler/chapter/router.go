package chapter

import (
	"github.com/gin-gonic/gin"
	"strconv"
	"net/http"
)

func RegisterRouter(r *gin.Engine) {
	g := r.Group("/chapter")
	g.GET("/:chapter", func(c *gin.Context) {
		chapterIdStr := c.Param("chapter")
		chapterId, err := strconv.Atoi(chapterIdStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"message": "课程不存在",
			})
			return
		}
		c.JSON(GetChapter(int64(chapterId)))
	})
}
