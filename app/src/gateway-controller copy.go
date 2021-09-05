package src

import (
	"fmt"
	"ipfs/inc"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AddGateway(c *gin.Context) {
	var param struct {
		Url string
	}
	err := c.BindJSON(&param)
	if err != nil {
		fmt.Println("some thing wrong")
	}
	inc.SetConfig("ipfsgateways", param.Url)
	c.String(http.StatusOK, "%s", "更新成功")
}

func SaveGateway(c *gin.Context) {
	var param struct {
		Url string
	}
	err := c.BindJSON(&param)
	if err != nil {
		fmt.Println("some thing wrong")
	}
	inc.SetConfig("ipfsgateway", param.Url)
	c.String(http.StatusOK, "%s", "更新成功")
}
