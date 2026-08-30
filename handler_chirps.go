package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/nukz87/chirpy/internal/auth"
	"github.com/nukz87/chirpy/internal/database"
)

type Chirp struct {
	ID       uuid.UUID `json:"id"`
	CreateAt time.Time `json:"create_at"`
	UpdateAt time.Time `json:"update_at"`
	Body     string    `json:"body"`
	UserID   uuid.UUID `json:"user_id"`
}

func validateChirp(text string) (string, error) {
	profaneWords := []string{"kerfuffle", "sharbert", "fornax"}

	if len(text) > 140 {
		return "", errors.New("Chirp is too long")
	}
	loweredBody := strings.ToLower(text)
	splittedBody := strings.Fields(text)
	splittedLoweredBody := strings.Fields(loweredBody)

	for i, word := range splittedLoweredBody {
		for _, profane := range profaneWords {
			if word == profane {
				splittedBody[i] = "****"
				break
			}
		}
	}

	return strings.Join(splittedBody, " "), nil
}

func (cfg *apiConfig) handlerCreateChirps(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	if err := decoder.Decode(&params); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload", err)
		return
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload", err)
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid authorization", err)
		return
	}

	cleanedBody, err := validateChirp(params.Body)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload", err)
		return
	}

	chirp, err := cfg.db.CreateChirp(r.Context(), database.CreateChirpParams{
		Body:   cleanedBody,
		UserID: userID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create chirp", err)
		return
	}

	respondWithJson(w, 201, Chirp{
		ID:       chirp.ID,
		CreateAt: chirp.CreatedAt,
		UpdateAt: chirp.UpdatedAt,
		Body:     chirp.Body,
		UserID:   chirp.UserID,
	})
}

func (cfg *apiConfig) handlerGetAllChirps(w http.ResponseWriter, r *http.Request) {
	dbChirps, err := cfg.db.GetAllChirps(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get all chirps", err)
		return
	}

	chirps := make([]Chirp, len(dbChirps))
	for i, c := range dbChirps {
		chirps[i] = Chirp{
			ID:       c.ID,
			CreateAt: c.CreatedAt,
			UpdateAt: c.UpdatedAt,
			Body:     c.Body,
			UserID:   c.UserID,
		}
	}

	respondWithJson(w, 200, chirps)
}

func (cfg *apiConfig) handlerGetChirpByID(w http.ResponseWriter, r *http.Request) {
	chirpIDString := r.PathValue("chirpID")

	chirpID, err := uuid.Parse(chirpIDString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid chirp ID", err)
		return
	}

	chirp, err := cfg.db.GetChirpByID(r.Context(), chirpID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Chirp not found", err)
		return
	}
	respondWithJson(w, 200, Chirp{
		ID:       chirp.ID,
		CreateAt: chirp.CreatedAt,
		UpdateAt: chirp.UpdatedAt,
		Body:     chirp.Body,
		UserID:   chirp.UserID,
	})
}
