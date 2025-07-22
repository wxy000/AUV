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
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	DB     *gorm.DB // 已导出为公共变量
	DBType string   // 数据库类型：mysql 或 sqlite
)

func InitDB(cfg *config.AppConfig) error {
	var dsn string
	var dialector gorm.Dialector

	switch cfg.DB.Driver {
	case "mysql":
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			cfg.DB.Mysql.User, cfg.DB.Mysql.Password, cfg.DB.Mysql.Host, cfg.DB.Mysql.Port, cfg.DB.Mysql.Name)
		dialector = mysql.Open(dsn)
	case "sqlite":
		dsn = cfg.DB.Sqlite.Path
		// SQLite默认连接参数
		if dsn == "" {
			dsn = "file:auv.db?cache=shared&_fk=1"
		}
		dialector = sqlite.Open(dsn)
	default:
		return fmt.Errorf("unsupported database driver: %s", cfg.DB.Driver)
	}

	var err error
	DB, err = gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to connect database: %w", err)
	}

	DBType = cfg.DB.Driver
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

	tableOptions := ""
	if DBType == "mysql" {
		tableOptions = "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4"
	}
	err := DB.Set("gorm:table_options", tableOptions).AutoMigrate(&models.User{})
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
