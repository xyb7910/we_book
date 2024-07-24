package web

import (
	"fmt"
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-contrib/sessions"
	"net/http"
	"time"
	"we_book/internal/domain"
	"we_book/internal/service"

	"github.com/gin-gonic/gin"
)

var biz = "login"

// UserHandler 我准备在它上面定义跟用户有关的路由
type UserHandler struct {
	emailExp    *regexp.Regexp
	passwordExp *regexp.Regexp
	svc         service.UserService
	codeSvc     service.CodeService
	jwtHandler
}

// NewUserHandler 一定要在main.go中调用这个函数，否则会出现路由注册失败的问题
func NewUserHandler(svc service.UserService, codeSvc service.CodeService) *UserHandler {
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
		codeSvc:     codeSvc,
	}
}

// RegisterRoutes 实现路由注册
func (u *UserHandler) RegisterRoutes(server *gin.Engine) {
	// 首先进行路由分组
	user := server.Group("/users")

	// 定义其他路由
	user.POST("/signup", u.SignUp)
	user.POST("/login", u.Login)
	user.GET("/profile", u.ProfileJWT)
	user.GET("/logout", u.Logout)
	user.GET("/edit", u.Edit)
	user.POST("/login_sms/code/send", u.SendLoginSMSCode)
	user.POST("/login_sms", u.VerifyLoginSMSCode)
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
	if err := u.setJWTToken(ctx, user.Id); err != nil {
		return
	}
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
	type UserProfile struct {
		NickName     string `json:"nick_name"`
		Birthday     string `json:"birthday"`
		Email        string `json:"email"`
		Introduction string `json:"introduction"`
	}
	// 通过 jwt 荷载的信息来获取用户的 id
	uc, ok := ctx.MustGet("claims").(*UserClaims)
	if !ok {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	user, err := u.svc.FindById(ctx, uc.Uid)
	if err != nil {
		ctx.String(http.StatusOK, "internal server error, FindById failed")
		return
	}
	ctx.JSON(http.StatusOK, UserProfile{
		NickName:     user.NickName,
		Birthday:     user.Birthday.Format("2006-01-02"),
		Email:        user.Email,
		Introduction: user.Introduction,
	})
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
func (u *UserHandler) Edit(ctx *gin.Context) {
	type EditReq struct {
		NickName     string `json:"nick_name"`
		Birthday     string `json:"birthday"`
		Introduction string `json:"introduction"`
	}

	var req EditReq
	if err := ctx.Bind(&req); err != nil {
		return
	}
	// 考虑 基于 jwt 携带参数来实现
	uc, ok := ctx.MustGet("claims").(*UserClaims)
	if !ok {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	// 对用户输入的信息进行校验
	birthday, err := time.Parse("2006-01-02", req.Birthday)
	if err != nil {
		ctx.String(http.StatusOK, "birthday format error")
		return
	}
	err = u.svc.Edit(ctx, domain.User{
		Id:           uc.Uid,
		NickName:     req.NickName,
		Birthday:     birthday,
		Introduction: req.Introduction,
	})
	if err != nil {
		ctx.String(http.StatusOK, "internal server error, Edit failed")
		return
	}
	ctx.String(http.StatusOK, "success")
}

func (u *UserHandler) SendLoginSMSCode(ctx *gin.Context) {
	type Request struct {
		Phone string `json:"phone"`
	}
	var req Request
	if err := ctx.Bind(&req); err != nil {
		return
	}
	err := u.codeSvc.Send(ctx, biz, req.Phone)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 4,
			Msg:  "send sms code failed",
		})
	}
	ctx.JSON(http.StatusOK, Result{
		Msg: "send success",
	})
}

func (u *UserHandler) VerifyLoginSMSCode(ctx *gin.Context) {
	type Request struct {
		Phone string
		Code  string
	}
	var req Request
	if err := ctx.Bind(&req); err != nil {
		return
	}
	ok, err := u.codeSvc.Verify(ctx, biz, req.Phone, req.Code)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "backend error",
		})
		return
	}
	if !ok {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "verify code failed",
		})
	}
	// 接下来，如果系统有这个手机号，则直接登录
	// 如果没有，则进行注册
	user, err := u.svc.FindOrCreate(ctx, req.Phone)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "backend error",
		})
	}
	//fmt.Println(user.Id)
	if err := u.setJWTToken(ctx, user.Id); err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "backend error",
		})
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Code: 2,
		Msg:  "login success",
	})

}
