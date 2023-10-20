package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"slices"
	"testing"
)

type vkPaymentResponse struct {
	Response json.RawMessage `json:"response"`
}

func (v *vkPaymentResponse) Write(data []byte) (n int, err error) {
	v.Response = append(v.Response, data...)
	return len(data), nil
}

func WriteJSON(w http.ResponseWriter, v any, code int) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	var response vkPaymentResponse
	err := json.NewEncoder(&response).Encode(v)
	if err != nil {
		return err
	}
	return json.NewEncoder(w).Encode(response)
}

func AssertResponse(t *testing.T, response any, body io.ReadCloser) {
	var bodyResponse vkPaymentResponse
	if decodeErr := json.NewDecoder(body).Decode(&bodyResponse); decodeErr != nil {
		t.Fatalf("failed assert response: decode body, err: %s", decodeErr)
	}
	body.Close()
	rawMessage, encodeErr := json.Marshal(response)
	if encodeErr != nil {
		t.Fatalf("failed assert response: encode response, err: %s", encodeErr)
	}
	equal := slices.Equal(bodyResponse.Response, rawMessage)
	if !equal {
		t.Fatal("response not equal")
	}
}
