package handlers

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/Arush71/url-shortener/internal/db"
	"github.com/Arush71/url-shortener/internal/helpers"
	"github.com/Arush71/url-shortener/internal/shortner"
)

type Handler struct {
	Q      *db.Queries
	DB     *sql.DB
	AppUrl string
}

func (handler *Handler) HandleShortening(w http.ResponseWriter, r *http.Request) {
	type reqT struct {
		Url *string `json:"url"`
	}
	var req reqT
	if er := helpers.ReadJson(r, &req); er != nil {
		helpers.WriteError(w, http.StatusBadRequest, helpers.ErrorResponse{
			Error: "BAD_REQUEST",
		})
		return
	}
	if req.Url == nil || *req.Url == "" {
		helpers.WriteError(w, http.StatusBadRequest, helpers.ErrorResponse{
			Error: "BAD_REQUEST",
		})
		return
	}
	url, err := url.Parse(*req.Url)
	if err != nil || url.Host == "" || url.Scheme == "" {
		helpers.WriteError(w, http.StatusBadRequest, helpers.ErrorResponse{Error: "BAD_REQUEST"})
		return
	}
	tx, err := handler.DB.BeginTx(r.Context(), nil)
	if err != nil {
		helpers.WriteError(w, http.StatusInternalServerError, helpers.ErrorResponse{Error: "INTERNAL_SERVER_ERROR"})
		return
	}
	defer tx.Rollback()
	qtx := handler.Q.WithTx(tx)
	id, err := qtx.GetNextURLID(r.Context())
	if err != nil {
		helpers.WriteError(w, http.StatusInternalServerError, helpers.ErrorResponse{Error: "INTERNAL_SERVER_ERROR"})
		return
	}
	code := shortner.GetCodeFromId(id)
	if err = qtx.CreateUrl(r.Context(), db.CreateUrlParams{
		ID:          id,
		Code:        code,
		OriginalUrl: *req.Url,
		Counter:     0,
	}); err != nil {
		helpers.WriteError(w, http.StatusInternalServerError, helpers.ErrorResponse{Error: "INTERNAL_SERVER_ERROR"})
		return
	}
	if err = tx.Commit(); err != nil {
		helpers.WriteError(w, http.StatusInternalServerError, helpers.ErrorResponse{Error: "INTERNAL_SERVER_ERROR"})
		return
	}
	type CreateShortURLResp struct {
		ShortURL string `json:"short_url"`
	}
	helpers.WriteJson(w, http.StatusCreated, CreateShortURLResp{
		ShortURL: fmt.Sprintf("%s/%s", handler.AppUrl, code),
	})
}

func (handler *Handler) Redirect(w http.ResponseWriter, r *http.Request) {
	code := r.PathValue("code")
	if code == "" {
		helpers.WriteError(w, http.StatusBadRequest, helpers.ErrorResponse{Error: "BAD_REQUEST"})
		return
	}
	originalUrl, err := handler.Q.UpdateAndRedirect(r.Context(), code)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			helpers.WriteError(w, http.StatusNotFound, helpers.ErrorResponse{Error: "BAD_REQUEST"})
			return
		}
		helpers.WriteError(w, http.StatusInternalServerError, helpers.ErrorResponse{Error: "INTERNAL_SERVER_ERROR"})
		return
	}
	http.Redirect(w, r, originalUrl, http.StatusFound)
}

func (handler *Handler) Stats(w http.ResponseWriter, r *http.Request) {
	code := r.PathValue("code")
	if code == "" {
		helpers.WriteError(w, http.StatusBadRequest, helpers.ErrorResponse{Error: "BAD_REQUEST"})
		return
	}
	stats, err := handler.Q.GetStats(r.Context(), code)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			helpers.WriteError(w, http.StatusNotFound, helpers.ErrorResponse{Error: "NOT_FOUND"})
			return
		}
		helpers.WriteError(w, http.StatusInternalServerError, helpers.ErrorResponse{Error: "INTERNAL_SERVER_ERROR"})
		return
	}
	helpers.WriteJson(w, http.StatusOK, stats)
}
