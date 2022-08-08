package serviceconfig

const (
	serviceConfigPath          = "service_config_path"
	defaultConfigPath          = "service"
	localhostConfigPath        = "localhost_service_config_path"
	localhostDefaultConfigPath = "localhost_service"
)

const (
	envServiceName            = "SERVICE_NAME"
	envEnv                    = "ENV"
	envPort                   = "PORT"
	envRedisHost              = "REDIS_HOST"
	envRedisPort              = "REDIS_PORT"
	envLoginPrivateKey        = "LOGIN_PRIVATE_KEY"
	envLoginPublicKey         = "LOGIN_PUBLIC_KEY"
	envAccessTokenExpiration  = "ACCESS_TOKEN_EXPIRATION_MILLIS"
	envRefreshTokenExpiration = "REFRESH_TOKEN_EXPIRATION_MILLIS"
	envDBHost                 = "DB_HOST"
	envDBPort                 = "DB_PORT"
	envDBUsername             = "DB_USERNAME"
	envDBPassword             = "DB_PASSWORD"
)
