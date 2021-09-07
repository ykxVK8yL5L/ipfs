package src

import (
	"ipfs/inc"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Index(c *gin.Context) {
	c.HTML(http.StatusOK, "index/index.html", gin.H{
		"ipfsgateway":  inc.GetConfig("ipfsgateway"),
		"ipfsgateways": inc.GetConfig("ipfsgateways"),
		"subtopic":     inc.GetConfig("subtopic"),
	})
}
