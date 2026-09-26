package request

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/go-playground/validator/v10"
)

var requestValidator = validator.New()

type validatable interface {
	Validate() error
}

func DecodeAndValidateRequest(r *http.Request, dest any) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(dest); err != nil {
		return fmt.Errorf(
			"decode json: %v",
			err,
		)
	}

	var extra any
	err := decoder.Decode(&extra)
	if !errors.Is(err, io.EOF) {

		return fmt.Errorf(
			"request body contains extra data: %v",
			err,
		)
	}

	v, ok := dest.(validatable)
	if ok {
		err = v.Validate()
	} else {
		err = requestValidator.Struct(dest)
	}

	if err != nil {
		return fmt.Errorf(
			"request validation: %v",
			err,
		)
	}

	return nil
}
