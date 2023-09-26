package entities

import "time"

type AccessToken struct {
	Token     string  `json:"token"`
	Type      string  `json:"type"`
	RoleId    string  `json:"roleId"`
	StoID     *string `json:"sto_id,omitempty"`
	ExpiredAt int64   `json:"expired_at"`
}

type TokenPermission struct {
	Token                      string   `json:"token"`
	Type                       string   `json:"type"`
	UserID                     string   `json:"user_id"`
	StoID                      *string  `json:"sto_id,omitempty"`
	ExpiredAt                  int64    `json:"expired_at"`
	RoleInternalID             string   `json:"role_internal_id"`
	PermissionMicroservice     string   `json:"permission_microservice"`
	PermissionMethod           string   `json:"permission_method"`
	PermissionRequiredParams   []Params `json:"permission_required_params"`
	PermissionRestrictedParams []Params `json:"permission_restricted_params"`
}

func (tp TokenPermission) CreateAccessToken() AccessToken {
	return AccessToken{
		Token:     tp.Token,
		RoleId:    tp.RoleInternalID,
		Type:      tp.Type,
		StoID:     tp.StoID,
		ExpiredAt: tp.ExpiredAt,
	}
}

type Role struct {
	InternalID  string       `json:"internal_id"`
	Type        string       `json:"type" validate:"required"`
	StoID       *string      `json:"sto_id,omitempty" validate:"required"`
	Permissions []Permission `json:"permissions" validate:"required"`
	Data        interface{}  `json:"data"`
}

type Params struct {
	Param  string   `json:"param" validate:"required"`
	Values []string `json:"values" validate:"required"`
	All    bool     `json:"all" validate:"required"`
}

type Permission struct {
	Microservice     string   `json:"microservice" validate:"required"`
	Method           string   `json:"method" validate:"required"`
	RequiredParams   []Params `json:"required_params"`
	RestrictedParams []Params `json:"restricted_params"`
}

type TokenExpirations struct {
	RefreshToken time.Duration
	AccessToken  time.Duration
}

type RoutineExecutionPeriods struct {
	Otp          time.Duration
	RefreshToken time.Duration
	AccessToken  time.Duration
}
