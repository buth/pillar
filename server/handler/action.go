package handler

import (
	"encoding/json"
	"github.com/coralproject/pillar/server/model"
	"github.com/coralproject/pillar/server/service"
	"net/http"
)

//ImportAction imports actions into the system
func ImportAction(w http.ResponseWriter, r *http.Request) {
	//Get the user from request
	jsonObject := model.Action{}
	json.NewDecoder(r.Body).Decode(&jsonObject)

	// Write content-type, status code and payload
	w.Header().Set("Content-Type", "application/json")
	dbObject, err := service.CreateAction(jsonObject)
	doRespond(w, dbObject, err)
}