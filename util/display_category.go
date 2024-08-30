package util

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"serviceNest/config"
	"serviceNest/model"
)

func DisplayCategory() {
	var category []model.Category

	file, err := ioutil.ReadFile(config.FILENAME)
	if err != nil {
		fmt.Println(err)
	}
	err = json.Unmarshal(file, &category)
	if err != nil {
		fmt.Println(err)
	}
	for i, service := range category {
		fmt.Printf("%d Name : %s Description : %s", i+1, service.Name, service.Description)
		fmt.Println()
	}
}
