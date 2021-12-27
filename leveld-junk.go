package main

import (
	"chainrunner/internal/util"
	"chainrunner/services"
	"fmt"
	"log"
)

func levldb() {
	db, err := services.NewLvlDB("/home/mithril/.mev/db/chainrunner-test")

	if err != nil {
		log.Fatal("Error initializing db", err)
	}

	defer util.Duration(util.Track("setKey"))
	err = db.SetByKey([]byte("hello"), []byte("world"))

	if err != nil {
		log.Fatal("error setting db", err)
	}

	defer util.Duration(util.Track("getKey"))
	data, err := db.GetByKey([]byte("hello"))

	if err != nil {
		log.Fatal("error setting db", err)
	}

	fmt.Printf("data: %v", string(data))
}