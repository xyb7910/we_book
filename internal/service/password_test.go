package service

import (
	"golang.org/x/crypto/bcrypt"
	"testing"
)

func TestPasswordEncrypt(t *testing.T) {
	pwd := []byte("123456")
	// 加密
	encryptpwd, err := bcrypt.GenerateFromPassword(pwd, bcrypt.DefaultCost)
	// 进行比较
	err = bcrypt.CompareHashAndPassword(encryptpwd, pwd)
	if err != nil {
		return
	}
}
