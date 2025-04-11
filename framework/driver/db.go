package driver

import (
	"net/url"

	"github.com/charmbracelet/log"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DB struct {
	*gorm.DB
}

func NewEmptyDB() *DB {
	return &DB{}
}

func NewDB(dsn string) *DB {
	log.Info("dsn: " + dsn)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Error),
	})
	if err != nil {
		panic(err)
	}
	return &DB{
		DB: db,
	}
}

func NewPostgresDB(dsn string) *DB {
	db, err := gorm.Open(postgres.New(postgres.Config{DSN: dsn, PreferSimpleProtocol: true}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		panic(err)
	}
	return &DB{
		DB: db,
	}
}

type MysqlConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DbName   string
	Charset  string
	Timezone string
}

func NewMysqlConfig(viper *viper.Viper) *MysqlConfig {
	timezone := viper.GetString("mysql.timezone")
	if timezone == "" {
		timezone = "Local"
	}

	// trasform to url
	timezone = url.QueryEscape(timezone)

	return &MysqlConfig{
		Host:     viper.GetString("mysql.host"),
		Port:     viper.GetString("mysql.port"),
		User:     viper.GetString("mysql.user"),
		Password: viper.GetString("mysql.password"),
		DbName:   viper.GetString("mysql.dbname"),
		Charset:  viper.GetString("mysql.charset"),
		Timezone: timezone,
	}
}

func (m *MysqlConfig) GetDSN() string {
	return m.User + ":" + m.Password + "@tcp(" + m.Host + ":" + m.Port + ")/" + m.DbName + "?charset=" + m.Charset + "&parseTime=true&loc=" + m.Timezone
}
