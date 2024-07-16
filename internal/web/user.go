package web

import (
	"fmt"
	regexp "github.com/dlclark/regexp2"
	"net/http"

	"github.com/gin-gonic/gin"
)

// UserHandler 我准备在它上面定义跟用户有关的路由
type UserHandler struct {
	emailExp    *regexp.Regexp
	passwordExp *regexp.Regexp
}

// NewUserHandler 一定要在main.go中调用这个函数，否则会出现路由注册失败的问题
func NewUserHandler() *UserHandler {
	const (
		emailRegexPattern    = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,4}$`
		passwordRegexPattern = `^(?=.*[A-Za-z])(?=.*\d)(?=.*[$@$!%*#?&])[A-Za-z\d$@$!%*#?&]{8,}$`
	)
	emailExp := regexp.MustCompile(emailRegexPattern, regexp.None)
	passwordExp := regexp.MustCompile(passwordRegexPattern, regexp.None)
	return &UserHandler{
		emailExp:    emailExp,
		passwordExp: passwordExp,
	}
}

// RegisterRoute 实现路由注册
func (u *UserHandler) RegisterRoute(server *gin.Engine) {
	// 首先进行路由分组
	user := server.Group("/users")

	// 定义其他路由
	user.POST("/signup", u.SignUp)
	user.POST("/signin", u.SignIn)
	user.GET("/profile", u.Profile)
	user.GET("/logout", u.Logout)
	user.GET("/edit", u.Edit)
}

// SignUp 实现 user 相关的 signup 接口
func (u *UserHandler) SignUp(ctx *gin.Context) {

	type SignUpReq struct {
		Email           string `json:"email"`
		ConfirmPassword string `json:"confirm_password"`
		Password        string `json:"password"`
	}

	var req SignUpReq
	// Bind 方法会跟据请求头中的 Content-Type 来自动选择解析方式，
	// 前端传递的数据格式必须与 Content-Type 一致
	// 自动将请求体中的数据绑定到 request 中
	// 如果绑定失败，则会返回错误 400
	if err := ctx.Bind(&req); err != nil {
		return
	}

	email := ctx.Query("email")
	password := ctx.Query("password")
	confirmPassword := ctx.Query("confirm_password")
	fmt.Println(email, password, confirmPassword)

	//邮箱格式校验
	if u.emailExp == nil {
		ctx.String(http.StatusInternalServerError, "email regexp error")
		return
	}
	isEmail, err := u.emailExp.MatchString(email)
	if err != nil {
		ctx.String(http.StatusInternalServerError, "email regexp error")
		return
	}
	if !isEmail {
		ctx.String(http.StatusBadRequest, "email format error")
		return
	}

	// 密码格式校验
	//if u.passwordExp == nil {
	//	ctx.String(http.StatusInternalServerError, "password regexp error")
	//	return
	//}
	isPassword, err := u.passwordExp.MatchString(password)
	if err != nil {
		ctx.String(http.StatusInternalServerError, "password regexp error")
		return
	}
	if !isPassword {
		ctx.String(http.StatusBadRequest, "password format error")
		return
	}

	// 检验两次密码是否一致
	if password != confirmPassword {
		ctx.String(http.StatusBadRequest, "password not match")
		return
	}
	ctx.String(http.StatusOK, "success")
}

// SignIn 实现 user 相关的 signin 接口
func (u *UserHandler) SignIn(context *gin.Context) {

}

// Profile 实现 user 相关的 profile 接口
func (u *UserHandler) Profile(context *gin.Context) {

}

// Logout 实现 user 相关的 logout 接口
func (u *UserHandler) Logout(context *gin.Context) {

}

// Edit 实现 user 相关的 edit 接口
func (u *UserHandler) Edit(context *gin.Context) {

}
