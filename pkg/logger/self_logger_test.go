package logger

import (
	"net/http"
	"testing"
)

var (
	baidu = "https://www.baidu.com"
)

func TestSelfLogger(t *testing.T) {
	l := InitSelfLogger()
	l.Debugf("Try to get %s", baidu)
	resp, err := http.Get(baidu)
	if err != nil {
		l.Errorf("Get %s failed, err:%v", baidu, err)
	} else {
		l.Infof("Get %s success, status:%d", baidu, resp.StatusCode)
		_ = resp.Body.Close()

	}
}
