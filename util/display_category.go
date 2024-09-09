package util

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"serviceNest/config"
	"serviceNest/model"
)

var Print = fmt.Printf
var ReadFile = ioutil.ReadFile

func DisplayCategory() {
	var category []model.Category

	file, err := ReadFile(config.FILENAME)
	if err != nil {
		fmt.Println(err)
	}
	err = json.Unmarshal(file, &category)
	if err != nil {
		fmt.Println(err)
	}
	for i, service := range category {
		Print("%d Name : %s Description : %s\n", i+1, service.Name, service.Description)
		fmt.Println()
	}
}
