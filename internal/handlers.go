package internal

import (
	"net/http"

	"github.com/Limpid-LLC/saiService"
	"github.com/Limpid-LLC/saiService/middlewares"
)

func (is InternalService) NewHandler() saiService.Handler {
	return saiService.Handler{
		"send_verify_code": saiService.HandlerElement{
			Name:        "Send OTP code",
			Description: "Generates and sends an OTP code",
			Function:    is.sendVerifyCodeHandler,
		},
		"send_reset_password_verify_code": saiService.HandlerElement{
			Name:        "Send OTP code",
			Description: "Generates and sends an OTP code",
			Function:    is.sendResetPasswordVerifyCodeHandler,
		},
		"sign_up": saiService.HandlerElement{
			Name:        "Sign up",
			Description: "Registers a new user",
			Function:    is.signUpHandler,
		},
		"check": saiService.HandlerElement{
			Name:        "Check token validity for request",
			Description: "Checks token validity for request",
			Function:    is.checkHandler,
		},
		"sign_in": saiService.HandlerElement{
			Name:        "Login",
			Description: "Login user",
			Function:    is.signInHandler,
		},
		"update_user": saiService.HandlerElement{
			Name:        "Update user",
			Description: "Updates user information",
			Function:    is.updateUserHandler,
			Middlewares: []saiService.Middleware{
				middlewares.CreateAuthMiddleware(is.AuthUrl, is.Name, "update_user"),
			},
		},
		"get_users": saiService.HandlerElement{
			Name:        "Get users",
			Description: "Fetches user data",
			Function:    is.getUsersHandler,
			Middlewares: []saiService.Middleware{
				middlewares.CreateAuthMiddleware(is.AuthUrl, is.Name, "get_users"),
			},
		},
		"delete_users": saiService.HandlerElement{
			Name:        "Delete users",
			Description: "Deletes users",
			Function:    is.deleteUsersHandler,
			Middlewares: []saiService.Middleware{
				middlewares.CreateAuthMiddleware(is.AuthUrl, is.Name, "delete_users"),
			},
		},
		"create_role": saiService.HandlerElement{
			Name:        "Create role",
			Description: "Creates a new role",
			Function:    is.createRoleHandler,
			Middlewares: []saiService.Middleware{
				middlewares.CreateAuthMiddleware(is.AuthUrl, is.Name, "create_role"),
			},
		},
		"update_roles": saiService.HandlerElement{
			Name:        "Update role",
			Description: "Updates role",
			Function:    is.updateRolesHandler,
			Middlewares: []saiService.Middleware{
				middlewares.CreateAuthMiddleware(is.AuthUrl, is.Name, "update_roles"),
			},
		},
		"delete_roles": saiService.HandlerElement{
			Name:        "Delete role",
			Description: "Deletes role",
			Function:    is.deleteRolesHandler,
			Middlewares: []saiService.Middleware{
				middlewares.CreateAuthMiddleware(is.AuthUrl, is.Name, "delete_roles"),
			},
		},

		"attach_role": saiService.HandlerElement{
			Name:        "Attach role",
			Description: "Attaches role to user",
			Function:    is.attachRoleHandler,
			Middlewares: []saiService.Middleware{
				middlewares.CreateAuthMiddleware(is.AuthUrl, is.Name, "attach_role"),
			},
		},

		"detach_role": saiService.HandlerElement{
			Name:        "Detach role",
			Description: "Detaches role from user",
			Function:    is.detachRoleHandler,
			Middlewares: []saiService.Middleware{
				middlewares.CreateAuthMiddleware(is.AuthUrl, is.Name, "detach_role"),
			},
		},

		"get_roles": saiService.HandlerElement{
			Name:        "Get roles",
			Description: "Fetches roles",
			Function:    is.getRolesHandler,
			Middlewares: []saiService.Middleware{
				middlewares.CreateAuthMiddleware(is.AuthUrl, is.Name, "get_roles"),
			},
		},

		"test_cred": saiService.HandlerElement{
			Name:        "Test credentials",
			Description: "Tests credentials",
			Function: func(i interface{}, meta interface{}) (interface{}, int, error) {
				return NewOkResponse("OK")
			},
			Middlewares: []saiService.Middleware{
				middlewares.CreateAuthMiddleware(is.AuthUrl, is.Name, "test_cred"),
			},
		},

		"check_otp_code": saiService.HandlerElement{
			Name:        "Check otp code",
			Description: "Check otp code",
			Function:    is.checkOTPCodeHandler,
		},

		"restore_password": saiService.HandlerElement{
			Name:        "Restore password",
			Description: "Restore password",
			Function:    is.restorePasswordHandler,
		},
	}
}

type ResponseOk struct {
	Result interface{} `json:"result,omitempty"`
	Status string      `json:"status"`
}

func NewOkResponse(result interface{}) (ResponseOk, int, error) {
	return ResponseOk{
		Result: result,
		Status: "OK",
	}, http.StatusOK, nil
}

type ErrorResponse struct {
	ErrorType string `json:"ErrorType"`
	ErrorCode string `json:"ErrorCode"`
	Error     string `json:"Error"`
	Status    string `json:"Status"`
}

func NewErrorResponse(errorType string, errorCode string, errorText string) ErrorResponse {
	return ErrorResponse{
		ErrorType: errorType,
		ErrorCode: errorCode,
		Error:     errorText,
		Status:    "NOK",
	}
}
