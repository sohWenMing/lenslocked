package models

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/sohWenMing/lenslocked/helpers"
)

type Envs struct{}

type EmailEnvs struct {
	Host     string
	Port     int
	Username string
	Password string
}

func (e *Envs) LoadPgConfig() (config PgConfig, err error) {
	config = PgConfig{}
	dbHost, err := e.GetDBHost()
	if err != nil {
		return config, err
	}
	port, err := e.GetDBPort()
	if err != nil {
		return config, err
	}
	user, err := e.GetDBUser()
	if err != nil {
		return config, err
	}
	dbPassword, err := e.GetDBPassword()
	if err != nil {
		return config, err
	}
	dbName, err := e.GetDBName()
	if err != nil {
		return config, err
	}
	sslMode, err := e.GetDBSSLMode()
	if err != nil {
		return config, err
	}
	return PgConfig{
		dbHost, port, user, dbPassword, dbName, sslMode,
	}, nil
}

func (e *Envs) LoadEmailEnvs() (emailEnvs *EmailEnvs, err error) {
	emailHost, err := e.GetEmailHost()
	if err != nil {
		return nil, err
	}
	emailPort, err := e.GetEmailPort()
	if err != nil {
		return nil, err
	}
	emailUserName, err := e.GetEmailUsername()
	if err != nil {
		return nil, err
	}
	emailPassword, err := e.GetEmailPassword()
	if err != nil {
		return nil, err
	}
	return &EmailEnvs{
		Host:     emailHost,
		Port:     emailPort,
		Username: emailUserName,
		Password: emailPassword,
	}, nil
}

func (e *Envs) GetIsDev() (bool, error) {
	isDevVal, err := getIsDevVal()
	if err != nil {
		return false, err
	}
	return isDevVal, nil
}

func (e *Envs) GetCSRFSecretKey() (string, error) {
	csrfKey, err := getEnvVar("CSRFSECRETKEY")
	if err != nil {
		return "", err
	}
	return csrfKey, nil
}
func (e *Envs) GetBaseURL() (string, error) {
	baseUrl, err := getEnvVar("BASEURL")
	if err != nil {
		return "", err
	}
	return baseUrl, nil
}
func (e *Envs) GetEmailHost() (string, error) {
	host, err := getEnvVar("EMAILHOST")
	if err != nil {
		return "", err
	}
	return host, nil
}
func (e *Envs) GetEmailPassword() (string, error) {
	password, err := getEnvVar("EMAILPASSWORD")
	if err != nil {
		return "", err
	}
	return password, nil
}
func (e *Envs) GetEmailUsername() (string, error) {
	password, err := getEnvVar("EMAILUSERNAME")
	if err != nil {
		return "", err
	}
	return password, nil
}
func (e *Envs) GetDBHost() (string, error) {
	dbHost, err := getEnvVar("DBHOST")
	if err != nil {
		return "", err
	}
	return dbHost, nil
}

func (e *Envs) GetDBPort() (string, error) {
	dbPort, err := getEnvVar("DBPORT")
	if err != nil {
		return "", err
	}
	_, err = strconv.Atoi(dbPort)
	if err != nil {
		return "", errors.New("dbPort could not be converted to numerical value")
	}
	return dbPort, nil
}
func (e *Envs) GetDBUser() (string, error) {
	dbHost, err := getEnvVar("DBUSER")
	if err != nil {
		return "", err
	}
	return dbHost, nil
}
func (e *Envs) GetDBPassword() (string, error) {
	dbPassword, err := getEnvVar("DBPASSWORD")
	if err != nil {
		return "", err
	}
	return dbPassword, nil
}
func (e *Envs) GetDBName() (string, error) {
	dbName, err := getEnvVar("DBNAME")
	if err != nil {
		return "", err
	}
	return dbName, nil
}
func (e *Envs) GetDBSSLMode() (string, error) {
	sslMode, err := getEnvVar("DBSSLMODE")
	if err != nil {
		return "", err
	}
	return sslMode, nil
}
func (e *Envs) GetEmailPort() (int, error) {
	port, err := getEmailPort()
	return port, err
}

func LoadEnv(path string) (envs *Envs, err error) {
	err = godotenv.Load(path)
	if err != nil {
		return nil, err
	}
	return envs, nil
}

func getEnvVar(input string) (envVarString string, err error) {
	envVarString = os.Getenv(helpers.TrimSpaceToUpper(input))
	if envVarString == "" {
		return "", fmt.Errorf("env var with name %s could not be found", input)
	}
	return envVarString, nil
}

func getIsDevVal() (bool, error) {
	isDevString, err := getEnvVar("ISDEV")
	if err != nil {
		return false, err
	}
	isDevVal, err := strconv.ParseBool(isDevString)
	if err != nil {
		return false, errors.New("ISDEV in .env file could not be parsed to a boolean")
	}
	return isDevVal, nil
}

func getEmailPort() (int, error) {
	portString, err := getEnvVar("PORT")
	if err != nil {
		return 0, err
	}
	port, err := strconv.Atoi(portString)
	if err != nil {
		return 0, err
	}
	return port, nil
}
