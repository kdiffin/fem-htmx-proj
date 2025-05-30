package router

import (
	"encoding/json"
	"net/http"
	db "swagtask/internal/db/generated"
	"swagtask/internal/middleware"
	"swagtask/internal/tag"
	"swagtask/internal/template"
	"swagtask/internal/utils"
)

func SetupTagRoutes(mux *http.ServeMux, queries *db.Queries, templates *template.Template) {

	mux.Handle("GET /tags/{$}", middleware.HandlerWithVaultIdFromUser(queries, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tag.HandlerGetTags(w, r, queries, templates)
	})))
	mux.Handle("PUT /tags/{id}/{$}", middleware.HandlerWithVaultIdFromUser(queries, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tag.HandlerUpdateTag(w, r, queries, templates)
	})))

	mux.Handle("DELETE /tags/{id}/{$}", middleware.HandlerWithVaultIdFromUser(queries, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tag.HandlerDeleteTag(w, r, queries, templates)
	})))
	mux.Handle("POST /tags/{id}/tasks/{$}", middleware.HandlerWithVaultIdFromUser(queries, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tag.HandlerAddTaskToTag(w, r, queries, templates)
	})))
	mux.Handle("DELETE /tags/{id}/tasks/{$}", middleware.HandlerWithVaultIdFromUser(queries, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tag.HandlerRemoveTaskFromTag(w, r, queries, templates)
	})))
	mux.HandleFunc("GET /json", func(w http.ResponseWriter, r *http.Request) {
		// Prepare the response data
		response := map[string]string{
			"message": "hello world",
		}

		// Set content-type to JSON
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		// Encode and send JSON
		if err := json.NewEncoder(w).Encode(response); err != nil {
			utils.LogError("failed to encode json", err)
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}
	})

}
