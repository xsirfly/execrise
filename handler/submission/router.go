package submission

import (
	"github.com/gin-gonic/gin"
)

func RegisterRouter(r *gin.Engine) {
	g := r.Group("/submission")
	g.POST("/:course/:chapter", func(c *gin.Context) {
		var req submitReq
		c.ShouldBindUri(&req)
		c.ShouldBindJSON(&req)
		code, resp := Submit(&req)
		c.JSON(code, resp)
	})
	g.GET("/:token/status", func(c *gin.Context) {
		token := c.Param("token")
		c.JSON(GetSubmitStatus(token))
	})
}

