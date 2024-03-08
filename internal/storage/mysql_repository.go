package storage

import (
	"github.com/arielcr/payment-gateway/internal/config"
	"github.com/arielcr/payment-gateway/internal/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type MySQLRepository struct {
	db     *gorm.DB
	config config.Application
}

func NewMySQLRepository(config config.Application) *MySQLRepository {
	return &MySQLRepository{
		config: config,
	}
}

func (m *MySQLRepository) Connect() error {
	var err error
	c := m.config.Repository
	dsn := c.User + ":" + c.Password + "@tcp" + "(" + c.Host + ":" + c.Port + ")/" + c.DBName + "?" + "parseTime=true&loc=Local"

	if m.db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{}); err != nil {
		return err
	}
	return nil
}

func (m *MySQLRepository) CreatePayment(payment models.Payment) error {
	if result := m.db.Create(&payment); result.Error != nil {
		return result.Error
	}
	return nil
}
