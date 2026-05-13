package main

import (
	"net/http"
	"sort"
	"github.com/e-300/http-server-go/internal/database"
	"github.com/google/uuid"
)

// retriving chirp via chirp id in the request path
func (cfg *apiConfig) handlerChirpsGet(w http.ResponseWriter, r *http.Request) {
	chirpIdStr := r.PathValue("chirpID")
	chirpId, err := uuid.Parse(chirpIdStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Bad chirp ID", err)
		return
	}
	dbChirps, err := cfg.db.GetChirp(r.Context(), chirpId)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not get Chirp", err)
		return
	}

	respondWithJSON(w, http.StatusOK, Chirp{
		ID: dbChirps.ID,
		CreatedAt: dbChirps.CreatedAt,
		UpdatedAt: dbChirps.UpdatedAt,
		UserID: dbChirps.UserID,
		Body: dbChirps.Body,
	})
}

func sortChirps(r *http.Request, chirps []Chirp) []Chirp {
	sortParam := r.URL.Query().Get("sort")
	if sortParam == "desc" {
		sort.Slice(chirps, func(i, j int) bool {
			return chirps[i].CreatedAt.After(chirps[j].CreatedAt)
		})
	}
	return chirps
}

func getAuthorIdFromReq(r *http.Request) (uuid.UUID, error) {
	authorIdStr := r.URL.Query().Get("author_id")
	if authorIdStr == "" {
		return uuid.Nil, nil
	}
	authorID, err := uuid.Parse(authorIdStr)
	if err != nil {
		return uuid.Nil, err
	}
	return authorID, nil
}

func (cfg *apiConfig) handlerChirpsRetrieve(w http.ResponseWriter, r *http.Request) {
	authorID, err := getAuthorIdFromReq(r)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Invalid Author ID", err)
		return

	}

	dbChirps := []database.Post{}

	if authorID != uuid.Nil {
		dbChirps, err = cfg.db.GetChirpsForAuthor(r.Context(), authorID)
	} else {
		dbChirps, err = cfg.db.GetAllChirps(r.Context())
	}
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldnt retrieve posts", err)
		return
	}

	chirps := []Chirp{}
	for _, chirp := range dbChirps {
		chirps = append(chirps, Chirp{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			UserID:    chirp.UserID,
			Body:      chirp.Body,
		})
	}

	chirps = sortChirps(r, chirps)

	respondWithJSON(w, http.StatusOK, chirps)
}
