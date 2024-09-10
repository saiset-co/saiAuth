package internal

import (
	"github.com/Limpid-LLC/go-auth/internal/entities"
	"github.com/Limpid-LLC/go-auth/internal/repo"
	"github.com/Limpid-LLC/go-auth/internal/storage"
	"github.com/Limpid-LLC/saiService"
	"github.com/go-playground/validator/v10"
)

type InternalService struct {
	Context *saiService.Context
	Storage *storage.SaiStorage

	UsersRepository            *repo.UsersRepository
	TokenPermissionsRepository *repo.TokenPermissionsRepository

	Collection  string
	DefaultRole entities.Role
	Validate    *validator.Validate

	SmsEnabled    bool
	SmsServiceUrl string

	MasterKey string

	EmailEnabled    bool
	EmailServiceUrl string
	EmailSender     string

	Salt string

	TokenExpirations entities.TokenExpirations
	MasterToken      string

	RoutineExecutionPeriods entities.RoutineExecutionPeriods

	AuthUrl           string
	AuthFloodLimit    int
	AuthFloodDuration int

	Name string
}

func (is InternalService) Init() {
	go startCleanupRoutine(is.Context.Context, is.RoutineExecutionPeriods.Otp, is.removeExpiredOtpCodes)
	go startCleanupRoutine(is.Context.Context, is.RoutineExecutionPeriods.RefreshToken, is.removeExpiredRefreshTokens)
	go startCleanupRoutine(is.Context.Context, is.RoutineExecutionPeriods.AccessToken, is.TokenPermissionsRepository.RemoveExpiredTokens)
	go is.FloodClear()
}
