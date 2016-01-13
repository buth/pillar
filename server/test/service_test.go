package test

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/coralproject/pillar/server/model"
	"github.com/coralproject/pillar/server/service"
)

const dataUsers = "fixtures/users.json"
const dataAssets = "fixtures/assets.json"
const dataComments = "fixtures/comments.json"
const dataActions = "fixtures/actions.json"

func TestCreateAsset(t *testing.T) {
	file, err := os.Open(dataAssets)
	if err != nil {
		fmt.Println("opening config file", err.Error())
	}

	objects := []model.Asset{}
	jsonParser := json.NewDecoder(file)
	if err = jsonParser.Decode(&objects); err != nil {
		fmt.Println("Error reading asset data", err.Error())
	}

	for _, one := range objects {
		_, err := service.CreateAsset(&one)
		if err != nil {
			t.Fail()
		}
	}
}

func TestCreateUser(t *testing.T) {
	file, err := os.Open(dataUsers)
	if err != nil {
		fmt.Println("opening config file", err.Error())
	}

	objects := []model.User{}
	jsonParser := json.NewDecoder(file)
	if err = jsonParser.Decode(&objects); err != nil {
		fmt.Println("Error reading user data", err.Error())
	}

	for _, one := range objects {
		_, err := service.CreateUser(&one)
		if err != nil {
			t.Fail()
		}
	}
}

func TestCreateComments(t *testing.T) {
	file, err := os.Open(dataComments)
	if err != nil {
		fmt.Println("opening config file", err.Error())
	}

	objects := []model.Comment{}
	jsonParser := json.NewDecoder(file)
	if err = jsonParser.Decode(&objects); err != nil {
		fmt.Println("Error reading user data", err.Error())
	}

	for _, one := range objects {
		_, err := service.CreateComment(&one)
		if err != nil {
			t.Fail()
		}
	}
}

func TestCreateActions(t *testing.T) {
	file, err := os.Open(dataActions)
	if err != nil {
		fmt.Println("opening config file", err.Error())
	}

	objects := []model.Action{}
	jsonParser := json.NewDecoder(file)
	if err = jsonParser.Decode(&objects); err != nil {
		fmt.Println("Error reading user data", err.Error())
	}

	for _, one := range objects {
		_, err := service.CreateAction(&one)
		if err != nil {
			t.Fail()
		}
	}
}