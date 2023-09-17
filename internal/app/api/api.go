package api

import (
	"encoding/json"
	"io"
	"net/http"

	"testFIO/cmd/config"
	"testFIO/cmd/loggers"
	"testFIO/internal/storage/model"
)

type External interface {
	AddGender(u model.User) (model.User, error)
	AddNational(u model.User) (model.User, error)
	AddAge(u model.User) (model.User, error)
}

type ExternalApi struct {
	client *http.Client
	cfg    config.ServerConfig
	logger loggers.Logger
}

func NewExternalApi(cfg config.ServerConfig, logger loggers.Logger) *ExternalApi {
	client := &http.Client{}
	return &ExternalApi{
		client,
		cfg,
		logger,
	}
}

func (e *ExternalApi) AddGender(u model.User) (model.User, error) {
	req, err := http.NewRequest(http.MethodGet, e.cfg.AddrGender+u.Name, nil)
	if err != nil {
		e.logger.LogErr(err, "Failed to request")
		return u, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	resp, err := e.client.Do(req)
	if err != nil {
		e.logger.LogErr(err, "Failed to do request")
		return u, err
	}
	var gender model.UserApiGender
	r, err := io.ReadAll(resp.Body)
	if err != nil {
		e.logger.LogErr(err, "Failed to read body")
		return u, err
	}
	if err = json.Unmarshal(r, &gender); err != nil {
		e.logger.LogErr(err, "Failed to read body")
		return u, err
	}
	u.Gender = gender.Gender
	resp.Body.Close()
	return u, nil
}

func (e *ExternalApi) AddNational(u model.User) (model.User, error) {
	req, err := http.NewRequest(http.MethodGet, e.cfg.AddrNational+u.Name, nil)
	if err != nil {
		e.logger.LogErr(err, "Failed to request")
		return u, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	resp, err := e.client.Do(req)
	if err != nil {
		e.logger.LogErr(err, "Failed to do request")
		return u, err
	}
	var national model.UserApiNational
	r, err := io.ReadAll(resp.Body)
	if err != nil {
		e.logger.LogErr(err, "Failed to read body")
		return u, err
	}
	if err = json.Unmarshal(r, &national); err != nil {
		e.logger.LogErr(err, "Failed to read body")
		return u, err
	}
	var max float64
	for i := 0; i < len(national.Country); i++ {
		if national.Country[i].Probability > max {
			max = national.Country[i].Probability
		}
	}
	for i := 0; i < len(national.Country); i++ {
		if national.Country[i].Probability == max {
			u.National = national.Country[i].CountryID
			break
		}
	}
	resp.Body.Close()
	return u, nil
}
func (e *ExternalApi) AddAge(u model.User) (model.User, error) {
	req, err := http.NewRequest(http.MethodGet, e.cfg.AddrAge+u.Name, nil)
	if err != nil {
		e.logger.LogErr(err, "Failed to request")
		return u, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	resp, err := e.client.Do(req)
	if err != nil {
		e.logger.LogErr(err, "Failed to do request")
		return u, err
	}
	var age model.UserApiAge
	r, err := io.ReadAll(resp.Body)
	if err != nil {
		e.logger.LogErr(err, "Failed to read body")
		return u, err
	}
	if err = json.Unmarshal(r, &age); err != nil {
		e.logger.LogErr(err, "Failed to read body")
		return u, err
	}
	u.Age = age.Age
	resp.Body.Close()
	return u, nil
}
