package main

import (
	"fmt"
	"io/ioutil"
	"ipfs/src"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

// @title Swagger Example API
// @version 1.0
// @description This is a sample server Petstore server.
// @termsOfService http://swagger.io/terms/
// @host localhost
// @BasePath /

func main() {

	if _, err := os.Stat("config/config.ini"); os.IsNotExist(err) {
		fmt.Println("配置文件不存在")
		body, _ := ioutil.ReadFile("config.ini")
		file, err := os.Create("config/config.ini")
		if err != nil {
			log.Fatal(err)
		}
		file.WriteString(string(body))
		file.Close()
	}

	router := gin.Default()
	router.LoadHTMLGlob("templates/**/*")
	router.StaticFS("/public", http.Dir("./public"))
	router.Use(src.CORS())

	router.GET("/", src.Index)
	router.POST("/save", src.SaveUploads)
	router.POST("/del", src.DelUploads)
	router.POST("/saveEdit", src.SaveEdit)
	router.GET("/uploads", src.GetUploads)
	router.GET("/sync", src.SyncRemote)
	router.POST("/addtask", src.SaveTask)

	router.GET("/m3u", src.CreateM3u)

	router.POST("/addgateway", src.AddGateway)
	router.POST("/savegate", src.SaveGateway)

	router.Run(":8080")
}
