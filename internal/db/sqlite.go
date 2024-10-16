package db

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"sync"
	"time"
)

type dbMp struct {
	mp map[string]*gorm.DB
	sync.RWMutex
}

var mp = dbMp{mp: map[string]*gorm.DB{}}

func newDB(config Config) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(config.Path), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	// 设置连接池参数
	sqlDB.SetMaxIdleConns(config.IdleConn)                                    // 设置最大空闲连接数
	sqlDB.SetMaxOpenConns(config.MaxConn)                                     // 设置最大连接数
	sqlDB.SetConnMaxLifetime(time.Duration(config.MaxLifeTime) * time.Minute) // 设置连接的最大生命周期

	return db, nil
}

func Get(dsn string) (*gorm.DB, error) {
	mp.RLock()
	db, ok := mp.mp[dsn]
	mp.RUnlock()
	if ok {
		return db, nil
	}
	mp.Lock()
	defer mp.Unlock()
	db, err := newDB(Config{
		Path:        dsn,
		IdleConn:    2,
		MaxConn:     100,
		MaxLifeTime: 30,
	})
	mp.mp[dsn] = db
	return db, err
}
