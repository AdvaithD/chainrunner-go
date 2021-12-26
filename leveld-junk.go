package main

import (
	"chainrunner/services"
	"chainrunner/util"
	"fmt"
	"log"
)

func levldb() {
	db, err := services.NewLvlDB("/home/mithril/.mev/db/chainrunner-test")

	if err != nil {
		log.Fatalf("Error initializing db", err)
	}

	defer util.Duration(util.Track("setKey"))
	err = db.SetByKey([]byte("hello"), []byte("world"))

	if err != nil {
		log.Fatalf("error setting db", err)
	}

	defer util.Duration(util.Track("getKey"))
	data, err := db.GetByKey([]byte("hello"))

	if err != nil {
		log.Fatalf("error setting db", err)
	}

	fmt.Printf("data: %v", string(data))
}