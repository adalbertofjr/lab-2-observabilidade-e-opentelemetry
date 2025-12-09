package api

import (
	"io"
	"net/http"
)

type CEPHandler struct {
	CEP string `json:"cep"`
}

func NewCEPHandler() *CEPHandler {
	return &CEPHandler{}
}

func (c *CEPHandler) CEPValidate(w http.ResponseWriter, r *http.Request) {
	req := r.Body
	defer req.Close()

	body, err := io.ReadAll(req)
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(body)
}
