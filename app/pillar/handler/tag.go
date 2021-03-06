package handler

import (
	"encoding/json"
	"github.com/coralproject/pillar/pkg/crud"
	"net/http"
)

func GetTags(w http.ResponseWriter, r *http.Request) {
	jsonObject := crud.Tag{}
	json.NewDecoder(r.Body).Decode(&jsonObject)

	// Write content-type, status code and payload
	w.Header().Set("Content-Type", "application/json")
	dbObject, err := crud.GetTags()
	doRespond(w, dbObject, err)
}

func UpsertTag(w http.ResponseWriter, r *http.Request) {
	//Get the tag from request
	jsonObject := crud.Tag{}
	json.NewDecoder(r.Body).Decode(&jsonObject)

	// Write content-type, status code and payload
	w.Header().Set("Content-Type", "application/json")
	dbObject, err := crud.UpsertTag(&jsonObject)
	doRespond(w, dbObject, err)
}

func DeleteTag(w http.ResponseWriter, r *http.Request) {
	//Get the tag from request
	jsonObject := crud.Tag{}
	json.NewDecoder(r.Body).Decode(&jsonObject)

	// Write content-type, status code and payload
	w.Header().Set("Content-Type", "application/json")
	err := crud.DeleteTag(&jsonObject)
	doRespond(w, nil, err)
}

