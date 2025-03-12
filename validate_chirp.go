package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func handlerValidateChirp(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	type requestBody struct {
		Body string `json:"body"`
	}

	type responseBody struct {
		Valid bool `json:"valid"`
	}

	data, err := io.ReadAll(r.Body)
	if err != nil {
		respondWithError(w, 500, "couldn't read request")
		return
	}

	fmt.Println(string(data))

	params := requestBody{}
	err = json.Unmarshal(data, &params)
	if err != nil {
		respondWithError(w, 500, "couldn't unmarshal parameters")
		return
	}

	if len(params.Body) > 140 {
		respondWithError(w, 400, "Chirp is too long")
	}

	respondWithJSON(w, 200, responseBody{
		Valid: true,
	})

}
