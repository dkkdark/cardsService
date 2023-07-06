package repository

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type ServiceImpl struct {
	db           *gorm.DB
	dbConnection string
	dbType       string
}

type InitParams struct {
	DBConnection string
	DBType       string
}

func New(params *InitParams) (*ServiceImpl, error) {
	serviceImpl := &ServiceImpl{}
	db, err := gorm.Open(params.DBType, params.DBConnection)
	if err != nil {
		return nil, fmt.Errorf("error open %s connection, err: %w", params.DBType, err)
	}
	serviceImpl.db = db
	serviceImpl.dbConnection = params.DBConnection
	serviceImpl.dbType = params.DBType

	if err = serviceImpl.checkConnection(); err != nil {
		return nil, err
	}

	return serviceImpl, nil
}

func (s *ServiceImpl) pingDb() error {
	ping := Ping{}
	err := s.db.Raw("SELECT 1+1 as result").Scan(&ping).Error
	return err
}

func (s *ServiceImpl) checkConnection() error {
	err := s.pingDb()
	if err != nil {
		err := s.db.Close()
		if err != nil {
			return fmt.Errorf("can't close db connection in checkConnection method, err: %+v", err)
		}
		s.db, err = gorm.Open(s.dbType, s.dbConnection)
		if err != nil {
			return fmt.Errorf("can't re open db connection in checkConnection method, err: %+v", err)
		}
	}
	return nil
}
