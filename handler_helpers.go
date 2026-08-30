package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func respondWithError(w http.ResponseWriter, code int, msg string, err error) {
	type errReturn struct {
		Error string `json:"error"`
	}

	errorReturn := errReturn{
		Error: msg,
	}

	errReturnData, err := json.Marshal(errorReturn)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}

	//fmt.Println(err)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(errReturnData)
}

func respondWithJson(w http.ResponseWriter, code int, payload interface{}) {
	respondData, err := json.Marshal(payload)
	if err != nil {
		respondWithError(w, 400, "Something went wrong", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(respondData)
}
