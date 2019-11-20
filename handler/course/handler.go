package course

import (
	"github.com/gin-gonic/gin"
	"exercise/database"
	"exercise/utils"
	"net/http"
)

func GetCourse(c *gin.Context) {
	courses, err := database.GetCourses()
	if err != nil {
		utils.AbortRequesrtWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, courses)

}
