package main

import (
	"blog-go/config"
	"blog-go/database"
	"blog-go/log"
	"fmt"
	"testing"
)

func TestMyTest(t *testing.T) {
	config.Setup()
	log.Setup()
	database.Setup()
	rdb, rctx := database.GetRedisClient()
	r, err := rdb.HGetAll(rctx, "test").Result()
	if err != nil {
		t.Error(err)
	} else {
		t.Log(r)
		fmt.Println(r)
	}
}
