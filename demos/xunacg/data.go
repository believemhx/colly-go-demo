package xunacg

import (
	"encoding/json"
	"io/ioutil"
)

type XunAcgUser struct {
	Name   string
	Uid    int
	Cookie string
	Count  int
	Status bool
}
type Config struct {
	Users []XunAcgUser
}

func GetData() (config Config, err error) {

	content, _ := ioutil.ReadFile("./config/xunacg.json")

	err = json.Unmarshal(content, &config)
	return
}
