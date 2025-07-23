package codec

import (
	"encoding/json"
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
