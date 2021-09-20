package gorm

import (
	"errors"

	"github.com/datti-to/purrmannplus-backend/config"
	db_errors "github.com/datti-to/purrmannplus-backend/database/errors"
	provider_models "github.com/datti-to/purrmannplus-backend/database/models"
	"github.com/datti-to/purrmannplus-backend/database/providers/gorm/models"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
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

	db, err := gorm.Open(o(config.DATABASE_URI), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return &GormProvider{}, err
	}

	return &GormProvider{DB: db}, nil
}

func (g *GormProvider) CreateTables() error {
	var err error
	err = g.DB.AutoMigrate(&models.AccountDB{})
	if err != nil {
		return err
	}

	err = g.DB.AutoMigrate(&models.AccountInfoDB{})
	if err != nil {
		return err
	}

	err = g.DB.AutoMigrate(&models.SubstitutionDB{})
	if err != nil {
		return err
	}
	return nil
}

func (g *GormProvider) CloseDB() error {
	dialect, err := g.DB.DB()
	if err != nil {
		return err
	}
	defer dialect.Close()

	return nil
}

// Adds an account with it's credendials (username=authId, password=authPw) to the database
func (g *GormProvider) AddAccount(authId, authPw string) (provider_models.AccountDBModel, error) {

	accdb := models.AccountDB{
		AuthId: authId,
		AuthPw: authPw,
	}
	err := g.DB.Create(&accdb).Error
	return models.AccountDBToAccountDBModel(accdb), err
}

// Returns account object of given accountId
func (g *GormProvider) GetAccount(id string) (provider_models.AccountDBModel, error) {

	accdb := models.AccountDB{}

	if err := g.DB.Where("id = ?", id).First(&accdb).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return provider_models.AccountDBModel{}, &db_errors.ErrRecordNotFound
		}
		return provider_models.AccountDBModel{}, err
	}

	return models.AccountDBToAccountDBModel(accdb), nil
}

// Gets account using username (authId) and password (authPw)
func (g *GormProvider) GetAccountByCredentials(authId, authPw string) (provider_models.AccountDBModel, error) {

	accdb := models.AccountDB{}

	err := g.DB.Where("auth_id = ? AND auth_pw = ?", authId, authPw).First(&accdb).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return provider_models.AccountDBModel{}, &db_errors.ErrRecordNotFound
		}
		return provider_models.AccountDBModel{}, err
	}

	return models.AccountDBToAccountDBModel(accdb), nil
}

// what do think it does?
func (g *GormProvider) GetAccounts() ([]provider_models.AccountDBModel, error) {

	accdb := []models.AccountDB{}

	err := g.DB.Find(&accdb).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return []provider_models.AccountDBModel{}, &db_errors.ErrRecordNotFound
		}
		return []provider_models.AccountDBModel{}, err
	}

	var accs []provider_models.AccountDBModel
	for _, v := range accdb {
		accs = append(accs, models.AccountDBToAccountDBModel(v))
	}

	return accs, nil
}

// Does what you think it does...
func (g *GormProvider) DeleteAccount(id string) error {
	a := &models.AccountDB{}
	a.Id = id
	return g.DB.Where("id = ?", id).Delete(a).Error
}

// Adds a new account_info entry of an account to the database
func (g *GormProvider) AddAccountInfo(accountId, phoneNumber string) (provider_models.AccountInfoDBModel, error) {

	accInfo := models.AccountInfoDB{
		AccountId:   accountId,
		PhoneNumber: phoneNumber,
	}
	err := g.DB.Create(&accInfo).Error
	return models.AccountInfoDBToAccountInfoDBModel(accInfo), err
}

// Gets the account_info from a given accountId
func (g *GormProvider) GetAccountInfo(accountId string) (provider_models.AccountInfoDBModel, error) {
	accInfo := models.AccountInfoDB{}
	err := g.DB.Where("account_id = ?", accountId).First(&accInfo).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return provider_models.AccountInfoDBModel{}, &db_errors.ErrRecordNotFound
		}
		return provider_models.AccountInfoDBModel{}, err
	}
	return models.AccountInfoDBToAccountInfoDBModel(accInfo), err
}

// Adds an account to the substitution_updater table if not exists
func (g *GormProvider) AddAccountToSubstitutionUpdater(accountId string) error {
	s := models.SubstitutionDB{
		AccountId: accountId,
	}

	return g.DB.FirstOrCreate(&s).Error
}

// Removes an account from the substitution_updater table if exists
func (g *GormProvider) RemoveAccountFromSubstitutionUpdater(accountId string) error {

	if err := g.DB.Where("account_id = ?", accountId).Delete(&models.SubstitutionDB{}).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &db_errors.ErrRecordNotFound
		}
		return err
	}
	return nil
}

// Updates the substitution of a given account
func (g *GormProvider) SetSubstitutions(accountId string, entries map[string][]string) (provider_models.SubstitutionDBModel, error) {

	subdb := models.SubstitutionDB{}

	if err := g.DB.Where("account_id = ?", accountId).First(&subdb).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return provider_models.SubstitutionDBModel{}, &db_errors.ErrRecordNotFound
		}
		return provider_models.SubstitutionDBModel{}, err
	}

	var e models.Entries = entries

	subdb.Entries = &e

	return models.SubstitutionDBToSubstitutionDBModel(subdb), g.DB.Save(&subdb).Error
}

//Returns the substitution of a given account
func (g *GormProvider) GetSubstitutions(accountId string) (provider_models.SubstitutionDBModel, error) {
	subdb := models.SubstitutionDB{}

	err := g.DB.Where("account_d = ?", accountId).First(&subdb).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return provider_models.SubstitutionDBModel{}, &db_errors.ErrRecordNotFound
		}
		return provider_models.SubstitutionDBModel{}, err
	}

	return models.SubstitutionDBToSubstitutionDBModel(subdb), nil
}
