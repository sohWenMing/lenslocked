package models

import "golang.org/x/crypto/bcrypt"

// generates a hash of the plainTextPassword passed in, will return "" and error if error occurs during hashing
func GenerateBcryptHash(plainTextPassword string) (hash string, err error) {
	hashBytes, err := bcrypt.GenerateFromPassword(
		[]byte(plainTextPassword),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return "", err
	}
	return string(hashBytes), nil
}
