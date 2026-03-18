package handlers

import (
	"net/http"

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
	shortUrl := handler.Storage.ShortenUrl()
	handler.Storage.StoreUrl(shortUrl, *req.Url)
	type CreateShortURLResp struct {
		ShortURL string `json:"short_url"`
	}
	helpers.WriteJson(w, http.StatusCreated, CreateShortURLResp{
		ShortURL: shortUrl,
	})
}

func (handler *Handler) Redirect(w http.ResponseWriter, r *http.Request) {
	code := r.PathValue("code")
	if code == "" {
		helpers.WriteError(w, http.StatusBadRequest, helpers.ErrorResponse{Error: "Path must be valid."})
		return
	}
	destination := handler.Storage.GetDestination(code)
	if destination == nil {
		helpers.WriteError(w, http.StatusNotFound, helpers.ErrorResponse{Error: "Url Not Found"})
		return
	}
	// todo: think about redirecting the user to the destination, and how to validate the url, and handle delineations.
}
