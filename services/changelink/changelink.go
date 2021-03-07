package main

import (
	"fmt"
	"github.com/bennettaur/changelink/services/changelink/models"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
)

func main() {
	//lineRange := []models.LineRange{{models.UNBOUNDED, models.UNBOUNDED}}
	//changelink := models.NewLink("This file", "github.com", "services/changelink/watcher.go", lineRange)
	//
	//// Make sure pass the model by reference.
	//err := mgm.Coll(changelink).Create(changelink)

	link := &models.Watcher{}
	coll := mgm.Coll(link)

	err := coll.First(bson.M{}, link)

	if err != nil {
		panic(err)
	}

	fmt.Print(link)
}
