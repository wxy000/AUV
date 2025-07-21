package handlers

import (
	"AUV/common/response"
	"AUV/config"
	"AUV/db/repository"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func Login(c *gin.Context) {
	var creds struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&creds); err != nil {
		response.Fail(c, http.StatusBadRequest, "无效请求格式")
		return
	}

	user, err := repository.UserRepo.GetByUsername(creds.Username)
	if err != nil {
		response.Fail(c, http.StatusUnauthorized, "用户不存在")
		return
	}

	// 验证密码
	// 避免变量名冲突，使用 passwordErr 代替 err
	passwordErr := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(creds.Password))
	if passwordErr != nil {
		response.Fail(c, http.StatusUnauthorized, "密码错误")
		return
	}

	// 管理员权限验证
	if user.Role != "admin" && (user.Username != config.Cfg.Admin.Username || bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(config.Cfg.Admin.Password)) != nil) {
		response.Fail(c, http.StatusForbidden, "非管理员禁止登录")
		return
	}

	// 更新最后登录时间
	user.LastLogin = time.Now()
	if updateErr := repository.UserRepo.Update(user); updateErr != nil {
		response.Fail(c, http.StatusInternalServerError, "系统错误")
		return
	}

	// 生成JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":      user.ID,
		"username": user.Username,
		"exp":      time.Now().Add(time.Hour * time.Duration(config.Cfg.JWT.ExpiresHours)).Unix(),
	})

	tokenString, err := token.SignedString([]byte(config.Cfg.JWT.Secret))
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "令牌生成失败")
		return
	}

	response.Success(c, gin.H{
		"access_token": tokenString,
		"expires_in":   config.Cfg.JWT.ExpiresHours * 3600,
		"user_id":      user.ID,
	})
}

func RefreshToken(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		response.Fail(c, http.StatusUnauthorized, "无效令牌")
		return
	}

	// 通过userID查询用户信息
	user, err := repository.UserRepo.GetByID(userID.(uint))
	if err != nil {
		response.Fail(c, http.StatusUnauthorized, "用户不存在")
		return
	}

	claimsInterface, _ := c.Get("jwtClaims")
	claims, ok := claimsInterface.(jwt.MapClaims)
	if !ok {
		response.Fail(c, http.StatusUnauthorized, "无效令牌")
		return
	}

	exp, ok := claims["exp"].(float64)
	if !ok {
		response.Fail(c, http.StatusUnauthorized, "无效令牌")
		return
	}
	expirationTime := time.Unix(int64(exp), 0)

	refreshWindow := time.Hour * time.Duration(config.Cfg.JWT.RefreshWindowHours)
	if time.Now().After(expirationTime.Add(refreshWindow)) {
		response.Fail(c, http.StatusUnauthorized, "令牌已过期，无法刷新")
		return
	}

	newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":      userID,
		"username": user.Username,
		"exp":      time.Now().Add(time.Hour * time.Duration(config.Cfg.JWT.ExpiresHours)).Unix(),
	})

	tokenString, err := newToken.SignedString([]byte(config.Cfg.JWT.Secret))
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "令牌生成失败")
		return
	}

	response.Success(c, gin.H{
		"access_token": tokenString,
		"expires_in":   config.Cfg.JWT.ExpiresHours * 3600,
		"user_id":      user.ID,
	})
}
