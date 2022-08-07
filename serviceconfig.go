package serviceconfig

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

// ServiceConfig is the structure to hold all the configs related to a common golang service
type ServiceConfig struct {
	data serviceConfig
}

// not allow accessing fields directly by not exported serviceConfig
type serviceConfig struct {
	Database    DatabaseConfig `json:"database" yaml:"database"`
	HTTP        HTTPServer     `json:"http" yaml:"http"`
	Redis       RedisConfig    `json:"redis" yaml:"redis"`
	Auth        AuthConfig     `json:"auth" yaml:"auth"`
	MailService MailServiceKey `json:"mail_service" yaml:"mail_service"`
	Env         string         `json:"env" yaml:"env"`
	Logger      LoggerConfig   `json:"logger" yaml:"logger"`
}

// DatabaseConfig is a structure to store all the database config
type DatabaseConfig struct {
	HostAddress string `json:"host_address" yaml:"host_address"`
	Port        string `json:"port" yaml:"port"`
	Name        string `json:"name" yaml:"name"`
	User        string `json:"user" yaml:"user"`
	Password    string `json:"password" yaml:"password"`
}

// HTTPServer include host and port when run a service
type HTTPServer struct {
	BaseURL string `json:"base_url" yaml:"base_url"`
	Host    string `json:"host" yaml:"host"`
	Port    string `json:"port" yaml:"port"`
}

// RedisConfig include host and port to a single redis instance usage
type RedisConfig struct {
	Host string `json:"host" yaml:"host"`
	Port string `json:"port" yaml:"port"`
}

// AuthConfig include config for an authentication service
type AuthConfig struct {
	PrivateKey                   string `json:"private_key" yaml:"private_key"`
	PublicKey                    string `json:"public_key" yaml:"public_key"`
	AccessTokenExpirationMillis  uint32 `json:"access_token_expiration_millis" yaml:"access_token_expiration_millis"`
	RefreshTokenExpirationMillis uint32 `json:"refresh_token_expiration_millis" yaml:"refresh_token_expiration_millis"`
	CleanUpSessionMillis         uint32 `json:"clean_up_session_millis" yaml:"clean_up_session_millis"`
}

type LoggerConfig struct {
	ConfigPath string `json:"config_path" yaml:"config_path"`
}

// MailServiceKey include related information if want to use a provided mail service
type MailServiceKey struct {
	ApiKey string `json:"api_key" yaml:"api_key"`
}

// ServiceConfigInterface is an interface to export all the methods of the ServiceConfig struct
type ServiceConfigInterface interface {
	GetDataBaseConfig() DatabaseConfig
	GetAuthConfig() AuthConfig
	GetRedisConfig() RedisConfig
	GetMailServiceConfig() MailServiceKey
	GetServerPort() string
	GetEnv() string
	GetLoggerConfigPath() string
}

var (
	globalServiceConfig  *ServiceConfig
	defaultServiceConfig ServiceConfigInterface
)

func (s *ServiceConfig) GetDataBaseConfig() DatabaseConfig {
	return s.data.Database
}

func (s *ServiceConfig) GetAuthConfig() AuthConfig {
	return s.data.Auth
}

func (s *ServiceConfig) GetRedisConfig() RedisConfig {
	return s.data.Redis
}

func (s *ServiceConfig) GetMailServiceConfig() MailServiceKey {
	return s.data.MailService
}

func (s *ServiceConfig) GetServerPort() string {
	return s.data.HTTP.Port
}

func (s *ServiceConfig) GetEnv() string {
	return s.data.Env
}

func (s *ServiceConfig) GetLoggerConfigPath() string {
	return s.data.Logger.ConfigPath
}

func init() {
	globalServiceConfig = &ServiceConfig{}
	initServiceDefaultConfigValue(globalServiceConfig)
	defaultServiceConfig = globalServiceConfig
	initServiceConfig(globalServiceConfig)
	defaultServiceConfig = globalServiceConfig
}

func initServiceDefaultConfigValue(s *ServiceConfig) {}

func initServiceConfig(s *ServiceConfig) {
	serviceConfigPath := getServiceConfigPath()
	if serviceConfigPath == "" {
		msg := "ERROR: service config file not found"
		fmt.Println(msg)
		return
	}

	parseServiceConfig(serviceConfigPath, s)

	localhostConfigPath := getLocalhostConfigPath()
	if localhostConfigPath != "" {
		parseServiceConfig(localhostConfigPath, s)
	}

	parseEnv(s)
}

func parseEnv(s *ServiceConfig) {
	env := strings.ToLower(os.Getenv(envEnv))
	if env != "" {
		s.data.Env = env
	}
	serverPort := os.Getenv(envPort)
	if serverPort != "" {
		s.data.HTTP.Port = serverPort
	}

	redisHost := os.Getenv(envRedisHost)
	if redisHost != "" {
		s.data.Redis.Host = redisHost
	}
	redisPort := os.Getenv(envRedisPort)
	if redisPort != "" {
		s.data.Redis.Port = redisPort
	}

	dbHost := os.Getenv(envDBHost)
	if dbHost != "" {
		s.data.Database.HostAddress = dbHost
	}
	dbPort := os.Getenv(envDBPort)
	if dbPort != "" {
		s.data.Database.Port = dbPort
	}
	dbUsername := os.Getenv(envDBUsername)
	if dbUsername != "" {
		s.data.Database.User = dbUsername
	}
	dbPassword := os.Getenv(envDBPassword)
	if dbPassword != "" {
		s.data.Database.Password = dbPassword
	}

	privateKey := os.Getenv(envLoginPrivateKey)
	if privateKey != "" {
		s.data.Auth.PrivateKey = privateKey
	}
	publicKey := os.Getenv(envLoginPublicKey)
	if publicKey != "" {
		s.data.Auth.PublicKey = publicKey
	}
	accessTokenExpiration := os.Getenv(envAccessTokenExpiration)
	if accessTokenExpiration != "" {
		expiryTime, _ := strconv.Atoi(accessTokenExpiration)
		if expiryTime != 0 {
			s.data.Auth.AccessTokenExpirationMillis = uint32(expiryTime)
		}
	}
	refreshTokenExpiration := os.Getenv(envRefreshTokenExpiration)
	if refreshTokenExpiration != "" {
		expiryTime, _ := strconv.Atoi(refreshTokenExpiration)
		if expiryTime != 0 {
			s.data.Auth.RefreshTokenExpirationMillis = uint32(expiryTime)
		}
	}
}

func getLocalhostConfigPath() string {
	localhostConfigPath := os.Getenv(localhostConfigPath)
	if localhostConfigPath != "" {
		return localhostConfigPath
	}
	return getDefaultServiceConfigPath(localhostDefaultConfigPath)
}

func parseServiceConfig(configPath string, serviceConfig *ServiceConfig) {
	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		return
	}
	suffix := getSuffix(configPath)
	switch suffix {
	case "json":
		parseJSON(data, serviceConfig)
	case "yaml", "yml":
		parseYAML(data, serviceConfig)
	}
}

func parseYAML(data []byte, serviceConfig *ServiceConfig) {
	err := yaml.Unmarshal(data, &serviceConfig.data)
	serverPort := serviceConfig.data.HTTP.Port
	if !isValidPort(serverPort) {
		panic("invalid port number")
	}
	if err != nil {
		panic("can not parse yaml-yml service config file")
	}
}

func parseJSON(data []byte, serviceConfig *ServiceConfig) {
	err := json.Unmarshal(data, &serviceConfig.data)
	serverPort := serviceConfig.data.HTTP.Port
	if !isValidPort(serverPort) {
		panic("invalid port number")
	}
	if err != nil {
		panic("can not parse json service config file")
	}
}

func isValidPort(port string) bool {
	pattern := `^(^$|[0-9]{1,4}|[1-5][0-9]{4}|6[0-4][0-9]{3}|65[0-4][0-9]{2}|655[0-2][0-9]|6553[0-5])$`
	isMatched, _ := regexp.MatchString(pattern, port)
	return isMatched
}

func getSuffix(configPath string) string {
	arr := strings.Split(configPath, ".")
	suffix := arr[len(arr)-1]
	return suffix
}

func getServiceConfigPath() string {
	serviceConfigPath := os.Getenv(serviceConfigPath)
	if serviceConfigPath != "" {
		return serviceConfigPath
	}
	return getDefaultServiceConfigPath(defaultConfigPath)
}

func getDefaultServiceConfigPath(fileName string) string {
	paths := []string{
		path.Join("etc", fileName+".json"),
		path.Join("etc", fileName+".yaml"),
		path.Join("etc", fileName+".yml"),
	}
	for _, p := range paths {
		if checkFileExist(p) {
			return p
		}
	}
	return ""
}

func checkFileExist(filePath string) bool {
	if _, err := os.Stat(filePath); err == nil {
		return true
	} else if os.IsNotExist(err) {
		return false
	}
	return false
}

// GetServiceConfig returns the service config interface
func GetServiceConfig() ServiceConfigInterface {
	return defaultServiceConfig
}
