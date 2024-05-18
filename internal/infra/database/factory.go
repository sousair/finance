package database

import (
	"context"
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type myLogger struct{}

func (l myLogger) Info(ctx context.Context, msg string, data ...interface{})  {}
func (l myLogger) Warn(ctx context.Context, msg string, data ...interface{})  {}
func (l myLogger) Error(ctx context.Context, msg string, data ...interface{}) {}
func (l myLogger) LogMode(level logger.LogLevel) logger.Interface             { return l }
func (l myLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	sql, _ := fc()
	if err == nil {
		fmt.Println(sql + "\n\n")
	}
}

func NewPostgres(connectionURL string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(connectionURL), &gorm.Config{
		Logger: myLogger{},
	})
	if err != nil {
		return nil, err
	}

	return db, nil
}
