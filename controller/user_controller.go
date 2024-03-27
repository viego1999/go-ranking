package controller

import (
	"strconv"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"go-ranking/models"
	"go-ranking/pkg/logger"
	"go-ranking/utils"
)

type UserController struct{}

func (u UserController) Register(ctx *gin.Context) {
	// 接受用户名，密码，确认密码
	params := make(map[string]interface{})
	err := ctx.BindJSON(&params)
	if err != nil {
		ReturnError(ctx, 5001, "服务器内部错误："+err.Error())
		return
	}

	username := utils.GetOrDefault(params, "username", "").(string)
	password := utils.GetOrDefault(params, "password", "").(string)
	confirmPwd := utils.GetOrDefault(params, "confirmPassword", "").(string)

	if username == "" || password == "" || confirmPwd == "" {
		ReturnError(ctx, 4001, "请输入正确的信息")
		return
	}

	if password != confirmPwd {
		ReturnError(ctx, 4001, "密码和确认密码不一致")
		return
	}

	user, err := models.GetUserByUsername(username)
	if user.Id != 0 {
		ReturnError(ctx, 4001, "用户名已存在")
		return
	}

	_, err = models.AddUser(username, utils.EncryMd5(password))
	if err != nil {
		ReturnError(ctx, 4001, "注册用户失败："+err.Error())
		return
	}

	ReturnSuccess(ctx, 0, "Register success", user.Id, 1)
}

func (u UserController) Login(ctx *gin.Context) {
	// 接受用户名 密码
	params := make(map[string]interface{})
	err := ctx.BindJSON(&params)
	if err != nil {
		ReturnError(ctx, 5001, "服务器内部错误："+err.Error())
		return
	}

	username := utils.GetOrDefault(params, "username", "").(string)
	password := utils.GetOrDefault(params, "password", "").(string)

	if username == "" || password == "" {
		ReturnError(ctx, 4001, "请输入正确的信息")
		return
	}

	user, err := models.GetUserByUsername(username)
	if user.Id == 0 {
		ReturnError(ctx, 4004, "用户名或密码不正确")
		return
	}

	if user.Password != utils.EncryMd5(password) {
		ReturnError(ctx, 4004, "用户名或密码不正确")
		return
	}

	session := sessions.Default(ctx)
	session.Set("login:"+strconv.Itoa(user.Id), user.Id)
	err = session.Save()
	if err != nil {
		ReturnError(ctx, 5000, "session save failed!")
		return
	}

	ReturnSuccess(ctx, 0, "登录成功", models.UserApi{Id: user.Id, Username: user.Username}, 1)
}

func (u UserController) GetUserInfo(ctx *gin.Context) {
	logger.Write("日志信息", "user")
	// 路径方式获取请求参数
	idStr := ctx.Param("id")
	id, _ := strconv.Atoi(idStr)
	user, _ := models.GetUserById(id)
	ReturnSuccess(ctx, 0, "success", user, 1)
}

func (u UserController) GetUserList(ctx *gin.Context) {
	logger.Write("日志信息 - GetUserList", "user")
	users, err := models.GetUserList()
	if err != nil {
		ReturnError(ctx, 404, "not found related user list info, err: "+err.Error())
		return
	}
	ReturnSuccess(ctx, 0, "get user list success", users, int64(len(users)))
}
