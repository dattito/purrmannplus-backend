package gorm

import (
	"errors"

	app_models "github.com/dattito/purrmannplus-backend/app/models"
	"github.com/dattito/purrmannplus-backend/config"
	db_errors "github.com/dattito/purrmannplus-backend/database/errors"
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

	err = g.DB.AutoMigrate(&models.MoodleUserAssignmentsDB{})
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
func (g *GormProvider) AddAccount(username, password string) (app_models.Account, error) {

	accdb := models.AccountDB{
		Username: username,
		Password: password,
	}
	err := g.DB.Create(&accdb).Error
	return accdb.ToAccount(), err
}

// Returns account object of given accountId
func (g *GormProvider) GetAccount(id string) (app_models.Account, error) {

	accdb := models.AccountDB{}

	if err := g.DB.First(&accdb, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return app_models.Account{}, &db_errors.ErrRecordNotFound
		}
		return app_models.Account{}, err
	}

	return accdb.ToAccount(), nil
}

// Gets account using username (authId) and password (authPw)
func (g *GormProvider) GetAccountByCredentials(username, password string) (app_models.Account, error) {

	accdb := models.AccountDB{}

	err := g.DB.First(&accdb, "auth_id = ? AND auth_pw = ?", username, password).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return app_models.Account{}, &db_errors.ErrRecordNotFound
		}
		return app_models.Account{}, err
	}

	return accdb.ToAccount(), nil
}

// what do think it does?
func (g *GormProvider) GetAccounts() ([]app_models.Account, error) {

	accdb := []models.AccountDB{}

	err := g.DB.Find(&accdb).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return []app_models.Account{}, &db_errors.ErrRecordNotFound
		}
		return []app_models.Account{}, err
	}

	var accs []app_models.Account
	for _, v := range accdb {
		accs = append(accs, v.ToAccount())
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
func (g *GormProvider) AddAccountInfo(accountId, phoneNumber string) (app_models.AccountInfo, error) {

	accInfo := models.AccountInfoDB{
		AccountId:   accountId,
		PhoneNumber: phoneNumber,
	}
	err := g.DB.Create(&accInfo).Error
	return accInfo.ToAccountInfo(), err
}

// Gets the account_info from a given accountId
func (g *GormProvider) GetAccountInfo(accountId string) (app_models.AccountInfo, error) {
	accInfo := models.AccountInfoDB{}
	err := g.DB.First(&accInfo, "account_id = ?", accountId).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return app_models.AccountInfo{}, &db_errors.ErrRecordNotFound
		}
		return app_models.AccountInfo{}, err
	}
	return accInfo.ToAccountInfo(), err
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

func (g *GormProvider) AddAccountToSubstitution(accountId, authId, authPw string) error {

	if _, err := g.GetSubstitutions(accountId); err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	substitution := models.SubstitutionDB{
		AccountId: accountId,
		AuthId:    authId,
		AuthPw:    authPw,
		Entries:   &models.Entries{},
		NotSetYet: true,
	}

	return g.DB.Create(&substitution).Error
}

// Updates the substitution of a given account
func (g *GormProvider) SetSubstitutions(accountId string, entries map[string][]string, NotSetYet bool) error {

	var entriesE models.Entries = entries

	subdb := models.SubstitutionDB{
		AccountId: accountId,
	}

	if err := g.DB.FirstOrCreate(&subdb, "account_id = ?", accountId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &db_errors.ErrRecordNotFound
		}
		return err
	}

	subdb.Entries = &entriesE
	subdb.NotSetYet = NotSetYet

	err := g.DB.Save(&subdb).Error
	if err != nil {
		return err
	}

	return nil
}

//Returns the substitution of a given account
func (g *GormProvider) GetSubstitutions(accountId string) (app_models.Substitutions, error) {
	subdb := models.SubstitutionDB{}

	err := g.DB.First(&subdb, "account_id = ?", accountId).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return app_models.Substitutions{}, &db_errors.ErrRecordNotFound
		}
		return app_models.Substitutions{}, err
	}

	return subdb.ToSubstitutions(), nil
}

// Returns all information needed to update the substitution of a given account in one list of models
func (g *GormProvider) GetAllSubstitutionInfos() ([]app_models.SubstitutionInfo, error) {
	m := []models.SubstitutionInfoDB{}

	g.DB.Model(models.AccountDB{}).Select("account_infos.phone_number", "accounts.id AS 'account_id'", "substitutions.auth_id", "substitutions.auth_pw", "substitutions.id AS 'substitutions_id'", "substitutions.entries", "substitutions.not_set_yet").Joins("INNER JOIN substitutions ON substitutions.account_id = accounts.id").Joins("INNER JOIN account_infos ON account_infos.account_id = accounts.id").Scan(&m)

	var mm []app_models.SubstitutionInfo
	for _, v := range m {
		mm = append(mm, v.ToSubstitutionInfo())
	}

	return mm, nil
}

// Returns the accountId, auth_id, auth_pw, phone_number, substitutions_id and the substitutions of a given account
func (g *GormProvider) GetSubstitutionInfos(accountId string) (app_models.SubstitutionInfo, error) {
	m := models.SubstitutionInfoDB{}
	err := g.DB.Model(models.AccountDB{}).Select("account_infos.phone_number", "accounts.id AS 'account_id'", "substitutions.auth_id", "substitutions.auth_pw", "substitutions.id AS 'substitutions_id'", "substitutions.entries", "substitutions.not_set_yet").Joins("INNER JOIN substitutions ON substitutions.account_id = accounts.id").Joins("INNER JOIN account_infos ON account_infos.account_id = accounts.id").Where("accounts.id = ?", accountId).Scan(&m).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return app_models.SubstitutionInfo{}, &db_errors.ErrRecordNotFound
		}
		return app_models.SubstitutionInfo{}, err
	}
	return m.ToSubstitutionInfo(), nil
}

func (g *GormProvider) SetMoodleAssignments(accountId string, assignmentIds []int, notSetYet bool) error {

	var assignmentIdsE models.AssignmentIds = assignmentIds

	m := models.MoodleUserAssignmentsDB{
		AccountId: accountId,
	}

	if err := g.DB.FirstOrCreate(&m, "account_id = ?", accountId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &db_errors.ErrRecordNotFound
		}
		return err
	}

	m.AssignmentIds = &assignmentIdsE
	m.NotSetYet = notSetYet

	err := g.DB.Save(&m).Error
	if err != nil {
		return err
	}

	return nil
}

func (g *GormProvider) RemoveAccountFromMoodleAssignmentUpdater(accountId string) error {

	if err := g.DB.Delete(&models.MoodleUserAssignmentsDB{}, "account_id = ?", accountId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &db_errors.ErrRecordNotFound
		}
		return err
	}
	return nil
}

func (g *GormProvider) GetMoodleAssignments(accountId string) (app_models.MoodleAssignments, error) {
	m := models.MoodleUserAssignmentsDB{}

	err := g.DB.First(&m, "account_id = ?", accountId).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return app_models.MoodleAssignments{}, &db_errors.ErrRecordNotFound
		}
		return app_models.MoodleAssignments{}, err
	}

	return m.ToMoodleAssignments(), nil
}

func (g *GormProvider) GetAllMoodleAssignmentInfos() ([]app_models.MoodleAssignmentInfo, error) {
	m := []models.MoodleAssignmentInfoDB{}

	g.DB.Model(models.AccountDB{}).Select("accounts.auth_id", "accounts.auth_pw", "account_infos.phone_number", "accounts.id AS 'account_id'", "moodle_user_assignments.assignment_ids", "moodle_user_assignments.not_set_yet").Joins("INNER JOIN moodle_user_assignments ON moodle_user_assignments.account_id = accounts.id").Joins("INNER JOIN account_infos ON account_infos.account_id = accounts.id").Scan(&m)

	var mm []app_models.MoodleAssignmentInfo
	for _, v := range m {
		mm = append(mm, v.ToMoodleAssignmentInfo())
	}

	return mm, nil
}

func (g *GormProvider) GetMoodleAssignmentInfos(accountId string) (app_models.MoodleAssignmentInfo, error) {
	m := models.MoodleAssignmentInfoDB{}
	err := g.DB.Model(models.AccountDB{}).Select("accounts.auth_id", "accounts.auth_pw", "account_infos.phone_number", "accounts.id AS 'account_id'", "moodle_user_assignments.assignment_ids", "moodle_user_assignments.not_set_yet").Joins("INNER JOIN moodle_user_assignments ON moodle_user_assignments.account_id = accounts.id").Joins("INNER JOIN account_infos ON account_infos.account_id = accounts.id").Where("accounts.id = ?", accountId).Scan(&m).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return app_models.MoodleAssignmentInfo{}, &db_errors.ErrRecordNotFound
		}
		return app_models.MoodleAssignmentInfo{}, err
	}
	return m.ToMoodleAssignmentInfo(), nil
}
