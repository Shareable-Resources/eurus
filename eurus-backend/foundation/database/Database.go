package database

import (
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	"eurus-backend/foundation/crypto"
	"eurus-backend/foundation/log"

	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

type DatabaseType uint16

const (
	DatabasePostgres DatabaseType = iota
	DatabaseSqlite
)

type ReadOnlyDatabase struct {
	IP           string `json: "ip"`
	Port         int    `json:"port"`
	UserName     string `json:"userName"`
	Password     string `json:"password"`
	DBName       string `json:"dbName"`
	SchemaName   string `json:"schemaName"`
	aesKey       string
	IdleConns    int `json:"dbIdleConns"`
	MaxOpenConns int `json:"dbMaxOpenConns"`
	dbManager    *gorm.DB
	DbType       DatabaseType // default postgres;
	Logger       *logrus.Logger
}

type Database struct {
	ReadOnlyDatabase
}

func (me *ReadOnlyDatabase) SetAESKey(key string) {
	me.aesKey = key
}

func (me *ReadOnlyDatabase) decryptPassword() (string, error) {
	password, err := crypto.DecryptAESFromBase64(me.Password, me.aesKey)
	decryptedPasswordBytes, _ := base64.StdEncoding.DecodeString(password)
	password = string(decryptedPasswordBytes)
	if err != nil {
		return "", err
	}
	return password, nil
}

func (me *ReadOnlyDatabase) Validate() {
	var err error = nil
	if me.Logger == nil {
		me.Logger = log.GetLogger(log.Name.Root)
	}

	if me.DbType == DatabasePostgres {
		if me.IP == "" {
			err = errors.New("DBServerIP should be provided!")
			me.Logger.Error(err.Error())
			panic(err)
		} else if me.Port == 0 {
			err = errors.New("DBServerPort should be provided!")
			me.Logger.Error(err.Error())
			panic(err)

		} else if me.UserName == "" {
			err = errors.New("DBUserName should be provided!")
			me.Logger.Error(err.Error())
			panic(err)

		} else if me.Password == "" {
			err = errors.New("DBPassword should be provided!")
			me.Logger.Error(err.Error())
			panic(err)
		}
	} else {
		if me.DBName == "" {
			err = errors.New("DBName is sqlite file name, must be provided")
			me.Logger.Error(err.Error())
			panic(err)
		}
	}
}

func (me *ReadOnlyDatabase) GetConn() (*gorm.DB, error) {
	if me.dbManager == nil {
		if me.DbType == DatabasePostgres {
			if me.dbManager == nil {
				password, err := me.decryptPassword()
				if err != nil {
					me.Logger.Errorln("Unable to decrypt DB password: ", err.Error())
					return nil, err
				}
				dsn := fmt.Sprintf("host=%s port=%v user=%s dbname=%s sslmode=disable password=%s", me.IP, me.Port, me.UserName, me.DBName, password)
				newLogger := logger.New(
					me.Logger, // io writer
					logger.Config{
						SlowThreshold: time.Second, // Slow SQL threshold
						LogLevel:      logger.Info, // Log level
						Colorful:      false,       // Disable color
					},
				)

				config := gorm.Config{
					Logger: newLogger,
				}
				if me.SchemaName != "" {
					config.NamingStrategy = schema.NamingStrategy{
						TablePrefix: me.SchemaName + ".", // table prefix for schema
					}
				}

				me.dbManager, err = gorm.Open(postgres.Open(dsn), &config)
				if err != nil {
					me.Logger.Errorln("Connect DB failed: ", err.Error())
					return nil, err
				}
				dbo, err := me.dbManager.DB()
				if err != nil {
					me.Logger.Error("Getting DB object error: ", err.Error())
					return nil, err
				}
				if me.IdleConns == 0 {
					me.IdleConns = 10
				}
				if me.MaxOpenConns == 0 {
					me.MaxOpenConns = 100
				}
				dbo.SetMaxIdleConns(me.IdleConns)
				dbo.SetMaxOpenConns(me.MaxOpenConns)
			}

		} else {
			if me.dbManager == nil {
				newLogger := logger.New(
					me.Logger, // io writer
					logger.Config{
						SlowThreshold: time.Second, // Slow SQL threshold
						LogLevel:      logger.Warn, // Log level
						Colorful:      false,       // Disable color
					},
				)
				config := gorm.Config{
					DisableAutomaticPing: true,
					Logger:               newLogger,
				}

				var err error
				me.dbManager, err = gorm.Open(sqlite.Open(me.DBName), &config)
				if err != nil {
					log.GetLogger(log.Name.Root).Errorln("Connect DB failed: ", err.Error())
					return nil, err
				}
				dbo, err := me.dbManager.DB()
				if err != nil {
					me.Logger.Error("Getting DB object error: ", err.Error())
					return nil, err
				}
				if me.IdleConns == 0 {
					me.IdleConns = 5
				}
				if me.MaxOpenConns == 0 {
					me.MaxOpenConns = 10
				}
				dbo.SetMaxIdleConns(me.IdleConns)
				dbo.SetMaxOpenConns(me.MaxOpenConns)
			}
		}
	}
	return me.dbManager.Session(&gorm.Session{}), nil
}
