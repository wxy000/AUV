package db

import (
	"AUV/config"
	"AUV/models"
	"errors"
	"fmt"
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB // 已导出为公共变量

func InitDB(cfg *config.AppConfig) error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.DB.User, cfg.DB.Password, cfg.DB.Host, cfg.DB.Port, cfg.DB.Name)

	var err error

	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to connect database: %w", err)
	}

	sqlDB, _ := DB.DB()

	// 设置连接池，空闲连接
	sqlDB.SetMaxIdleConns(50)
	// 打开链接
	sqlDB.SetMaxOpenConns(100)
	// SetConnMaxLifetime 设置了连接可复用的最大时间。
	sqlDB.SetConnMaxLifetime(time.Hour)

	log.Println("数据库连接成功")

	return nil
}

// AutoMigrate 自动迁移表结构
func AutoMigrate() {
	if DB == nil {
		log.Fatal("数据库连接未初始化，请先调用InitDB")
	}

	err := DB.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4").AutoMigrate(&models.User{})
	if err != nil {
		log.Fatalf("数据库迁移失败: %v", err)
	}

	// 初始化管理员账户
	initAdminAccount()
}

// 初始化管理员账户
func initAdminAccount() {
	adminUser := config.Cfg.Admin
	if adminUser.Username == "" || adminUser.Password == "" {
		log.Println("管理员配置不完整，跳过初始化")
		return
	}

	var existing models.User
	err := DB.Where("role = ?", "admin").First(&existing).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		// 管理员账户不存在，继续创建流程
	} else if err != nil {
		log.Fatalf("管理员账户查询失败: %v", err)
		return
	} else {
		log.Println("管理员账户已存在，跳过初始化")
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(adminUser.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("密码加密失败: %v", err)
	}

	tx := DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Create(&models.User{
		Username: adminUser.Username,
		Password: string(hash),
		Role:     "admin",
	}).Error; err != nil {
		tx.Rollback()
		log.Fatalf("管理员账户创建失败: %v", err)
	}

	tx.Commit()
	log.Println("默认管理员账户创建成功")
}
