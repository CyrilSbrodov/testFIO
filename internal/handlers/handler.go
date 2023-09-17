/*
Package handlers работа с эндпоинтами
*/
package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"testFIO/cmd/config"
	"testFIO/cmd/loggers"
	"testFIO/internal/storage"
	"testFIO/internal/storage/model"
	"testFIO/internal/storage/model/filter"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Handlers interface {
	Register(router *chi.Mux)
}

type Handler struct {
	storage.Service
	logger loggers.Logger
	cfg    config.ServerConfig
}

// Register создание роутеров
func (h *Handler) Register(r *chi.Mux) {
	r.Get("/value/", filter.MiddleWare(h.GetUsers(), h.cfg.Limit))
	r.Get("/value/*", h.GetUser())
	r.Post("/add/", h.CollectUsers())
	r.Post("/update/", h.UpdateUser())
	r.Post("/delete/", h.DeleteUser())
	r.Get("/ping", h.PingDB())
	r.Mount("/debug", middleware.Profiler())
}

func NewHandler(service storage.Service, logger loggers.Logger, cfg config.ServerConfig) Handlers {
	return &Handler{
		service,
		logger,
		cfg,
	}
}

// CollectUsers хендлердобавления новых пользователей
func (h *Handler) CollectUsers() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		h.logger.LogInfo("Request:", "New request", "accept new request")
		content, err := io.ReadAll(r.Body)
		if err != nil {
			h.logger.LogErr(err, "")
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte(err.Error()))
			return
		}
		defer r.Body.Close()

		var u model.User
		if err := json.Unmarshal(content, &u); err != nil {
			h.logger.LogErr(err, "")
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte(err.Error()))
			return
		}
		if u.Name == "" || u.Surname == "" {
			h.logger.LogErr(fmt.Errorf("first or last name is not filled in"), "")
			rw.WriteHeader(http.StatusBadRequest)
			rw.Write([]byte(fmt.Sprintf("first or last name is not filled in")))
			return
		}

		err = h.Service.Collect(u)
		if err != nil {
			h.logger.LogErr(err, "")
			rw.WriteHeader(http.StatusBadRequest)
			rw.Write([]byte(err.Error()))
			return
		}

		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(http.StatusOK)
		rw.Write([]byte("ok"))
	}
}

// GetUsers хендлер получения данных из gauge and counter в формате JSON
func (h *Handler) GetUsers() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		h.logger.LogInfo("Request:", "New request", "accept new request")
		filterOptions := r.Context().Value(filter.OptionsContextKey).(filter.Options)
		name := r.URL.Query().Get("name")
		if name != "" {
			err := filterOptions.AddField(model.Name, filter.OperatorSubString, fmt.Sprintf("'%s'", name), filter.DataTypeStr)
			if err != nil {
				h.logger.LogErr(err, "wrong operator")
				rw.WriteHeader(http.StatusBadRequest)
				rw.Write([]byte("Bad operator"))
				return
			}
		}
		surname := r.URL.Query().Get("surname")
		if surname != "" {
			err := filterOptions.AddField(model.Surname, filter.OperatorSubString, fmt.Sprintf("'%s'", surname), filter.DataTypeStr)
			if err != nil {
				h.logger.LogErr(err, "wrong operator")
				rw.WriteHeader(http.StatusBadRequest)
				rw.Write([]byte("Bad operator"))
				return
			}
		}
		patronymic := r.URL.Query().Get(model.Patronymic)
		if patronymic != "" {
			err := filterOptions.AddField(model.Patronymic, filter.OperatorSubString, fmt.Sprintf("'%s'", patronymic), filter.DataTypeStr)
			if err != nil {
				h.logger.LogErr(err, "wrong operator")
				rw.WriteHeader(http.StatusBadRequest)
				rw.Write([]byte("Bad operator"))
				return
			}
		}
		age := r.URL.Query().Get(model.Age)
		if age != "" {
			operator := filter.OperatorEq
			value := age
			if strings.Index(age, ":") != -1 {
				split := strings.Split(age, ":")
				operator = split[0]
				value = split[1]
			}
			err := filterOptions.AddField(model.Age, operator, value, filter.DataTypeInt)
			if err != nil {
				h.logger.LogErr(err, "wrong operator")
				rw.WriteHeader(http.StatusBadRequest)
				rw.Write([]byte("Bad operator"))
				return
			}
		}
		gender := r.URL.Query().Get(model.Gender)
		if gender != "" {
			err := filterOptions.AddField(model.Gender, filter.OperatorSubString, fmt.Sprintf("'%s'", gender), filter.DataTypeStr)
			if err != nil {
				h.logger.LogErr(err, "wrong operator")
				rw.WriteHeader(http.StatusBadRequest)
				rw.Write([]byte("Bad operator"))
				return
			}
		}
		national := r.URL.Query().Get(model.National)
		if national != "" {
			err := filterOptions.AddField(model.National, filter.OperatorSubString, fmt.Sprintf("'%s'", national), filter.DataTypeStr)
			if err != nil {
				h.logger.LogErr(err, "wrong operator")
				rw.WriteHeader(http.StatusBadRequest)
				rw.Write([]byte("Bad operator"))
				return
			}
		}

		answer, err := h.Service.GetList(filterOptions)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}
		h.logger.LogInfo("Request:", "", "send answer")
		rw.WriteHeader(http.StatusOK)
		rw.Write(answer)
	}
}

// GetUser хендлер получения данных из gauge and counter
func (h *Handler) GetUser() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		//проверка и разбивка URL
		url := strings.Split(r.URL.Path, "/")
		if len(url) < 2 {
			rw.WriteHeader(http.StatusBadRequest)
			rw.Write([]byte("incorrect request"))
			return
		}
		method := url[1]
		if method != "value" {
			rw.WriteHeader(http.StatusNotFound)
			rw.Write([]byte("method is wrong"))
			return
		}
		id := url[2]

		value, err := h.Service.GetByID(id)
		if err != nil {
			h.logger.LogErr(err, "")
			rw.WriteHeader(http.StatusNotFound)
			rw.Write([]byte("incorrect id"))
			return
		}
		rw.WriteHeader(http.StatusOK)
		rw.Write(value)
		return

	}
}

func (h *Handler) UpdateUser() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		h.logger.LogInfo("Request:", "New request", "accept new request")
		content, err := io.ReadAll(r.Body)
		if err != nil {
			h.logger.LogErr(err, "")
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte(err.Error()))
			return
		}
		defer r.Body.Close()

		var u model.User
		if err := json.Unmarshal(content, &u); err != nil {
			h.logger.LogErr(err, "")
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte(err.Error()))
			return
		}

		err = h.Service.UpdateByID(u)
		if err != nil {
			h.logger.LogErr(err, "")
			rw.WriteHeader(http.StatusBadRequest)
			rw.Write([]byte(err.Error()))
			return
		}

		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(http.StatusOK)
		rw.Write([]byte("ok"))
	}
}

// DeleteUser удаление юзера из БД
func (h *Handler) DeleteUser() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		h.logger.LogInfo("Request:", "New request", "accept new request")
		content, err := io.ReadAll(r.Body)
		if err != nil {
			h.logger.LogErr(err, "")
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte(err.Error()))
			return
		}
		defer r.Body.Close()

		var u model.User
		if err := json.Unmarshal(content, &u); err != nil {
			h.logger.LogErr(err, "")
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte(err.Error()))
			return
		}

		err = h.Service.DeleteByID(u)
		if err != nil {
			h.logger.LogErr(err, "")
			rw.WriteHeader(http.StatusBadRequest)
			rw.Write([]byte(err.Error()))
			return
		}

		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(http.StatusOK)
		rw.Write([]byte("ok"))
	}
}

func (h *Handler) PingDB() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		err := h.Service.Ping()
		if err != nil {
			h.logger.LogErr(err, "")
			http.Error(rw, "", http.StatusInternalServerError)
		}
		rw.WriteHeader(http.StatusOK)
	}
}
