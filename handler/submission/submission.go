package submission

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"path/filepath"
	"io/ioutil"
	"net/http"
)

type submitReq struct {
	Code string `json:"code"`
}
func Submit(c *gin.Context) {
	var req submitReq
	c.BindJSON(&req)
	submissionId := uuid.New()

	path := filepath.Join("/opt/code", submissionId.String(), "Main.java")
	if err := ioutil.WriteFile(path, []byte(req.Code), 0666); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
	}




}
