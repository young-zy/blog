package databases

import (
	"fmt"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"blog/conf"
)

var db *gorm.DB

func init() {
	var err error
	config := conf.Config.Database
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=true", config.Username, config.Password, config.Address, config.Port, config.DBName, config.Charset)
	log.Printf("connecting %s", dsn)
	//dsn := "user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"

	myLogger := NewMyLogger(log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second, // 慢 SQL 阈值
			LogLevel:      logger.Warn, // Log level
			Colorful:      true,        // 彩色打印
		},
	)

	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: myLogger,
	})

	if err != nil {
		log.Fatal(err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal(err)
	}
	sqlDB.SetMaxOpenConns(30)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(time.Hour)
}

func GetTransaction() *gorm.DB {
	return db.Begin()
}

func Close() {
	sqlDB, err := db.DB()
	if err != nil {
		panic("failed to close db")
	}
	err = sqlDB.Close()
	if err != nil {
		panic("failed to close db")
	}
	log.Println("db connection closed")
}
