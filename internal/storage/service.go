package storage

import (
	"encoding/json"

	"testFIO/cmd/loggers"
	"testFIO/internal/storage/model"
	"testFIO/internal/storage/model/filter"
)

type Service struct {
	storage Storage
	logger  loggers.Logger
}

func NewService(storage Storage) *Service {
	return &Service{
		storage: storage,
	}
}
func (s *Service) Collect(u model.User) error {
	return s.storage.CollectUser(u)
}

func (s *Service) GetList(options filter.Options) ([]byte, error) {
	if options.Fields() != nil {
		users, err := s.storage.GetAllFiltered(options)
		if err != nil {
			s.logger.LogErr(err, "failed to get all users from db")
			return nil, err
		}
		usersJSON, err := json.Marshal(users)
		if err != nil {
			s.logger.LogErr(err, "failed to convert to JSON")
			return nil, err
		}
		return usersJSON, nil
	}
	users, err := s.storage.GetAll(options)
	if err != nil {
		s.logger.LogErr(err, "failed to get all users from db")
		return nil, err
	}
	usersJSON, err := json.Marshal(users)
	if err != nil {
		s.logger.LogErr(err, "failed to convert to JSON")
		return nil, err
	}
	return usersJSON, nil
}

func (s *Service) GetByID(id string) ([]byte, error) {
	u, err := s.storage.GetUserByID(id)
	if err != nil {
		s.logger.LogErr(err, "failed to get all users from db")
		return nil, err
	}
	usersJSON, err := json.Marshal(u)
	if err != nil {
		s.logger.LogErr(err, "failed to convert to JSON")
		return nil, err
	}
	return usersJSON, nil
}

func (s *Service) UpdateByID(u model.User) error {
	return s.storage.ChangeUser(u)
}
func (s *Service) DeleteByID(u model.User) error {
	return s.storage.DeleteUser(u)
}
func (s *Service) Ping() error {
	return s.storage.PingClient()
}
