package constant

type CurrencyActionType string

const (
	AdminUserStatusDisable = 0
	AdminUserStatusEnable  = 1

	UserStatusDisable = 0
	UserStatusEnable  = 1

	AdminAccessToken = "admin_access_token"
	ApiAccessToken   = "api_access_token"
)
