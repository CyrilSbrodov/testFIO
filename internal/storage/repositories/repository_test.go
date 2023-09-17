package repositories

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"testFIO/internal/storage/model"
)

func TestRepository_CollectUser(t *testing.T) {
	repo, _ := NewRepository(&CFG, nil)
	type fields struct {
		repo *Repository
	}
	type args struct {
		u model.User
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "test counter OK",
			fields: fields{
				repo: repo,
			},
			args: args{u: model.User{
				ID:      "test",
				Name:    "test01",
				Surname: "test02",
				Age:     10,
			}},
		},
		{
			name: "test collect OK",
			fields: fields{
				repo: repo,
			},
			args: args{u: model.User{
				ID:      "test",
				Name:    "test01",
				Surname: "test02",
				Age:     10,
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.fields.repo.CollectUser(tt.args.u)
			assert.NoError(t, err)
		})
	}
}

func TestRepository_ChangeUser(t *testing.T) {
	repo, _ := NewRepository(&CFG, nil)
	repo.Users = nil
	repo.Users = make(map[string]model.User)
	repo.Users["test"] = model.User{
		ID:         "test",
		Name:       "test01",
		Surname:    "test02",
		Patronymic: "",
		Age:        10,
		Gender:     "",
		National:   "",
		Error:      nil,
	}
	type fields struct {
		repo *Repository
	}

	type args struct {
		u model.User
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "test update OK",
			fields: fields{
				repo: repo,
			},
			args: args{u: model.User{
				ID:         "test",
				Name:       "1",
				Surname:    "2",
				Patronymic: "",
				Age:        3,
				Gender:     "",
				National:   "",
				Error:      nil,
			},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.fields.repo.ChangeUser(tt.args.u)
			assert.NoError(t, err)
			u, err := tt.fields.repo.GetUserByID(tt.args.u.ID)
			assert.NoError(t, err)
			assert.Equal(t, u, tt.args.u)
		})
	}
}

func TestRepository_DeleteUser(t *testing.T) {
	repo, _ := NewRepository(&CFG, nil)
	repo.Users = nil
	repo.Users = make(map[string]model.User)
	repo.Users["test"] = model.User{
		ID:         "test",
		Name:       "test01",
		Surname:    "test02",
		Patronymic: "",
		Age:        10,
		Gender:     "",
		National:   "",
		Error:      nil,
	}
	type fields struct {
		repo *Repository
	}
	type args struct {
		u model.User
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "test OK",
			fields: fields{
				repo: repo,
			},
			args: args{u: model.User{
				ID: "test",
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.fields.repo.DeleteUser(tt.args.u)
			assert.NoError(t, err)
		})
	}
}

func TestRepository_GetUserByID(t *testing.T) {
	repo, _ := NewRepository(&CFG, nil)
	repo.Users = nil
	repo.Users = make(map[string]model.User)
	repo.Users["test"] = model.User{
		ID:         "test",
		Name:       "test01",
		Surname:    "test02",
		Patronymic: "",
		Age:        10,
		Gender:     "",
		National:   "",
		Error:      nil,
	}
	type fields struct {
		repo *Repository
	}
	type args struct {
		u model.User
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "test OK",
			fields: fields{
				repo: repo,
			},
			args: args{u: model.User{
				ID:         "test",
				Name:       "test01",
				Surname:    "test02",
				Patronymic: "",
				Age:        10,
				Gender:     "",
				National:   "",
				Error:      nil,
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u, err := tt.fields.repo.GetUserByID(tt.args.u.ID)
			assert.NoError(t, err)
			assert.Equal(t, tt.args.u, u)
		})
	}
}
