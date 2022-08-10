package queue

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"testing"
	"time"
)

type helloJob struct {
	si StatusInfo
}

func (h *helloJob) InitStatus() {
	h.si.Name = "helloJob"
	h.si.Status = None
}

func (h *helloJob) Run(si chan StatusInfo) (StatusInfo, error) {
	h.si.Status = Running
	si <- h.si

	//time.Sleep(time.Second * 5)

	h.si.Status = Succeeded
	si <- h.si

	return h.si, nil
}

func (h *helloJob) Status() StatusInfo {
	return h.si
}

type hiJob struct {
	si StatusInfo
}

func (h *hiJob) InitStatus() {
	h.si.Name = "hiJob"
	h.si.Status = None
}

func (h *hiJob) Run(si chan StatusInfo) (StatusInfo, error) {
	h.si.Status = Running
	si <- h.si

	//time.Sleep(time.Second * 5)

	h.si.Status = Failed
	err := fmt.Errorf("hiJob error")
	h.si.Error = err.Error()
	si <- h.si

	return h.si, err
}

func (h *hiJob) Status() StatusInfo {
	return h.si
}

func TestAll(t *testing.T) {

	hello := helloJob{}
	hi := hiJob{}

	db := getGormDB()
	//id := uuid.New().String()
	id := "4e54212a-b5e4-4c3b-8b10-7e13b2f82aa2"

	_, err := NewQueue(id, &hello, &hi)
	if err != nil {
		t.Fatal(err)
	}

	q, err := GetQueue(id)
	if err != nil {
		t.Error(err)
	}

	// load status
	err = q.LoadStatus(db)
	if err != nil {
		fmt.Println("load status error:", err)
	}

	// defer save status
	defer func() {
		err = q.SaveStatus(db)
		if err != nil {
			t.Error(err)
		}
	}()

	// start queue, in a new goroutine
	done := make(chan struct{})
	go q.Start(done)

	// view status, in a new goroutine
	go func() {
		for {
			time.Sleep(time.Second)
			info, err := q.GetStatus()
			if err != nil {
				t.Error(err)
			}
			for _, v := range info {
				fmt.Print(v, "; ")
			}
			fmt.Println()
		}
	}()

	// wait
	<-done

	// check the status after finishing
	info, err := q.GetStatus()
	if err != nil {
		t.Error(err)
	}
	fmt.Println(info)
}

func getGormDB() *gorm.DB {
	// get sqlite db
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	return db
}

func TestQueueJobDuplicate(t *testing.T) {
	hello := helloJob{}
	hi := hiJob{}
	hello2 := helloJob{}

	id := "4e54212a-b5e4-4c3b-8b10-7e13b2f82aa2"

	_, err := NewQueue(id, &hello, &hi, &hello2)
	assert.Error(t, err)
}