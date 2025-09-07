package models

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/sohWenMing/lenslocked/helpers"
)

type Envs struct {
}

func (e *Envs) GetIsDev() (bool, error) {
	isDevVal, err := getIsDevVal()
	if err != nil {
		return false, err
	}
	return isDevVal, nil
}

func (e *Envs) GetCSRFSecretKey() (string, error) {
	csrfKey, err := getEnvVar("SECRETKEY")
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
