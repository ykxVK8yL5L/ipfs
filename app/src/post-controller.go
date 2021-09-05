package src

import (
	"encoding/json"
	"fmt"
	"ipfs/inc"
	"ipfs/model"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type DatatableResult struct {
	Total    int64               `json:"recordsTotal"`
	Filtered int64               `json:"recordsFiltered"`
	Data     []map[string]string `json:"data"`
}

func GetUploads(c *gin.Context) {
	db = inc.GetDB()

	table := "uploads"
	var total, filtered int64
	var posts []model.Upload
	var items []map[string]string

	query := db.Table(table)
	query = query.Offset(QueryOffset(c))
	query = query.Limit(QueryLimit(c))
	query = query.Order(QueryOrder(c))
	query = query.Scopes(SearchScope(c), DateTimeScope(c))

	if err := query.Find(&posts).Error; err != nil {
		c.AbortWithStatus(404)
		log.Println(err)
		return
	}

	// Filtered data count
	query = query.Offset(0)
	query.Table(table).Count(&filtered)

	// Total data count
	db.Table(table).Count(&total)

	for _, element := range posts {
		iteminfo := make(map[string]string)
		b, _ := json.Marshal(element)
		iteminfo["info"] = (string(b))
		//b, _ := json.Marshal(iteminfo)
		items = append(items, iteminfo)
	}

	result := DatatableResult{
		Total:    total,
		Filtered: filtered,
		Data:     items,
	}
	c.JSON(200, result)
}
func SaveUploads(c *gin.Context) {

	var upload model.Upload
	err := c.BindJSON(&upload)
	if err != nil {
		fmt.Println("some thing wrong")
	}
	result := db.Create(&upload)

	if result.Error != nil {
		c.String(http.StatusOK, "%s", "数据插入失败")
		panic("数据插入失败")
	}
	c.String(http.StatusOK, "%s", "更新成功")

}

func SaveEdit(c *gin.Context) {
	var param struct {
		Id   int64
		Name string
		Cid  string
	}
	err := c.BindJSON(&param)
	if err != nil {
		fmt.Println("some thing wrong")
	}
	result := db.Model(&model.Upload{}).Where("id = ?", param.Id).Updates(model.Upload{Name: param.Name, Cid: param.Cid})
	if result.Error != nil {
		c.String(http.StatusOK, "%s", "数据插入失败")
		panic("数据更新失败")
	}
	c.String(http.StatusOK, "%s", "更新成功")
}
func DelUploads(c *gin.Context) {
	var param struct {
		Id int64
	}
	err := c.BindJSON(&param)
	if err != nil {
		fmt.Println("some thing wrong")
	}
	result := db.Delete(&model.Upload{}, param.Id)
	if result.Error != nil {
		c.String(http.StatusOK, "%s", "数据插入失败")
		panic("数据删除失败")
	}
	c.String(http.StatusOK, "%s", "更新成功")
}

func SearchScope(c *gin.Context) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		query := db
		search := c.QueryMap("search")
		fmt.Println(search)
		if search["value"] != "" {
			query = query.Where("name LIKE ? ", "%"+search["value"]+"%")
			query = query.Or("cid LIKE ? ", "%"+search["value"]+"%")
		}
		return query
	}
}

func CreateM3u(c *gin.Context) {
	var uploads []model.Upload
	contents := "#EXTM3U\n"
	db.Find(&uploads).Order("updated_at desc")
	for _, element := range uploads {
		contents += "#EXTINF:-1 ," + element.Name + "\n" + inc.GetConfig("ipfsgateway") + "/" + element.Cid + "/" + element.Name + "\n"
	}

	file, err := os.Create("public/videos.m3u")
	if err != nil {
		log.Fatal(err)
	}
	file.WriteString(contents)
	file.Close()
	c.String(http.StatusOK, "%s", "生成成功")
}
