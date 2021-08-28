package gorm

import (
	"errors"

	app_models "github.com/datti-to/purrmannplus-backend/app/models"
	"github.com/datti-to/purrmannplus-backend/config"
	db_models "github.com/datti-to/purrmannplus-backend/database/provider/gorm/models"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type GormProvider struct {
	DB *gorm.DB
}

func NewGormProvider() (*GormProvider, error) {

	type Open func(string) gorm.Dialector
	var o Open

	switch config.DATABASE_TYPE {
	case "POSTGRES":
		o = postgres.Open
	case "MYSQL":
		o = mysql.Open
	case "SQLITE":
		o = sqlite.Open
	default:
		return &GormProvider{}, errors.New("DATABASE_TYPE env has to one of ('POSTGRES', 'MYSQL', 'SQLITE')")
	}

	db, err := gorm.Open(o(config.DATABASE_URI), &gorm.Config{})
	if err != nil {
		return &GormProvider{}, err
	}

	return &GormProvider{DB: db}, nil
}

func (g *GormProvider) CloseDB() error {
	dialect, err := g.DB.DB()
	if err != nil {
		return err
	}
	defer dialect.Close()

	return nil
}

func (g *GormProvider) AddAccount(account *app_models.Account) error {

	accdb := db_models.AccountToAccountDB(*account)
	err := g.DB.Create(&accdb).Error
	if err != nil {
		return err
	}
	account.ID = accdb.ID
	return nil
}

func (g *GormProvider) GetAccount(id string) (app_models.Account, error) {

	accdb := db_models.AccountDB{}

	err := g.DB.First(&accdb, id).Error

	if err != nil {
		return app_models.Account{}, err
	}

	return db_models.AccountDBToAccount(accdb), nil
}

func (g *GormProvider) GetAccounts() ([]app_models.Account, error) {

	accdb := []db_models.AccountDB{}

	err := g.DB.Find(&accdb).Error

	if err != nil {
		return []app_models.Account{}, err
	}

	var accounts []app_models.Account
	for _, a := range accdb {
		accounts = append(accounts, db_models.AccountDBToAccount(a))
	}

	return accounts, nil
}

func (g *GormProvider) UpdateAccount(account app_models.Account) error {

	accdb := db_models.AccountToAccountDB(account)

	return g.DB.Save(accdb).Error
}

func (g *GormProvider) DeleteAccount(id string) error {

	accdb := db_models.AccountDB{}

	err := g.DB.First(&accdb, id).Error

	if err != nil {
		return err
	}

	return g.DB.Delete(&accdb).Error
}
