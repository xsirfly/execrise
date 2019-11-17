package submission

import "github.com/gin-gonic/gin"

func registerRouter(r *gin.Engine) {
	r.POST("/submission", Submit)
}
