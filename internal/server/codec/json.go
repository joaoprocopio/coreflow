package codec

import (
	"coreflow/internal/server/validator"
	"encoding/json"
	"fmt"
	"net/http"
)

func WriteEncodedJSON[T any](w http.ResponseWriter, r *http.Request, s int, v T) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(s)

	if err := json.NewEncoder(w).Encode(v); err != nil {
		return err
	}

	return nil
}

func ReadDecodedJSON[T any](r *http.Request) (T, error) {
	var v T

	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		return v, err
	}

	return v, nil
}

func ReadDecodedValidJSON[T validator.Validator](r *http.Request) (T, validator.Problems, error) {
	var v T

	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		return v, nil, err
	}

	if problems := v.Valid(r.Context()); len(problems) > 0 {
		return v, problems, fmt.Errorf("invalid %T: %d problems", v, len(problems))
	}

	return v, nil, nil
}
