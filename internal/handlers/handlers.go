package handlers

import (
	"net/http"

	"github.com/Arush71/url-shortener/internal/helpers"
)

func HandleShortening(w http.ResponseWriter, r *http.Request) {
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
	type CreateShortURLResp struct {
		ShortURL string `json:"short_url"`
	}
	helpers.WriteJson(w, http.StatusCreated, CreateShortURLResp{
		ShortURL: "short",
	})
}
