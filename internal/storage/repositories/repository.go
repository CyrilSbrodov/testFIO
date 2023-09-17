package repositories

import (
	"database/sql"
	"fmt"
	"sync"

	"testFIO/cmd/config"
	"testFIO/cmd/loggers"
	"testFIO/internal/storage/model"
	"testFIO/internal/storage/model/filter"
)

// Repository структура репозитория.
type Repository struct {
	Users  map[string]model.User
	Dsn    string
	logger loggers.Logger
	sync   sync.Mutex
}

func (r *Repository) GetAll(options filter.Options) ([]model.User, error) {
	var users []model.User
	if options.Fields() == nil {
		for _, user := range r.Users {
			users = append(users, user)
		}
		return users, nil
	}
	return users, nil
}
func (r *Repository) GetAllFiltered(options filter.Options) ([]model.User, error) {
	var users []model.User
	if options.Fields() == nil {
		for _, user := range r.Users {
			users = append(users, user)
		}
		return users, nil
	}
	return users, nil
}
func (r *Repository) GetUserByID(id string) (model.User, error) {
	user, ok := r.Users[id]
	if !ok {
		r.logger.LogErr(fmt.Errorf("id not found %s", id), "id not found")
		err := fmt.Errorf("id not found %s", id)
		return user, err
	}
	return user, nil
}
func (r *Repository) ChangeUser(u model.User) error {
	_, ok := r.Users[u.ID]
	if !ok {
		err := fmt.Errorf("user not found")
		r.logger.LogErr(err, "user not found")
		return err
	}
	r.Users[u.ID] = u
	return nil
}
func (r *Repository) DeleteUser(u model.User) error {
	_, ok := r.Users[u.ID]
	if !ok {
		err := fmt.Errorf("user not found")
		r.logger.LogErr(err, "user not found")
		return err
	}
	delete(r.Users, u.ID)
	return nil
}

// NewRepository создание нового репозитория.
func NewRepository(cfg *config.ServerConfig, logger *loggers.Logger) (*Repository, error) {
	users := model.UsersStore

	return &Repository{
		Users: users,
		Dsn:   cfg.DatabaseDSN,
	}, nil
}

// CollectUser сохранение метрики.
func (r *Repository) CollectUser(u model.User) error {

	entry, ok := r.Users[u.ID]
	if !ok {
		r.Users[u.ID] = u
		return nil
	}
	r.Users[u.ID] = entry
	return nil

}

// PingClient проверка клиента.
func (r *Repository) PingClient() error {
	db, err := sql.Open("postgres", r.Dsn)
	if err != nil {
		r.logger.LogErr(err, "not connection")
		return err
	}

	return db.Ping()
}
