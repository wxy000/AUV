package models

import (
	"time"

	"gorm.io/gorm"
)

type GenderType string

const (
	Male   GenderType = "male"
	Female GenderType = "female"
	Other  GenderType = "unknown"
)

type User struct {
	gorm.Model
	Username  string     `json:"username" gorm:"size:20;not null;unique;comment:用户名"`
	RealName  string     `json:"realname" gorm:"size:100;comment:真实姓名"`
	Password  string     `json:"-" gorm:"size:100;not null;default:123456;comment:密码"`
	Gender    GenderType `json:"gender" gorm:"size:10;default:'unknown';not null;comment:性别"`
	Phone     string     `json:"phone" gorm:"size:20;comment:手机号码"`
	Avatar    string     `json:"avatar" gorm:"size:255;comment:头像"`
	Email     string     `json:"email" gorm:"size:100;comment:邮箱"`
	LastLogin time.Time  `json:"lastLogin" gorm:"index:idx_users_last_login;comment:上次登录时间"`
	IsActive  bool       `json:"isActive" gorm:"default:true;not null;comment:账号是否活跃"`
	Role      string     `json:"role" gorm:"size:20;not null;default:user;comment:角色-admin管理员user普通用户..."`
	Remark    string     `json:"remark" gorm:"type:text;comment:备注"`
}
