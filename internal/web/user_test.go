package web

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
	"testing"
	"we_book/internal/domain"
	"we_book/internal/service"
	svcmocks "we_book/internal/service/mocks"
)

func TestUserHandler_SignUp(t *testing.T) {
	testCases := []struct {
		name string
		mock func(ctrl *gomock.Controller) service.UserService

		reqBody    string
		wantedCode int
		wantedBody string
	}{

		{
			name: "register success",
			mock: func(ctrl *gomock.Controller) service.UserService {
				usersvc := svcmocks.NewMockUserService(ctrl)
				usersvc.EXPECT().SignUp(gomock.Any(), gomock.AssignableToTypeOf(domain.User{})).Return(nil)
				return usersvc
			},

			reqBody: `
		{
			"email": "123@qq.com",
			"password": "Bogzc2002@",
			"confirm_password": "Bogzc2002@"
		}
		`,
			wantedCode: http.StatusOK,
			wantedBody: "success",
		},

		{
			name: "email format error",
			mock: func(ctrl *gomock.Controller) service.UserService {
				usersvc := svcmocks.NewMockUserService(ctrl)
				return usersvc
			},

			reqBody: `
					{
						"email": "123@q",
						"password": "hello#world123",
						"confirmPassword": "hello#world123"
					}
					`,
			wantedCode: http.StatusOK,
			wantedBody: "email format error",
		},
		{
			name: "password format error",
			mock: func(ctrl *gomock.Controller) service.UserService {
				usersvc := svcmocks.NewMockUserService(ctrl)
				return usersvc
			},

			reqBody: `
					{
						"email": "user@test.com",
						"password": "1223214124",
						"confirmPassword": "1223214124"
					}
					`,
			wantedCode: http.StatusOK,
			wantedBody: "password format error",
		},
		{
			name: "password not match",
			mock: func(ctrl *gomock.Controller) service.UserService {
				usersvc := svcmocks.NewMockUserService(ctrl)
				return usersvc
			},

			reqBody: `
					{
						"email": "user@test.com",
						"password": "Bogzc2002@",
						"confirmPassword": "Bogzc2002!"
					}
					`,
			wantedCode: http.StatusOK,
			wantedBody: "password not match",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			// 创建一个 gin server
			server := gin.Default()
			// 不会使用到 code
			h := NewUserHandler(tc.mock(ctrl), nil)
			h.RegisterRoutes(server)

			req, err := http.NewRequest(http.MethodPost,
				"/users/signup", bytes.NewBuffer([]byte(tc.reqBody)))
			require.NoError(t, err)
			// 数据类型为 JSON
			req.Header.Set("Content-Type", "application/json")

			resp := httptest.NewRecorder()
			server.ServeHTTP(resp, req)
			fmt.Println(resp.Code)

			assert.Equal(t, tc.wantedCode, resp.Code)
			assert.Equal(t, tc.wantedBody, resp.Body.String())

		})
	}

}

func TestMock(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	usersvc := svcmocks.NewMockUserService(ctrl)
	usersvc.EXPECT().SignUp(gomock.Any(), gomock.Any()).Return(errors.New("mock error"))
	err := usersvc.SignUp(context.Background(), domain.User{
		Email: "user@example.com",
	})
	t.Log(err)
}
