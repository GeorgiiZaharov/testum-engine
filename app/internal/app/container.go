package app

import (
	"go.uber.org/zap"

	"testum-engine/app/internal/adapter/db"
	ldapadap "testum-engine/app/internal/adapter/ldap"
	storageadap "testum-engine/app/internal/adapter/storage"

	"testum-engine/app/internal/config"

	answerserv "testum-engine/app/internal/service/core/answer"
	authserv "testum-engine/app/internal/service/core/auth"
	latexvalidationserv "testum-engine/app/internal/service/core/latexvalidator"
	resultserv "testum-engine/app/internal/service/core/result"
	validationserv "testum-engine/app/internal/service/core/validation"
)

type Container struct {
	Logger         *zap.Logger
	DB             *db.DB
	BasePictureURL string

	// adapters
	LDAP    *ldapadap.LdapAdapter
	Storage *storageadap.StorageAdapter

	// core services
	AnswerService          *answerserv.CheckService
	AuthService            *authserv.Service
	LatexValidationService *latexvalidationserv.Validator
	ResultService          *resultserv.CalculationService
	ValidationService      *validationserv.Parser
}

func NewContainer(cfg config.Config) (*Container, error) {
	// =====================
	// logger
	// =====================
	logger, err := zap.NewProduction()
	if err != nil {
		return nil, err
	}

	// =====================
	// database
	// =====================
	database, err := db.NewDB(db.DBOptions{
		Path: cfg.DB.Path,
	})
	if err != nil {
		return nil, err
	}

	if _, err := database.Exec("PRAGMA foreign_keys = ON"); err != nil {
		return nil, err
	}

	migCfg := db.MigrationConfig{
		Dir:     "./migrations",
		Dialect: "sqlite",
	}

	if err := db.RunMigrations(database.DB.DB, migCfg); err != nil {
		return nil, err
	}

	logger.Info("database migrations applied successfully")

	// =====================
	// adapters
	// =====================
	ldapAdapter := ldapadap.NewLdapAdapter(cfg.Ldap, logger)
	storageAdapter := storageadap.NewStorageAdapter(storageadap.OSFileSystem{})

	// =====================
	// core services
	// =====================
	answerService := answerserv.NewCheckService()
	authService := authserv.New(cfg.App.Secret, logger)
	latexValidationService := latexvalidationserv.New()
	resultService := resultserv.NewCalculationService()
	validationService := validationserv.NewParser(logger)

	return &Container{
		Logger: logger,
		DB:     database,

		BasePictureURL: cfg.App.PictureBaseURL,

		LDAP:    &ldapAdapter,
		Storage: storageAdapter,

		AnswerService:          &answerService,
		AuthService:            &authService,
		LatexValidationService: latexValidationService,
		ResultService:          &resultService,
		ValidationService:      validationService,
	}, nil
}
