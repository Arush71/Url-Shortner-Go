package handlers

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/Arush71/url-shortener/internal/helpers"
	"github.com/Arush71/url-shortener/internal/shortner"
)

type Handler struct {
	Storage *shortner.Storage
}

func (handler *Handler) HandleShortening(w http.ResponseWriter, r *http.Request) {
	type reqT struct {
		Url *string `json:"url"`
	}
	var req reqT
	if er := helpers.ReadJson(r, &req); er != nil {
		helpers.WriteError(w, http.StatusBadRequest, helpers.ErrorResponse{
			Error: "bad request",
		})
		return
	}
	if req.Url == nil || *req.Url == "" {
		helpers.WriteError(w, http.StatusBadRequest, helpers.ErrorResponse{
			Error: "the url should be provided.",
		})
		return
	}
	url, err := url.Parse(*req.Url)
	if err != nil || url.Host == "" || url.Scheme == "" {
		helpers.WriteError(w, http.StatusBadRequest, helpers.ErrorResponse{Error: "Invalid Url."})
		return
	}
	handler.Storage.Mu.Lock()
	shortUrl := handler.Storage.ShortenUrl()
	handler.Storage.Store(shortUrl, *req.Url)
	handler.Storage.Mu.Unlock()
	type CreateShortURLResp struct {
		ShortURL string `json:"short_url"`
	}
	helpers.WriteJson(w, http.StatusCreated, CreateShortURLResp{
		ShortURL: fmt.Sprintf("http://localhost:8080/%s", shortUrl),
	})
}

func (handler *Handler) Redirect(w http.ResponseWriter, r *http.Request) {
	code := r.PathValue("code")
	if code == "" {
		helpers.WriteError(w, http.StatusBadRequest, helpers.ErrorResponse{Error: "Path must be valid."})
		return
	}
	handler.Storage.Mu.Lock()
	destination := handler.Storage.GetDestination(code)
	if destination == nil {
		helpers.WriteError(w, http.StatusNotFound, helpers.ErrorResponse{Error: "Url Not Found"})
		handler.Storage.Mu.Unlock()
		return
	}
	handler.Storage.UpdateCounter(code)
	handler.Storage.Mu.Unlock()
	http.Redirect(w, r, *destination, http.StatusFound)
}

func (handler *Handler) Stats(w http.ResponseWriter, r *http.Request) {
	code := r.PathValue("code")
	if code == "" {
		helpers.WriteError(w, http.StatusBadRequest, helpers.ErrorResponse{Error: "Path must be valid."})
		return
	}
	handler.Storage.Mu.Lock()
	url, counter, err := handler.Storage.GetStats(code)
	handler.Storage.Mu.Unlock()
	if err == nil {
		type ResStat struct {
			Url     string `json:"url"`
			Counter int    `json:"counter"`
		}
		helpers.WriteJson(w, http.StatusOK, ResStat{
			Url:     url,
			Counter: counter,
		})
		return
	}
	helpers.WriteError(w, http.StatusNotFound, helpers.ErrorResponse{Error: err.Error()})
}
