package storage

import (
	"testFIO/internal/storage/model"
	"testFIO/internal/storage/model/filter"
)

// Storage интерфейс хранилища.
type Storage interface {
	CollectUser(u model.User) error                              //получение юзеров.
	GetUserByID(id string) (model.User, error)                   //выгрузка юзера
	GetAll(options filter.Options) ([]model.User, error)         //выгрузка всех юзеров.
	GetAllFiltered(options filter.Options) ([]model.User, error) //выгрузка юзеров с фильтрами
	ChangeUser(u model.User) error                               //изменение юзера
	DeleteUser(u model.User) error                               //удаление юзера
	PingClient() error                                           //ping
}
