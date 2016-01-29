package service

import (
	"encoding/json"
	"fmt"
	"github.com/coralproject/pillar/config"
	"github.com/coralproject/pillar/pkg/model"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"net/http"
	"os"
	"time"
)

// AppError encapsulates application specific error
type AppError struct {
	Error   error
	Message string
	Code    int
}

var (
	mgoSession *mgo.Session
)

// MongoManager encapsulates a mongo session with all relevant collections
type MongoManager struct {
	Session   *mgo.Session
	Assets    *mgo.Collection
	Users     *mgo.Collection
	Actions   *mgo.Collection
	Comments  *mgo.Collection
	Tags      *mgo.Collection
	TagTarget *mgo.Collection
}

//Close closes the mongodb session; must be called, else the session remain open
func (manager *MongoManager) Close() {
	manager.Session.Close()
}

//export MONGODB_URL=mongodb://localhost:27017/coral
func init() {
	session, err := mgo.Dial(config.GetContext().MongoURL)
	if err != nil {
		log.Fatalf("Error connecting to mongo database: %s", err)
	}

	mgoSession = session

	//url and source.id on Asset
	mgoSession.DB("").C(model.CollectionAction).EnsureIndexKey("source.id")

	//url and source.id on Asset
	mgoSession.DB("").C(model.CollectionAsset).EnsureIndexKey("source.id")
	mgoSession.DB("").C(model.CollectionAsset).EnsureIndexKey("url")

	//source.id on User
	mgoSession.DB("").C(model.CollectionUser).EnsureIndexKey("source.id")

	//source.id on Comment
	mgoSession.DB("").C(model.CollectionComment).EnsureIndexKey("source.id")

	//name on Tag
	mgoSession.DB("").C(model.CollectionTag).EnsureIndexKey("name")

	//target_id, name and target,
	mgoSession.DB("").C(model.CollectionTag).EnsureIndexKey("target_id", "name", "target")
}

func initDB() {
	file, err := os.Open("dbindex.json")
	if err != nil {
		log.Fatalf("Error opening file %s\n", err.Error())
	}

	objects := []model.Index{}
	jsonParser := json.NewDecoder(file)
	if err = jsonParser.Decode(&objects); err != nil {
		log.Fatalf("Error reading index information %v\n", err)
	}

	for _, one := range objects {
		if err := CreateIndex(&one); err != nil {
			log.Fatalf("Error creating indexes %v\n", err)
		}
	}
}

//GetMongoManager returns a cloned MongoManager
func GetMongoManager() *MongoManager {

	manager := MongoManager{}
	manager.Session = mgoSession.Clone()
	manager.Users = manager.Session.DB("").C(model.CollectionUser)
	manager.Assets = manager.Session.DB("").C(model.CollectionAsset)
	manager.Actions = manager.Session.DB("").C(model.CollectionAction)
	manager.Comments = manager.Session.DB("").C(model.CollectionComment)
	manager.Tags = manager.Session.DB("").C(model.CollectionTag)
	manager.TagTarget = manager.Session.DB("").C(model.CollectionTagTarget)

	return &manager
}

// UpdateMetadata updates metadata for an entity
func UpdateMetadata(object *model.Metadata) (interface{}, *AppError) {

	manager := GetMongoManager()
	defer manager.Close()

	collection := manager.Session.DB("").C(object.Target)
	var dbEntity bson.M
	collection.FindId(object.TargetID).One(&dbEntity)
	if len(dbEntity) == 0 {
		collection.Find(bson.M{"source.id": object.Source.ID}).One(&dbEntity)
	}

	if len(dbEntity) == 0 {
		message := fmt.Sprintf("Cannot update metadata for [%+v]\n", object)
		return nil, &AppError{nil, message, http.StatusInternalServerError}
	}

	collection.Update(
		bson.M{"_id": dbEntity["_id"]},
		bson.M{"$set": bson.M{"metadata": object.Metadata}},
	)

	return dbEntity, nil
}

// CreateIndex creates indexes to various entities
func CreateIndex(object *model.Index) *AppError {
	manager := GetMongoManager()
	defer manager.Close()

	err := manager.Session.DB("").C(object.Target).EnsureIndex(object.Index)
	if err != nil {
		message := fmt.Sprintf("Error creating index [%+v]", object)
		return &AppError{err, message, http.StatusInternalServerError}
	}

	return nil
}

// UpsertTag adds/updates tags to the master list
func UpsertTag(object *model.Tag) (*model.Tag, *AppError) {
	manager := GetMongoManager()
	defer manager.Close()

	//set created-date for the new ones
	var dbEntity model.Tag
	if manager.Tags.FindId(object.Name).One(&dbEntity); dbEntity.Name == "" {
		object.DateCreated = time.Now()
	}

	object.DateUpdated = time.Now()
	_, err := manager.Tags.UpsertId(object.Name, object)
	if err != nil {
		message := fmt.Sprintf("Error creating tag [%+v]", object)
		return nil, &AppError{err, message, http.StatusInternalServerError}
	}
	fmt.Printf("Tag: %+v\n\n", object)

	return object, nil
}

// CreateTagStats creates TagStat entries for an entity
func CreateTagStats(manager *MongoManager, tags []string, tt *model.TagTarget) error {

	for _, name := range tags {

		tt.ID = bson.NewObjectId()
		tt.Name = name
		tt.DateCreated = time.Now()

		//skip the same entry, if exists
		dbEntity := model.TagTarget{}
		manager.TagTarget.Find(bson.M{"target_id": tt.TargetID, "name": name, "target": tt.Target}).One(&dbEntity)
		if dbEntity.ID != "" {
			continue
		}

		if err := manager.TagTarget.Insert(tt); err != nil {
			return err
		}
	}

	return nil
}