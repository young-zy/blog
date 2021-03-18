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

var DefaultDb *gorm.DB

var Default *Transaction

type Transaction struct {
	tx *gorm.DB
}

func (tx *Transaction) Rollback() {
	tx.tx.Rollback()
}

func (tx *Transaction) Commit() {
	tx.tx.Commit()
}

func init() {
	var err error
	config := conf.Config.Database
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=true", config.Username, config.Password, config.Address, config.Port, config.DBName, config.Charset)
	//dsn := "user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"

	myLogger := NewMyLogger(log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second, // 慢 SQL 阈值
			LogLevel:      logger.Warn, // Log level
			Colorful:      true,        // 彩色打印
		},
	)

	DefaultDb, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: myLogger,
	})

	if err != nil {
		log.Fatal(err)
	}

	Default = &Transaction{tx: DefaultDb}

	sqlDB, err := DefaultDb.DB()
	if err != nil {
		log.Fatal(err)
	}
	sqlDB.SetMaxOpenConns(30)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(time.Hour)
}

func GetTransaction() *Transaction {
	return &Transaction{tx: DefaultDb.Begin()}
}

func Close() {
	sqlDB, err := DefaultDb.DB()
	if err != nil {
		panic("failed to close db")
	}
	err = sqlDB.Close()
	if err != nil {
		panic("failed to close db")
	}
	log.Println("db connection closed")
}
