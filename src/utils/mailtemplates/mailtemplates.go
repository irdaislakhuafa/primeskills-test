package mailtemplates

import (
	"os"

	"github.com/irdaislakhuafa/go-sdk/codes"
	"github.com/irdaislakhuafa/go-sdk/errors"
	"github.com/irdaislakhuafa/go-sdk/strformat"
)

const TEMPLATE_DIR = "docs/templates/email"
const (
	REGISTER_VERIFICATION = TEMPLATE_DIR + "/register_verification.html"
	RESET_PASSWORD        = TEMPLATE_DIR + "/reset_password.html"
)

func ReadAndParse(file string, params any) (string, error) {
	fileBytes, err := os.ReadFile(file)
	if err != nil {
		return "", errors.NewWithCode(codes.CodeStorageNoFile, "Cannot read file '%v'", file)
	}

	result, err := strformat.T(string(fileBytes), params)
	if err != nil {
		return "", errors.NewWithCode(codes.CodeInternalServerError, "Cannot parse template '%s'", file)
	}

	return result, nil
}
