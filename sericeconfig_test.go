package serviceconfig

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func unsetEnvVars(t *testing.T, envVars ...string) {
	for _, v := range envVars {
		err := os.Unsetenv(v)
		assert.NoErrorf(t, err, "clean up env variables failed: could not unset env %s", v)
	}
}

func TestInitializeServiceConfig(t *testing.T) {
	//defer unsetEnvVars(t, serviceConfigPath)
	_ = os.Setenv(serviceConfigPath, "./test_configs/service.yml")

	expected := &ServiceConfig{
		serviceConfig{
			Database: DatabaseConfig{
				HostAddress: "123.456.789.10",
				Port:        "5432",
				Name:        "fitx",
				User:        "fitx",
				Password:    "123456",
			},
			HTTP: HTTPServer{
				BaseURL: "localhost",
				Host:    "localhost",
				Port:    "9876",
			},
			Redis: RedisConfig{
				Host: "localhost",
				Port: "6379",
			},
			Auth: AuthConfig{
				PrivateKey:                   "private_key",
				PublicKey:                    "public_key",
				AccessTokenExpirationMillis:  300000,
				RefreshTokenExpirationMillis: 2592000000,
			},
			MailService: MailServiceKey{
				ApiKey: "api_key",
			},
			Env: "test",
			Logger: LoggerConfig{
				ConfigPath: "config_path",
			},
		},
	}

	serviceConfig := &ServiceConfig{}
	initServiceDefaultConfigValue(serviceConfig)
	initServiceConfig(serviceConfig)
	assert.Equal(t, expected, serviceConfig, "test initializeServiceConfig failed")
}

func TestInitializeServiceConfigLocalhost(t *testing.T) {
	defer unsetEnvVars(t, serviceConfigPath, localhostConfigPath)
	_ = os.Setenv(serviceConfigPath, "./test_configs/service.yml")
	_ = os.Setenv(localhostConfigPath, "./test_configs/localhost_service.yml")

	expected := &ServiceConfig{
		serviceConfig{
			Database: DatabaseConfig{
				HostAddress: "localhost",
				Port:        "5432",
				Name:        "localhost_fitx",
				User:        "localhost_fitx",
				Password:    "123456",
			},
			HTTP: HTTPServer{
				BaseURL: "localhost",
				Host:    "localhost",
				Port:    "9876",
			},
			Redis: RedisConfig{
				Host: "localhost",
				Port: "6379",
			},
			Auth: AuthConfig{
				PrivateKey:                   "localhost_private_key",
				PublicKey:                    "localhost_public_key",
				AccessTokenExpirationMillis:  300000,
				RefreshTokenExpirationMillis: 2592000000,
			},
			MailService: MailServiceKey{
				ApiKey: "localhost_api_key",
			},
			Env: "test",
			Logger: LoggerConfig{
				ConfigPath: "localhost_config_path",
			},
		},
	}

	serviceConfig := &ServiceConfig{}
	initServiceDefaultConfigValue(serviceConfig)
	initServiceConfig(serviceConfig)

	assert.Equal(t, expected, serviceConfig, "Test initializeServiceConfig localhost failed")
}

func TestParseJSON(t *testing.T) {
	expected := &ServiceConfig{
		serviceConfig{
			Database: DatabaseConfig{
				HostAddress: "123.456.789.10",
				Port:        "5432",
				Name:        "fitx",
				User:        "fitx",
				Password:    "123456",
			},
			HTTP: HTTPServer{
				BaseURL: "localhost",
				Host:    "localhost",
				Port:    "9876",
			},
			Redis: RedisConfig{
				Host: "localhost",
				Port: "6379",
			},
			Auth: AuthConfig{
				PrivateKey:                   "private_key",
				PublicKey:                    "public_key",
				AccessTokenExpirationMillis:  300000,
				RefreshTokenExpirationMillis: 2592000000,
			},
			MailService: MailServiceKey{
				ApiKey: "api_key",
			},
			Env: "test",
			Logger: LoggerConfig{
				ConfigPath: "config_path",
			},
		},
	}

	filePath := "./test_configs/service.json"
	data, _ := ioutil.ReadFile(filePath)
	serviceConfig := &ServiceConfig{}
	parseJSON(data, serviceConfig)

	assert.Equal(t, expected, serviceConfig, "test parseJSON failed")
}

func TestParseYAML(t *testing.T) {
	filePaths := []string{
		"./test_configs/service.yml",
		"./test_configs/service.yaml",
	}

	expected := &ServiceConfig{
		serviceConfig{
			Database: DatabaseConfig{
				HostAddress: "123.456.789.10",
				Port:        "5432",
				Name:        "fitx",
				User:        "fitx",
				Password:    "123456",
			},
			HTTP: HTTPServer{
				BaseURL: "localhost",
				Host:    "localhost",
				Port:    "9876",
			},
			Redis: RedisConfig{
				Host: "localhost",
				Port: "6379",
			},
			Auth: AuthConfig{
				PrivateKey:                   "private_key",
				PublicKey:                    "public_key",
				AccessTokenExpirationMillis:  300000,
				RefreshTokenExpirationMillis: 2592000000,
			},
			MailService: MailServiceKey{
				ApiKey: "api_key",
			},
			Env: "test",
			Logger: LoggerConfig{
				ConfigPath: "config_path",
			},
		},
	}

	for _, filePath := range filePaths {
		data, _ := ioutil.ReadFile(filePath)
		serviceConfig := &ServiceConfig{}
		parseYAML(data, serviceConfig)

		assert.Equal(t, expected, serviceConfig, "test parseYAML failed")
	}
}

func TestGetSuffix(t *testing.T) {
	tests := []struct {
		fileName string
		expected string
	}{
		{"nam.json", "json"},
		{"nam.yaml", "yaml"},
		{"nam.yml", "yml"},
		{"/etc/service.yml", "yml"},
		{"/etc/service.json", "json"},
		{"/etc/service.yaml", "yaml"},
	}

	for _, test := range tests {
		assert.Equal(t, test.expected, getSuffix(test.fileName), "test getSuffix failed")
	}
}

func TestCheckFileExist(t *testing.T) {
	tests := []struct {
		fileName string
		expected bool
	}{
		{"./test_configs/localhost_service.yml", true},
		{"./test_configs/service.json", true},
		{"./test_configs/service.yml", true},
		{"./test_configs/nam.yml", false},
	}

	for _, test := range tests {
		assert.Equal(t, test.expected, checkFileExist(test.fileName), "test checkFileExist failed")
	}
}

func TestServiceConfig_GetDataBaseConfig(t *testing.T) {
	expected := DatabaseConfig{
		HostAddress: "123.456.789.10",
		Port:        "5432",
		Name:        "fitx",
		User:        "fitx",
		Password:    "123456",
	}
	testServiceConfig := &ServiceConfig{
		serviceConfig{
			Database: expected,
		},
	}

	assert.Equal(t, expected, testServiceConfig.GetDataBaseConfig(), "test GetDatabaseConfig failed")
}

func TestServiceConfig_GetAuthConfig(t *testing.T) {
	expected := AuthConfig{
		PrivateKey:                   "private_key",
		PublicKey:                    "public_key",
		AccessTokenExpirationMillis:  300000,
		RefreshTokenExpirationMillis: 2592000000,
	}
	testServiceConfig := &ServiceConfig{
		serviceConfig{
			Auth: expected,
		},
	}

	assert.Equal(t, expected, testServiceConfig.GetAuthConfig(), "test GetAuthConfig failed")
}

func TestServiceConfig_GetRedisConfig(t *testing.T) {
	expected := RedisConfig{
		Host: "localhost_123",
		Port: "1234",
	}
	testServiceConfig := &ServiceConfig{
		serviceConfig{
			Redis: expected,
		},
	}

	assert.Equal(t, expected, testServiceConfig.GetRedisConfig(), "test GetRedisConfig failed")
}

func TestServiceConfig_GetMailServiceConfig(t *testing.T) {
	expected := MailServiceKey{
		ApiKey: "api_key",
	}
	testServiceConfig := &ServiceConfig{
		serviceConfig{
			MailService: expected,
		},
	}

	assert.Equal(t, expected, testServiceConfig.GetMailServiceConfig(), "test GetMailService failed")
}

func TestServiceConfig_GetServerPort(t *testing.T) {
	expected := HTTPServer{
		Port: "1234",
	}
	testServiceConfig := &ServiceConfig{
		serviceConfig{
			HTTP: expected,
		},
	}

	assert.Equal(t, expected.Port, testServiceConfig.GetServerPort(), "test GetServerPort failed")
}

func TestServiceConfig_GetEnv(t *testing.T) {
	expected := "live"
	testServiceConfig := &ServiceConfig{
		serviceConfig{
			Env: expected,
		},
	}

	assert.Equal(t, expected, testServiceConfig.GetEnv(), "test GetEnv failed")
}

func TestServiceConfig_GetLoggerConfigPath(t *testing.T) {
	expected := LoggerConfig{
		ConfigPath: "config_path.yml",
	}
	testServiceConfig := &ServiceConfig{
		serviceConfig{
			Logger: expected,
		},
	}

	assert.Equal(t, expected.ConfigPath, testServiceConfig.GetLoggerConfigPath(), "test GetLoggerConfigPath failed")
}
