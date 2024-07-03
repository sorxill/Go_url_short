package url

import (
	"errors"
	api_response "go_url_short/internal/http_server/api"
	"go_url_short/internal/http_server/random"
	"go_url_short/internal/storage"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

const aliasLength = 6

type Requests struct {
	URL   string `json:"url" validate:"required,url"` // required - поле обязательное.
	Alias string `json:"alias,omitempty"`             // omitepmpty - указание, что этот параметр не обязательно присутсвует в json.
}

type Response struct {
	api_response.Response
	Alias string `json:"alias,omitempty"`
}

type URLSaver interface {
	SaveURL(urlToSave string, alias string) error
}

func SaveNew(log *slog.Logger, urlSaver URLSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.save.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Requests

		req.URL = r.URL.Query().Get("url")
		req.Alias = r.URL.Query().Get("alias")

		if req.URL == "" {
			log.Error("failed to get url")
			render.JSON(w, r, api_response.Error("no correct URL"))

			return
		}

		// err := render.DecodeJSON(r.Body, &req)
		// if err != nil {
		// 	log.Error("failed to decode requests body", slog.Any("error:", err))

		// 	render.JSON(w, r, api_response.Error("failed to decode requests"))

		// 	return
		// }
		/*
			Эта реализация у меня не сработала, написал свою через получение параметра из URL строки.
			Реализация выше, проводит ту же проверку, но проще, через указание параметра в строке через знак ?url=...&alias=...
		*/

		log.Info("request body decoded", slog.Any("request:", req))

		if err := validator.New().Struct(req); err != nil {
			log.Error("invalid request", slog.Any("error: ", err))

			render.JSON(w, r, api_response.Error("invalid request"))

			return
		}

		alias := req.Alias
		if alias == "" {
			alias = random.NewRandomString(aliasLength)
		}

		err := urlSaver.SaveURL(req.URL, alias)
		if errors.Is(err, storage.ErrURLExists) {
			log.Info("url already exists", slog.String("url", req.URL))

			render.JSON(w, r, api_response.Error("url already exists"))

			return
		}

		if err != nil {
			log.Error("failed to add url", slog.Any("error: ", err))

			render.JSON(w, r, api_response.Error("falied to add url"))

			return
		}

		log.Info("url added", slog.Any("status: ", api_response.StatusOK))

		ResponseOk(w, r, alias)
	}
}

func ResponseOk(w http.ResponseWriter, r *http.Request, alias string) {
	render.JSON(w, r, Response{
		Response: api_response.OK(),
		Alias:    alias,
	})
}
