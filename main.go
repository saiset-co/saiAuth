package main

import (
	"encoding/json"
	"github.com/Limpid-LLC/go-auth/internal"
	"github.com/Limpid-LLC/go-auth/internal/entities"
	"github.com/Limpid-LLC/go-auth/internal/repo"
	"github.com/Limpid-LLC/go-auth/logger"
	"github.com/Limpid-LLC/saiService"
	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
	"github.com/saiset-co/sai-storage-mongo/external/adapter"
	"log"
	_ "net/http/pprof"
	"time"
)

func main() {
	name := "Auth"

	svc := saiService.NewService(name)

	svc.RegisterConfig("config.yml")

	logger.Logger = svc.Logger

	url := svc.GetConfig("common.storage.url", "").(string)
	token := svc.GetConfig("common.storage.token", "").(string)

	store := &adapter.SaiStorage{
		Url:   url,
		Token: token,
	}

	validate := validator.New()

	smsEnabled := svc.GetConfig("common.sms.enabled", false).(bool)
	smsUrl := svc.GetConfig("common.sms.url", "").(string)
	masterKey := svc.GetConfig("common.sms.master_key", "").(string)
	defaultRole := svc.GetConfig("default_role", "").(string)
	adminRole := svc.GetConfig("admin_role", "").(string)
	emailEnabled := svc.GetConfig("common.email.enabled", false).(bool)
	emailServiceUrl := svc.GetConfig("common.email.url", "").(string)
	emailSender := svc.GetConfig("common.email.sender", "").(string)

	var role entities.Role
	err := json.Unmarshal([]byte(defaultRole), &role)
	if err != nil {
		log.Fatalln(errors.Wrap(err, "Default role un-marshal error"))
	}

	var aRole entities.Role
	err = json.Unmarshal([]byte(adminRole), &aRole)
	if err != nil {
		log.Fatalln(errors.Wrap(err, "Admin role un-marshal error"))
	}

	salt := svc.GetConfig("common.encryption.salt", "").(string)

	if salt == "" {
		log.Fatalln(errors.Wrap(err, "Salt should be define in config"))
	}

	authUrl := svc.GetConfig("common.auth.url", "").(string)
	authFloodLimit := svc.GetConfig("common.auth.flood_limit", "").(int)
	authFloodDuration := svc.GetConfig("common.auth.flood_duration", "").(int)

	usersRepository := &repo.UsersRepository{
		Storage:    store,
		Collection: "users",
	}

	tokenPermissionsRepository := &repo.TokenPermissionsRepository{
		Storage:    store,
		Collection: "tokenPermissions",
	}

	is := internal.InternalService{
		Context: svc.Context,
		Storage: store,

		UsersRepository:            usersRepository,
		TokenPermissionsRepository: tokenPermissionsRepository,

		DefaultRole: role,
		AdminRole:   aRole,
		Validate:    validate,

		SmsEnabled:    smsEnabled,
		SmsServiceUrl: smsUrl,
		MasterKey:     masterKey,

		EmailEnabled:    emailEnabled,
		EmailServiceUrl: emailServiceUrl,
		EmailSender:     emailSender,

		Salt: salt,

		TokenExpirations: entities.TokenExpirations{
			AccessToken:  time.Duration(svc.GetConfig("tokens.expiration.access_token", 0).(int)),
			RefreshToken: time.Duration(svc.GetConfig("tokens.expiration.refresh_token", 0).(int)),
		},

		MasterToken: svc.GetConfig("tokens.token", "").(string),

		RoutineExecutionPeriods: entities.RoutineExecutionPeriods{
			Otp:          time.Duration(svc.GetConfig("tokens.routine_execution_period.otp", 0).(int)),
			AccessToken:  time.Duration(svc.GetConfig("tokens.routine_execution_period.access_token", 0).(int)),
			RefreshToken: time.Duration(svc.GetConfig("tokens.routine_execution_period.refresh_token", 0).(int)),
		},

		AuthUrl:           authUrl,
		AuthFloodLimit:    authFloodLimit,
		AuthFloodDuration: authFloodDuration,
		Name:              name,
	}

	svc.RegisterHandlers(
		is.NewHandler(),
	)

	svc.RegisterInitTask(is.Init)

	svc.Start()

}
