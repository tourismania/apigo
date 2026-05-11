package httpx

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/go-playground/validator/v10"
)

// ErrBadJSON wraps any JSON decoding failure.
var ErrBadJSON = errors.New("invalid JSON body")

// DecodeJSON decodes the request body into v and validates it (if a
// validator is supplied). The caller decides how to surface the error.
func DecodeJSON(r *http.Request, v any, validate *validator.Validate) error {
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(v); err != nil {
		if errors.Is(err, io.EOF) {
			return ErrBadJSON
		}
		return errors.Join(ErrBadJSON, err)
	}
	if validate != nil {
		if err := validate.Struct(v); err != nil {
			return err
		}
	}
	return nil
}
