package db

import (
	"GoRecordurbate/modules/file"
	"crypto/rand"
	"encoding/hex"

	"golang.org/x/crypto/bcrypt"
)

type API_secrets struct {
	Keys []api `json:keys`
}

type api struct {
	User string `json:user`
	Name string `json:name`
	Key  string `json:secret`
}

var API API_secrets

func (a *API_secrets) Load() error {
	return file.ReadJson(file.API_keys_file, &a)
}

func (a API_secrets) NewKey() api {
	return api{}
}

// GenerateAPIKey creates a secure random API key
func (a *api) GenerateAPIKey(length int) (string, error) {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func (a *api) HashAPIKey(apiKey string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(apiKey), bcrypt.DefaultCost)
	return string(hash), err
}

func VerifyAPIKey(hashedKey, apiKey string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashedKey), []byte(apiKey)) == nil
}
