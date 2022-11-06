package storage

import (
	"log"
	"testing"
	"time"
)
import "github.com/stretchr/testify/assert"

func TestFileStorage_GetAll(t *testing.T) {
	oFs := NewFileStorage()
	_, err := oFs.GetAll()
	assert.Nil(t, err)
}

func TestFileStorage_Watch(t *testing.T) {
	oFs := NewFileStorage()
	oFs.Init()
	go func() {
		tc := time.NewTicker(3 * time.Second)
		for range tc.C {
			Projects := oFs.Projects
			t.Log("current project info: ", Projects)
		}
	}()

	go func() {
		err := oFs.Watch()
		if err != nil {
			log.Println("receive event err:", err.Error())
		}
		log.Println("range ok....")
	}()

	go func() {
		log.Println("start range WatchEvent...")
		for c := range oFs.WatchEvent() {
			log.Println("receive channel ok:", c, c.Project)
		}
	}()

	t.Log("done")

	ddd := make(chan bool)
	<-ddd
}
