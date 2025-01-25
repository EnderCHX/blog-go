package blog

import (
	"blog-go/config"
	"blog-go/database"
	"blog-go/log"
	"fmt"
	"testing"
	"time"
)

func TestPassageID(t *testing.T) {
	passage := Passage{
		PassageId:      "test",
		AuthorUsername: "test",
		Title:          "test2",
		Content:        `# Go是世界上最好的语言`,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
		Deleted:        false,
	}
	passage.GenPassageId()
	fmt.Println(passage)
}

func TestGetPassages(t *testing.T) {
	config.Setup()
	log.Setup()
	database.Setup()
	var passages []Passage
	passages, err := GetPassagesByTag("test")
	if err != nil {
		t.Error(err)
	}
	fmt.Println(passages)
}
