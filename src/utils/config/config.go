package config

import (
	"encoding/json"
	"os"

	"github.com/irdaislakhuafa/go-sdk/codes"
	"github.com/irdaislakhuafa/go-sdk/errors"
	"github.com/irdaislakhuafa/go-sdk/files"
)

type (
	Meta struct {
		Title       string
		Description string
		Version     string
		Host        string
		BasePath    string
	}

	Gin struct {
		Port           string
		TimeoutSeconds int
		Mode           string
		Cors           struct {
			Mode string
		}
	}

	Log struct {
		Level string
	}

	DB struct {
		Master struct {
			Username string
			Password string
			Host     string
			Port     string
			DBName   string
			Ssl      bool
			Options  struct{}
		}
	}

	Secrets struct {
		Key string
	}

	Token struct {
		ExpirationMinutes int64
	}

	Config struct {
		Meta    Meta
		Gin     Gin
		Log     Log
		DB      DB
		Secrets Secrets
		Token   Token
	}
)

func ReadFileJSON(pathToFile string) (Config, error) {
	if !files.IsExist(pathToFile) {
		return Config{}, errors.NewWithCode(codes.CodeStorageNoFile, "file '%s' not found", pathToFile)
	}

	fileBytes, err := os.ReadFile(pathToFile)
	if err != nil {
		return Config{}, errors.NewWithCode(codes.CodeStorageNoFile, "cannot read file '%v': %v", pathToFile, err.Error())
	}
	result := Config{}
	if err := json.Unmarshal(fileBytes, &result); err != nil {
		return Config{}, errors.NewWithCode(codes.CodeJSONUnmarshalError, "cannot parse '%v': %v", pathToFile, err.Error())
	}

	return result, nil
}
