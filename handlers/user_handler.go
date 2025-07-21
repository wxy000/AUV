package handlers

import (
	"AUV/common/response"
	"AUV/config"
	"AUV/db/repository"
	"AUV/models"
	"net/http"
	"regexp"
	"strconv"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func CreateUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		response.Fail(c, http.StatusBadRequest, "无效的请求格式")
		return
	}

	// 清空前端可能传递的密码字段
	user.Password = ""

	// 手机号正则验证
	if user.Phone != "" {
		phoneRegex := regexp.MustCompile(`^1[3-9]\d{9}$`)
		if !phoneRegex.MatchString(user.Phone) {
			response.Fail(c, http.StatusBadRequest, "无效的手机号码格式")
			return
		}
	}

	// 邮箱正则验证
	if user.Email != "" {
		emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$`)
		if !emailRegex.MatchString(user.Email) {
			response.Fail(c, http.StatusBadRequest, "无效的邮箱地址格式")
			return
		}
	}

	// 设置默认密码并加密
	defaultPassword := config.Cfg.User.DefaultPassword
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(defaultPassword), bcrypt.DefaultCost)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "密码加密失败")
		return
	}
	user.Password = string(hashedPassword)

	if err := repository.UserRepo.Create(&user); err != nil {
		response.Fail(c, http.StatusInternalServerError, "创建用户失败")
		return
	}

	response.Success(c, user)
}

func GetUsers(c *gin.Context) {
	var users []models.User

	users, err := repository.UserRepo.GetAll()
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询用户失败")
		return
	}

	response.Success(c, users)
}

func GetInactiveUsers(c *gin.Context) {
	var users []models.User

	users, err := repository.UserRepo.GetAllInActive()
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询用户失败")
		return
	}

	response.Success(c, users)
}

func GetUser(c *gin.Context) {
	userId := c.Param("userId")
	// 将字符串类型的 userId 转换为 uint 类型
	id, err := strconv.ParseUint(userId, 10, 64)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "无效的用户ID")
		return
	}

	var user *models.User

	user, err = repository.UserRepo.GetByID(uint(id))
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "获取用户失败")
		return
	}

	response.Success(c, user)
}

func UpdateUserStatus(c *gin.Context) {
	userId := c.Param("userId")

	var req struct {
		IsActive bool `json:"isActive"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "无效的请求格式")
		return
	}

	if err := repository.UserRepo.UpdateUserStatus(userId, req.IsActive); err != nil {
		response.Fail(c, http.StatusInternalServerError, "更新用户状态失败")
		return
	}

	response.SuccessWithMessage(c, "用户状态更新成功", nil)
}

func DeleteUser(c *gin.Context) {
	userId := c.Param("userId")

	if err := repository.UserRepo.DeleteUser(userId); err != nil {
		response.Fail(c, http.StatusInternalServerError, "删除用户失败")
		return
	}

	response.SuccessWithMessage(c, "删除用户成功", nil)
}

func UpdateUser(c *gin.Context) {
	userId := c.Param("userId")
	// 将字符串类型的 userId 转换为 uint 类型
	id, err := strconv.ParseUint(userId, 10, 64)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "无效的用户ID")
		return
	}
	var updateData models.User

	if err := c.ShouldBindJSON(&updateData); err != nil {
		response.Fail(c, http.StatusBadRequest, "无效的请求格式")
		return
	}

	if err := repository.UserRepo.UpdateInfo(uint(id), &models.User{
		Username: updateData.Username,
		RealName: updateData.RealName,
		Gender:   updateData.Gender,
		Phone:    updateData.Phone,
		Email:    updateData.Email,
		Role:     updateData.Role,
		Remark:   updateData.Remark,
	}); err != nil {
		response.Fail(c, http.StatusInternalServerError, "更新用户失败")
		return
	}

	response.SuccessWithMessage(c, "用户更新成功", nil)
}
