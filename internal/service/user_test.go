package service

import (
	"context"
	"github.com/go-playground/assert/v2"
	"go.uber.org/mock/gomock"
	"golang.org/x/crypto/bcrypt"
	"testing"
	"we_book/internal/domain"
	"we_book/internal/repository"
	repomocks "we_book/internal/repository/mocks"
	_ "we_book/internal/service/mocks"
)

func TestUserService_Login(t *testing.T) {
	testCases := []struct {
		name string                                                  // test case name
		mock func(ctrl *gomock.Controller) repository.UserRepository // mock function

		// input args
		email    string
		password string

		// output result
		wantUser    domain.User
		wantedError error
	}{
		{
			name: "login success",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				repo.EXPECT().FindByEmail(gomock.Any(), "user@example.com").
					Return(domain.User{
						Email:    "user@example.com",
						Password: "$2a$10$xHXwtabXN3.aMZ.rJVw4aecb5Xxj5lqPc3t3ADh1h17cqLnCckEuS",
						Phone:    "1456436558",
					}, nil)
				return repo
			},

			email:    "user@example.com",
			password: "Bogzc2002@",

			wantUser: domain.User{
				Email:    "user@example.com",
				Password: "$2a$10$xHXwtabXN3.aMZ.rJVw4aecb5Xxj5lqPc3t3ADh1h17cqLnCckEuS",
				Phone:    "1456436558",
			},
			wantedError: nil,
		},
		{
			name: "user not exist",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				repo.EXPECT().FindByEmail(gomock.Any(), "user@example.com").
					Return(domain.User{
						Email:    "user@example.com",
						Password: "$2a$10$xHXwtabXN3.aMZ.rJVw4aecb5Xxj5lqPc3t3ADh1h17cqLnCckEuS",
						Phone:    "1456436558",
					}, nil)
				return repo
			},
			email:    "user@example.com",
			password: "Bogzc2002er@",

			wantUser:    domain.User{},
			wantedError: ErrInvalidUserOrPassword,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			// 具体的测试代码
			svc := NewUserService(tc.mock(ctrl))
			u, err := svc.Login(context.Background(), tc.email, tc.password)
			assert.Equal(t, tc.wantUser, u)
			assert.Equal(t, tc.wantedError, err)
		})
	}
}

func TestEncryptPassword(t *testing.T) {
	res, err := bcrypt.GenerateFromPassword([]byte("Bogzc2002@"), bcrypt.DefaultCost)
	if err == nil {
		t.Log(string(res))
	}
}

//func Test(t *testing.T) {
//	// 首先创建一个 mock 的控制器
//	ctrl := gomock.NewController(t)
//	// 每个测试用例结束后都要调用 Finish 方法
//	defer ctrl.Finish()
//	// 创建一个模拟对象
//	usersvc := svcmocks.NewMockUserService(ctrl)
//	// 设计一个模拟调用
//	// 设置预期调用 Login 方法时，返回一个错误
//	usersvc.EXPECT().SignUp(gomock.Any(), gomock.Any()).
//		Return(errors.New("mock error"))
//
//}
