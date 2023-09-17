package repositories

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
	"testFIO/cmd/config"
	"testFIO/cmd/loggers"
	"testFIO/internal/storage/model"
	"testFIO/internal/storage/model/filter"
	"testFIO/pkg/client/postgresql"
)

// PGSStore объявление структуры PostgreSQL.
type PGSStore struct {
	client postgresql.Client
	cfg    config.ServerConfig
	logger loggers.Logger
}

// NewPGSStore создание новой базы.
func NewPGSStore(client postgresql.Client, cfg *config.ServerConfig, logger *loggers.Logger) (*PGSStore, error) {
	return &PGSStore{
		client: client,
		cfg:    *cfg,
		logger: *logger,
	}, nil
}

// CollectUser Сохранение юзера в БД.
func (p *PGSStore) CollectUser(u model.User) error {
	q := `INSERT INTO users (firstname, lastname, patronymic, age, gender, nation)
    						VALUES ($1, $2, $3, $4, $5, $6)`
	if _, err := p.client.Exec(context.Background(), q, u.Name, u.Surname, u.Patronymic, u.Age, u.Gender, u.National); err != nil {
		p.logger.LogErr(err, "Failure to insert object into table")
		return err
	}
	return nil
}

// GetUserByID выгрузка юзера.
func (p *PGSStore) GetUserByID(id string) (model.User, error) {
	var u model.User
	q := `SELECT id, firstname, lastname, patronymic, age, gender, nation FROM users WHERE id = $1`
	if err := p.client.QueryRow(context.Background(), q, id).Scan(&u.ID, &u.Name, &u.Surname,
		&u.Patronymic, &u.Age, &u.Gender, &u.National); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			p.logger.LogErr(err, "Failure to select object from table")
			return u, err
		}
		p.logger.LogErr(err, "wrong id")
		return u, fmt.Errorf("missing user %s", id)
	}
	return u, nil
}

// GetAll выгрузка всех юзеров.
func (p *PGSStore) GetAll(options filter.Options) ([]model.User, error) {
	var users []model.User
	q := `SELECT * FROM users`
	rows, err := p.client.Query(context.Background(), q)
	if err != nil {
		p.logger.LogErr(err, "Failure to select object from table")
		return nil, err
	}
	for rows.Next() {
		var u model.User
		err = rows.Scan(&u.ID, &u.Name, &u.Surname, &u.Patronymic, &u.Age, &u.Gender, &u.National)
		if err != nil {
			p.logger.LogErr(err, "Failure to convert object from table")
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

// GetAllFiltered выгрузка юзеров с фильтрами
func (p *PGSStore) GetAllFiltered(options filter.Options) ([]model.User, error) {
	q := "SELECT * FROM users WHERE"
	for i, field := range options.Fields() {
		q += fmt.Sprintf(" %v %v %v", field.Name, field.Operator, field.Value)
		if i == len(options.Fields())-1 {
			q += fmt.Sprintf(" LIMIT %v", options.Limit())
			break
		} else {
			q += fmt.Sprintf(" AND")
		}
	}
	var users []model.User

	rows, err := p.client.Query(context.Background(), q)
	if err != nil {
		p.logger.LogErr(err, "Failure to select object from table")
		return nil, err
	}
	for rows.Next() {
		var u model.User
		err = rows.Scan(&u.ID, &u.Name, &u.Surname, &u.Patronymic, &u.Age, &u.Gender, &u.National)
		if err != nil {
			p.logger.LogErr(err, "Failure to convert object from table")
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

// ChangeUser изменение юзера.
func (p *PGSStore) ChangeUser(u model.User) error {
	q := `UPDATE users SET firstname = $1, lastname = $2, patronymic = $3, 
    							age = $4, gender = $5, nation = $6 WHERE id = $7`
	if _, err := p.client.Exec(context.Background(), q, u.Name, u.Surname, u.Patronymic,
		u.Age, u.Gender, u.National, u.ID); err != nil {
		p.logger.LogErr(err, "Failure to update object into table")
		return err
	}
	return nil
}

// DeleteUser удаление юзера.
func (p *PGSStore) DeleteUser(u model.User) error {
	q := `DELETE FROM users WHERE id = $1`
	if _, err := p.client.Exec(context.Background(), q, u.ID); err != nil {
		p.logger.LogErr(err, "Failure to delete object from table")
		return err
	}
	return nil
}

func (p *PGSStore) PingClient() error {
	return p.client.Ping(context.Background())
}
