package databases

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/young-zy/blog/conf"
)

// DefaultDb the default db instance the project uses
var DefaultDb *gorm.DB

// Default the default transaction exposed when no transaction is needed
var Default *Transaction

// Transaction wrapper of a transaction instance
type Transaction struct {
	*gorm.DB
}

// // Rollback rolls back the transaction
// func (tx Transaction) Rollback() {
//
// }
//
// // Commit commits the transaction
// func (tx *Transaction) Commit() {
// 	tx.Commit()
// }

func init() {
	var err error
	config := conf.Config.Database
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=true", config.Username, config.Password, config.Address, config.Port, config.DBName, config.Charset)
	// dsn := "user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"

	myLogger := newMyLogger(log.New(os.Stdout, "\r\n", log.LstdFlags),
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

	Default = &Transaction{DB: DefaultDb}

	sqlDB, err := DefaultDb.DB()
	if err != nil {
		log.Fatal(err)
	}
	sqlDB.SetMaxOpenConns(30)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(time.Hour)
}

// GetTransaction creates a transaction object
func GetTransaction() *Transaction {
	return &Transaction{DB: DefaultDb.Begin()}
}

// Close closes the database connection
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
