package model

import (
	"gorm.io/gorm"
)

type Upload struct {
	gorm.Model
	Name string
	Cid  string
	Size int64
}

type DataResult struct {
	Total    int64    `json:"recordsTotal"`
	Filtered int64    `json:"recordsFiltered"`
	Data     []Upload `json:"data"`
}
