package gorm

import (
	"errors"

	"github.com/dattito/purrmannplus-backend/config"
	db_errors "github.com/dattito/purrmannplus-backend/database/errors"
	provider_models "github.com/dattito/purrmannplus-backend/database/models"
	"github.com/dattito/purrmannplus-backend/database/providers/gorm/models"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type GormProvider struct {
	DB *gorm.DB
}

// Returns a GormProvider object with a connection to the database
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
		Logger: logger.Default.LogMode(logger.LogLevel(config.DATABASE_LOG_LEVEL)),
	})
	if err != nil {
		return &GormProvider{}, err
	}


	return &GormProvider{DB: db}, nil
}

// Creates all tables in the database using AutoMigrate()
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

// Closes the database connection
func (g *GormProvider) CloseDB() error {
	dialect, err := g.DB.DB()
	if err != nil {
		return err
	}
	defer dialect.Close()

	return nil
}

// Adds an account with it's credendials (username=authId, password=authPw) to the database
func (g *GormProvider) AddAccount(username, password string) (provider_models.AccountDBModel, error) {

	accdb := models.AccountDB{
		Username: username,
		Password: password,
	}
	err := g.DB.Create(&accdb).Error
	return models.AccountDBToAccountDBModel(accdb), err
}

// Returns account object of given accountId
func (g *GormProvider) GetAccount(id string) (provider_models.AccountDBModel, error) {

	accdb := models.AccountDB{}

	if err := g.DB.First(&accdb, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return provider_models.AccountDBModel{}, &db_errors.ErrRecordNotFound
		}
		return provider_models.AccountDBModel{}, err
	}

	return models.AccountDBToAccountDBModel(accdb), nil
}

// Gets account using username (authId) and password (authPw)
func (g *GormProvider) GetAccountByCredentials(username, password string) (provider_models.AccountDBModel, error) {

	accdb := models.AccountDB{}

	err := g.DB.First(&accdb, "auth_id = ? AND auth_pw = ?", username, password).Error

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
	return g.DB.Delete(a, "id = ?", id).Error
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
	err := g.DB.First(&accInfo, "account_id = ?", accountId).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return provider_models.AccountInfoDBModel{}, &db_errors.ErrRecordNotFound
		}
		return provider_models.AccountInfoDBModel{}, err
	}
	return models.AccountInfoDBToAccountInfoDBModel(accInfo), err
}

// Removes an account from the substitution_updater table if exists
func (g *GormProvider) RemoveAccountFromSubstitutionUpdater(accountId string) error {

	if err := g.DB.Delete(&models.SubstitutionDB{}, "account_id = ?", accountId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &db_errors.ErrRecordNotFound
		}
		return err
	}
	return nil
}

// Updates the substitution of a given account
func (g *GormProvider) SetSubstitutions(accountId string, entries map[string][]string, NotSetYet bool) (provider_models.SubstitutionDBModel, error) {

	var entriesE models.Entries = entries

	subdb := models.SubstitutionDB{
		AccountId: accountId,
	}

	if err := g.DB.FirstOrCreate(&subdb, "account_id = ?", accountId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return provider_models.SubstitutionDBModel{}, &db_errors.ErrRecordNotFound
		}
		return provider_models.SubstitutionDBModel{}, err
	}

	subdb.Entries = &entriesE
	subdb.NotSetYet = NotSetYet

	err := g.DB.Save(&subdb).Error
	if err != nil {
		return provider_models.SubstitutionDBModel{}, err
	}

	return models.SubstitutionDBToSubstitutionDBModel(subdb), nil
}

//Returns the substitution of a given account
func (g *GormProvider) GetSubstitutions(accountId string) (provider_models.SubstitutionDBModel, error) {
	subdb := models.SubstitutionDB{}

	err := g.DB.First(&subdb, "account_id = ?", accountId).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return provider_models.SubstitutionDBModel{}, &db_errors.ErrRecordNotFound
		}
		return provider_models.SubstitutionDBModel{}, err
	}

	return models.SubstitutionDBToSubstitutionDBModel(subdb), nil
}

// Returns all information needed to update the substitution of a given account in one list of models
func (g *GormProvider) GetAllAccountCredentialsAndPhoneNumberAndSubstitutions() ([]provider_models.AccountCredentialsAndPhoneNumberAndSubstitutionsDBModel, error) {
	m := []models.AccountCredentialsAndPhoneNumberAndSubstitutionsDB{}

	g.DB.Model(models.AccountDB{}).Select("accounts.auth_id", "accounts.auth_pw", "account_infos.phone_number", "accounts.id AS 'account_id'", "substitutions.id AS 'substitutions_id'", "substitutions.entries", "substitutions.not_set_yet").Joins("INNER JOIN substitutions ON substitutions.account_id = accounts.id").Joins("INNER JOIN account_infos ON account_infos.account_id = accounts.id").Scan(&m)

	var mm []provider_models.AccountCredentialsAndPhoneNumberAndSubstitutionsDBModel
	for _, v := range m {
		mm = append(mm, models.ACPSDBtoACPDSDBM(v))
	}

	return mm, nil
}

// Returns the accountId, auth_id, auth_pw, phone_number, substitutions_id and the substitutions of a given account
func (g *GormProvider) GetAccountCredentialsAndPhoneNumberAndSubstitutions(accountId string) (provider_models.AccountCredentialsAndPhoneNumberAndSubstitutionsDBModel, error) {
	m := models.AccountCredentialsAndPhoneNumberAndSubstitutionsDB{}
	g.DB.Model(models.AccountDB{}).Select("accounts.auth_id", "accounts.auth_pw", "account_infos.phone_number", "accounts.id AS 'account_id'", "substitutions.id AS 'substitutions_id'", "substitutions.entries", "substitutions.not_set_yet").Joins("INNER JOIN substitutions ON substitutions.account_id = accounts.id").Joins("INNER JOIN account_infos ON account_infos.account_id = accounts.id").Where("accounts.id = ?", accountId).Scan(&m)
	return models.ACPSDBtoACPDSDBM(m), nil
}
