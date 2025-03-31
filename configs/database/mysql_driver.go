package database

import (
	"fmt"
	"github.com/MarlonG1/delivery-backend/configs"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type MysqlDriver struct {
	Config *config.EnvConfig
}

func NewMysqlDriver(config *config.EnvConfig) *MysqlDriver {
	return &MysqlDriver{Config: config}
}

func (m *MysqlDriver) GetDSN() gorm.Dialector {
	return mysql.Open(m.GetStringConnection())
}

func (m *MysqlDriver) GetStringConnection() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=Local",
		m.Config.Database.User,
		m.Config.Database.Password,
		m.Config.Database.Host,
		m.Config.Database.Port,
		m.Config.Database.Name,
		m.Config.Database.Charset)
}

func (m *MysqlDriver) GetHost() string {
	return m.Config.Database.Host
}

func (m *MysqlDriver) GetDriverName() string {
	return "MySQL"
}
