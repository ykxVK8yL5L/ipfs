package inc

import (
	"gopkg.in/ini.v1"
)

func GetConfig(k string) string {
	cfg, err := ini.Load("config/config.ini")
	if err != nil {
		panic("无法读取配置文件，请检查配置文件")
	}
	val := cfg.Section("").Key(k).String()
	return val
}

func SetConfig(k string, v string) {
	cfg, err := ini.Load("config/config.ini")
	if err != nil {
		panic("无法读取配置文件，请检查配置文件")
	}
	cfg.Section("").Key(k).SetValue(v)
	cfg.SaveTo("config/config.ini")
}
