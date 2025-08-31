package models

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/sohWenMing/lenslocked/helpers"
)

type EnvVars struct {
	IsDev         bool
	CSRFSecretKey string
	BaseUrl       string
}

func LoadEnv(path string) (envVars EnvVars, err error) {
	err = godotenv.Load(path)
	if err != nil {
		return envVars, err
	}
	isDevVal, err := getIsDevVal()
	if err != nil {
		return envVars, err
	}
	csrfSecretKey, err := getCSRFKey()
	if err != nil {
		return envVars, err
	}
	baseURL, err := getBaseURL()
	if err != nil {
		return envVars, err
	}
	envVars.IsDev = isDevVal
	envVars.CSRFSecretKey = csrfSecretKey
	envVars.BaseUrl = baseURL
	return envVars, nil
}

func getEnvVar(input string) (envVarString string, err error) {
	envVarString = os.Getenv(helpers.TrimSpaceToUpper(input))
	if envVarString == "" {
		return "", fmt.Errorf("env var with name %s could not be found", input)
	}
	return envVarString, nil
}

func getBaseURL() (string, error) {
	baseUrl, err := getEnvVar("BASEURL")
	if err != nil {
		return "", err
	}
	return baseUrl, nil
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
func getCSRFKey() (string, error) {
	CSRFSecretKey, err := getEnvVar("SECRETKEY")
	if err != nil {
		return "", err
	}
	return CSRFSecretKey, nil
}
