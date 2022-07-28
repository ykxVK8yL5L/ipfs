package src

import (
	"context"
	"fmt"
	"io/ioutil"
	"ipfs/inc"
	"ipfs/model"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type FUpload struct {
	Name string
	Cid  string
	Size string
}

type Task struct {
	Url   string
	Isnow int
}

func SyncRemote(c *gin.Context) {
	findOptions := options.Find()
	uri := inc.GetConfig("mongdbsrc")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()
	// Ping the primary
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		panic(err)
	}
	collection := client.Database(inc.GetConfig("dbname")).Collection("web3")
	filter := bson.D{{"issync", 0}}
	cur, errf := collection.Find(context.TODO(), filter, findOptions)
	if errf != nil {
		log.Fatal(err)
	}
	// Finding multiple documents returns a cursor
	// Iterating through the cursor allows us to decode documents one at a time
	for cur.Next(context.TODO()) {
		// create a value into which the single document can be decoded
		var elem FUpload
		err := cur.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}
		size, _ := strconv.ParseInt(elem.Size, 10, 64)
		upload := model.Upload{Name: elem.Name, Cid: elem.Cid, Size: size}
		db.Create(&upload) // 通过数据的指针来创建
		//results = append(results,&elem)
	}

	result, err := collection.UpdateMany(context.TODO(), filter, bson.D{{"$set", bson.D{{"issync", "1"}}}})
	if err != nil {
		c.String(http.StatusOK, "%s", "更新失败")
	}
	c.String(http.StatusOK, "同步成功,共更新%d文件", result.ModifiedCount)

}

func SaveTask(c *gin.Context) {
	var param struct {
		Url   string
		Isnow int
	}
	err := c.BindJSON(&param)
	if err != nil {
		fmt.Println("some thing wrong")
	}
	uri := inc.GetConfig("mongdbsrc")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()
	// Ping the primary
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		panic(err)
	}
	collection := client.Database(inc.GetConfig("dbname")).Collection("task")
	task := Task{Url: param.Url, Isnow: param.Isnow}
	collection.InsertOne(context.TODO(), task)

	if param.Isnow == 1 {
		url := inc.GetConfig("workflowurl")
		method := "POST"
		payload := strings.NewReader(`{"ref":"main"}`)
		client := &http.Client{}
		req, err := http.NewRequest(method, url, payload)
		if err != nil {
			fmt.Println(err)
			return
		}
		req.Header.Add("Accept", "application/vnd.github.v3+json")
		req.Header.Add("Authorization", "Bearer "+inc.GetConfig("githubtoken"))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

		res, err := client.Do(req)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer res.Body.Close()

		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			fmt.Println(err)
			return
		}
		c.String(http.StatusOK, "添加成功,并已经成功触发%s", string(body))
	}
	c.String(http.StatusOK, "添加成功")
}
