package login

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

const (
	EndpointRoot                          = "/"
	EndpointHealthz                       = "/healthz"
	EndpointReadiness                     = "/ready"
	EndpointLogin                         = "/login"
	EndpointExternalLogin                 = "/login/externalidp"
	EndpointExternalLoginCallback         = "/login/externalidp/callback"
	EndpointExternalLoginCallbackFormPost = "/login/externalidp/callback/form"
	EndpointExternalLogout                = "/logout/externalidp"
	EndpointSAMLACS                       = "/login/externalidp/saml/acs"
	EndpointJWTAuthorize                  = "/login/jwt/authorize"
	EndpointJWTCallback                   = "/login/jwt/callback"
	EndpointLDAPLogin                     = "/login/ldap"
	EndpointLDAPCallback                  = "/login/ldap/callback"
	EndpointPasswordlessLogin             = "/login/passwordless"
	EndpointPasswordlessRegistration      = "/login/passwordless/init"
	EndpointPasswordlessPrompt            = "/login/passwordless/prompt"
	EndpointLoginName                     = "/loginname"
	EndpointUserSelection                 = "/userselection"
	EndpointChangeUsername                = "/username/change"
	EndpointPassword                      = "/password"
	EndpointInitPassword                  = "/password/init"
	EndpointChangePassword                = "/password/change"
	EndpointPasswordReset                 = "/password/reset"
	EndpointInitUser                      = "/user/init"
	EndpointInviteUser                    = "/user/invite"
	EndpointMFAVerify                     = "/mfa/verify"
	EndpointMFAPrompt                     = "/mfa/prompt"
	EndpointMFAInitVerify                 = "/mfa/init/verify"
	EndpointMFASMSInitVerify              = "/mfa/init/sms/verify"
	EndpointMFAOTPVerify                  = "/mfa/otp/verify"
	EndpointMFAInitU2FVerify              = "/mfa/init/u2f/verify"
	EndpointU2FVerification               = "/mfa/u2f/verify"
	EndpointMailVerification              = "/mail/verification"
	EndpointMailVerified                  = "/mail/verified"
	EndpointRegisterOption                = "/register/option"
	EndpointRegister                      = "/register"
	EndpointExternalRegister              = "/register/externalidp"
	EndpointExternalRegisterCallback      = "/register/externalidp/callback"
	EndpointRegisterOrg                   = "/register/org"
	EndpointLogoutDone                    = "/logout/done"
	EndpointLoginSuccess                  = "/login/success"
	EndpointExternalNotFoundOption        = "/externaluser/option"

	EndpointResources        = "/resources"
	EndpointDynamicResources = "/resources/dynamic"

	EndpointDeviceAuth       = "/device"
	EndpointDeviceAuthAction = "/device/{action}"

	EndpointLinkingUserPrompt = "/link/user"
)

var (
	IgnoreInstanceEndpoints = []string{
		EndpointResources + "/fonts",
		EndpointResources + "/images",
		EndpointResources + "/scripts",
		EndpointResources + "/themes",
	}
)

func CreateRouter(login *Login, interceptors ...func(http.Handler) http.Handler) chi.Router {
	router := chi.NewRouter()
	router.Use(interceptors...)
	router.Get(EndpointRoot, login.handleLogin)
	router.Get(EndpointHealthz, login.handleHealthz)
	router.Get(EndpointReadiness, login.handleReadiness)
	router.Get(EndpointLogin, login.handleLogin)
	router.Post(EndpointLogin, login.handleLogin)
	router.Get(EndpointExternalLogin, login.handleExternalLogin)
	router.Get(EndpointExternalLoginCallback, login.handleExternalLoginCallback)
	router.Post(EndpointExternalLoginCallbackFormPost, login.handleExternalLoginCallbackForm)
	router.Get(EndpointExternalLogout, login.handleExternalLogout)
	router.Get(EndpointSAMLACS, login.handleExternalLoginCallback)
	router.Post(EndpointSAMLACS, login.handleExternalLoginCallbackForm)
	router.Get(EndpointJWTAuthorize, login.handleJWTRequest)
	router.Get(EndpointJWTCallback, login.handleJWTCallback)
	router.Post(EndpointPasswordlessLogin, login.handlePasswordlessVerification)
	router.Get(EndpointPasswordlessRegistration, login.handlePasswordlessRegistration)
	router.Post(EndpointPasswordlessRegistration, login.handlePasswordlessRegistrationCheck)
	router.Post(EndpointPasswordlessPrompt, login.handlePasswordlessPrompt)
	router.Get(EndpointLoginName, login.handleLoginName)
	router.Post(EndpointLoginName, login.handleLoginNameCheck)
	router.Post(EndpointUserSelection, login.handleSelectUser)
	router.Post(EndpointChangeUsername, login.handleChangeUsername)
	router.Post(EndpointPassword, login.handlePasswordCheck)
	router.Get(EndpointInitPassword, login.handleInitPassword)
	router.Post(EndpointInitPassword, login.handleInitPasswordCheck)
	router.Get(EndpointPasswordReset, login.handlePasswordReset)
	router.Get(EndpointInitUser, login.handleInitUser)
	router.Post(EndpointInitUser, login.handleInitUserCheck)
	router.Get(EndpointInviteUser, login.handleInviteUser)
	router.Post(EndpointInviteUser, login.handleInviteUserCheck)
	router.Post(EndpointMFAVerify, login.handleMFAVerify)
	router.Get(EndpointMFAPrompt, login.handleMFAPromptSelection)
	router.Post(EndpointMFAPrompt, login.handleMFAPrompt)
	router.Post(EndpointMFAInitVerify, login.handleMFAInitVerify)
	router.Post(EndpointMFASMSInitVerify, login.handleRegisterSMSCheck)
	router.Get(EndpointMFAOTPVerify, login.handleOTPVerificationCheck)
	router.Post(EndpointMFAOTPVerify, login.handleOTPVerificationCheck)
	router.Post(EndpointMFAInitU2FVerify, login.handleRegisterU2F)
	router.Post(EndpointU2FVerification, login.handleU2FVerification)
	router.Get(EndpointMailVerification, login.handleMailVerification)
	router.Post(EndpointMailVerification, login.handleMailVerificationCheck)
	router.Post(EndpointChangePassword, login.handleChangePassword)
	router.Get(EndpointRegisterOption, login.handleRegisterOption)
	router.Post(EndpointRegisterOption, login.handleRegisterOptionCheck)
	router.Post(EndpointExternalNotFoundOption, login.handleExternalNotFoundOptionCheck)
	router.Get(EndpointRegister, login.handleRegister)
	router.Post(EndpointRegister, login.handleRegisterCheck)
	router.Get(EndpointExternalRegister, login.handleExternalRegister)
	router.Get(EndpointExternalRegisterCallback, login.handleExternalLoginCallback)
	router.Get(EndpointLogoutDone, login.handleLogoutDone)
	router.Get(EndpointDynamicResources, login.handleDynamicResources)
	router.PathPrefix(EndpointResources).Handler(login.handleResources()).Methods(http.MethodGet)
	router.Get(EndpointRegisterOrg, login.handleRegisterOrg)
	router.Post(EndpointRegisterOrg, login.handleRegisterOrgCheck)
	router.Get(EndpointLoginSuccess, login.handleLoginSuccess)
	router.Get(EndpointLDAPLogin, login.handleLDAP)
	router.Post(EndpointLDAPCallback, login.handleLDAPCallback)
	router.SkipClean(true).Handle("", http.RedirectHandler(HandlerPrefix+"/", http.StatusMovedPermanently))
	router.Get(EndpointDeviceAuth, login.handleDeviceAuthUserCode)
	router.Post(EndpointDeviceAuth, login.handleDeviceAuthUserCode)
	router.Get(EndpointDeviceAuthAction, login.handleDeviceAuthAction)
	router.Post(EndpointDeviceAuthAction, login.handleDeviceAuthAction)
	return router
}
