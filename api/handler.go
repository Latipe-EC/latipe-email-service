package api

import (
	"encoding/json"
	"log"
	"net/http"
)

type HandlerAPI struct {
}

func NewHandlerApi() HandlerAPI {
	return HandlerAPI{}
}

func (api *HandlerAPI) WelcomeHandler(w http.ResponseWriter, r *http.Request) {
	response := struct {
		Message string `json:"message"`
		Version string `json:"version"`
	}{
		Message: "Email service was developed by TienDat",
		Version: "v0.0.1",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Printf("Error encoding JSON response: %v", err)
	}
}
