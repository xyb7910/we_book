package main

import (
	"context"
	corn "github.com/robfig/cron/v3"
	"log"
	"testing"
	"time"
)

func TestTicker(t *testing.T) {
	ticker := time.NewTicker(time.Second * 2)
	// 注意停止 否则会一直循环
	defer ticker.Stop()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	// 注意取消
	defer cancel()
	for {
		select {
		case <-ctx.Done():
			t.Log("time out")
			goto END
		case now := <-ticker.C:
			t.Log(now.Format("2006-01-02 15:04:05"))
		}
	}
END:
	t.Log("end")
}

func TestTimer(t *testing.T) {
	tm := time.NewTimer(time.Second * 2)
	defer tm.Stop()
	go func() {
		for now := range tm.C {
			t.Log(now.Format("2006-01-02 15:04:05"))
		}
	}()
	time.Sleep(time.Second * 3)
}

type myJob struct {
}

func (m myJob) Run() {
	log.Println("run", time.Now().Format("2006-01-02 15:04:05"))
}

func TestCornExpr(t *testing.T) {
	expr := corn.New(corn.WithSeconds())

	//expr.AddJob("@every 1s", myJob{})
	expr.AddFunc("@every 3s", func() {
		t.Log("begin print")
		t.Log("run", time.Now().Format("2006-01-02 15:04:05"))
		time.Sleep(time.Second * 2)
		t.Log("end print")
	})
	expr.Start()
	time.Sleep(10 * time.Second)
	stop := expr.Stop()
	t.Log("stop", stop)
	<-stop.Done()
	t.Log("end")
}
