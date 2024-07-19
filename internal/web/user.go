package web

import (
	"fmt"
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-contrib/sessions"
	jwt "github.com/golang-jwt/jwt/v5"
	"net/http"
	"time"
	"we_book/internal/domain"
	"we_book/internal/service"

	"github.com/gin-gonic/gin"
)

// UserHandler 我准备在它上面定义跟用户有关的路由
type UserHandler struct {
	emailExp    *regexp.Regexp
	passwordExp *regexp.Regexp
	svc         *service.UserService
}

// NewUserHandler 一定要在main.go中调用这个函数，否则会出现路由注册失败的问题
func NewUserHandler(svc *service.UserService) *UserHandler {
	const (
		emailRegexPattern    = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,4}$`
		passwordRegexPattern = `^(?=.*[A-Za-z])(?=.*\d)(?=.*[$@$!%*#?&])[A-Za-z\d$@$!%*#?&]{8,}$`
	)
	emailExp := regexp.MustCompile(emailRegexPattern, regexp.None)
	passwordExp := regexp.MustCompile(passwordRegexPattern, regexp.None)
	return &UserHandler{
		emailExp:    emailExp,
		passwordExp: passwordExp,
		svc:         svc,
	}
}

// RegisterRoute 实现路由注册
func (u *UserHandler) RegisterRoute(server *gin.Engine) {
	// 首先进行路由分组
	user := server.Group("/users")

	// 定义其他路由
	user.POST("/signup", u.SignUp)
	user.POST("/login", u.Login)
	user.GET("/profile", u.ProfileJWT)
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
	//email, password, confirmPassword := req.Email, req.Password, req.ConfirmPassword
	fmt.Println(email, password, confirmPassword)

	//邮箱格式校验
	if u.emailExp == nil {
		ctx.String(http.StatusOK, "email regexp error")
		return
	}
	isEmail, err := u.emailExp.MatchString(req.Email)
	if err != nil {
		ctx.String(http.StatusOK, "email regexp error")
		return
	}
	if !isEmail {
		ctx.String(http.StatusOK, "email format error")
		return
	}

	// 密码格式校验
	//if u.passwordExp == nil {
	//	ctx.String(http.StatusInternalServerError, "password regexp error")
	//	return
	//}
	isPassword, err := u.passwordExp.MatchString(req.Password)
	if err != nil {
		ctx.String(http.StatusOK, "password regexp error")
		return
	}
	if !isPassword {
		ctx.String(http.StatusOK, "password format error")
		return
	}

	// 检验两次密码是否一致
	if req.Password != req.ConfirmPassword {
		ctx.String(http.StatusOK, "password not match")
		return
	}
	//ctx.String(http.StatusOK, "success")

	// 调用 service 层的方法 保存用户信息， 并且进行相关的校验
	err = u.svc.SignUp(ctx, domain.User{
		Email:    email,
		Password: password,
	})

	if err == service.ErrUserDuplicateEmail {
		ctx.String(http.StatusOK, "email already exists")
		return
	}
	if err != nil {
		ctx.String(http.StatusOK, "internal server error, SigUp failed")
		return
	}
	ctx.String(http.StatusOK, "success")
}

func (u *UserHandler) LoginJWT(ctx *gin.Context) {
	type LoginReq struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var req LoginReq

	if err := ctx.Bind(&req); err != nil {
		return
	}
	user, err := u.svc.Login(ctx, req.Email, req.Password)
	if err == service.ErrInvalidUserOrPassword {
		ctx.String(http.StatusOK, "invalid user or password")
		return
	}
	if err != nil {
		ctx.String(http.StatusOK, "internal server error, Login failed")
		return
	}
	// 使用 JWT 进行登录 并将 用户的 id 存储在 token 中
	claims := UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute)),
		},
		Uid: user.Id,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	tokenStr, err := token.SignedString([]byte("95osj3fUD7fo0mlYdDbncXz4VD2igvf0"))
	if err != nil {
		ctx.String(http.StatusInternalServerError, "internal server error, Login failed")
		return
	}
	ctx.Header("x-jwt-token", tokenStr)
	fmt.Println(user)
	ctx.String(http.StatusOK, "success")
	return
}

// Login 实现 user 相关的 Login 接口
func (u *UserHandler) Login(ctx *gin.Context) {
	type LoginReq struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var req LoginReq

	if err := ctx.Bind(&req); err != nil {
		return
	}
	email := ctx.Query("email")
	password := ctx.Query("password")
	fmt.Println(email, password)
	user, err := u.svc.Login(ctx, req.Email, req.Password)
	if err == service.ErrInvalidUserOrPassword {
		ctx.String(http.StatusOK, "invalid user or password")
		return
	}
	if err != nil {
		ctx.String(http.StatusOK, "internal server error, Login failed")
		return
	}
	// 设置 session
	sess := sessions.Default(ctx)
	sess.Set("user_id", user.Id)
	sess.Options(sessions.Options{
		Secure:   true,
		HttpOnly: true,
		MaxAge:   60,
	})
	err = sess.Save()
	if err != nil {
		return
	}
	ctx.String(http.StatusOK, "success")
	return
}

// ProfileJWT 实现 user 相关的 profile 接口
func (u *UserHandler) ProfileJWT(ctx *gin.Context) {
	c, _ := ctx.Get("claims")
	claims, ok := c.(*UserClaims)
	if !ok {
		ctx.String(http.StatusOK, "invalid token")
		return
	}
	println(claims.Uid)
	ctx.String(http.StatusOK, "This is profile")
}

// Profile 实现 user 相关的 profile 接口
func (u *UserHandler) Profile(ctx *gin.Context) {

	ctx.String(http.StatusOK, "This is profile")
}

// Logout 实现 user 相关的 logout 接口
func (u *UserHandler) Logout(ctx *gin.Context) {
	sess := sessions.Default(ctx)
	sess.Options(sessions.Options{
		Secure:   true,
		HttpOnly: true,
		MaxAge:   60,
	})
	sess.Save()
	ctx.String(http.StatusOK, "success")
}

// Edit 实现 user 相关的 edit 接口
func (u *UserHandler) Edit(context *gin.Context) {
	//type EditReq struct {
	//	NickName     string
	//	Birthday     string
	//	introduction string
	//}
	//
	//var req EditReq
	//if err := context.Bind(&req); err != nil {
	//	return
	//}
	//nickName := context.Query("nick_name")
	//birthday := context.Query("birthday")
	//introduction := context.Query("introduction")
	//err := u.svc.Edit(ctx, domain.UserInfo{
	//	NickName:     nickName,
	//	Birthday:     birthday,
	//	Introduction: introduction,
	//})
	//if err != nil {
	//	return
	//}
}

type UserClaims struct {
	jwt.RegisteredClaims

	// 声明自己的字段
	Uid int64 `json:"uid"`
}
