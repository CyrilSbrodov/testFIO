package repositories

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"testFIO/cmd/config"
	"testFIO/internal/storage/model"
	"testFIO/internal/storage/model/filter"
)

var (
	CFG config.ServerConfig
)

func TestMain(m *testing.M) {
	cfg := config.ServerConfigInit()
	CFG = *cfg
	CFG.DatabaseDSN = "postgres://postgres:postgres@postgres:5432/praktikum?sslmode=disable"
	os.Exit(m.Run())
}

func TestPGSStore_CollectUser(t *testing.T) {
	s, teardown := TestPGStore(t, CFG)
	defer teardown("users")
	err := s.CollectUser(model.User{
		ID:         "test01",
		Name:       "test02",
		Surname:    "test03",
		Patronymic: "test04",
		Age:        1,
		Gender:     "test05",
		National:   "test06",
		Error:      nil,
	})
	assert.NoError(t, err)

}

func TestPGSStore_ChangeUser(t *testing.T) {
	s, teardown := TestPGStore(t, CFG)
	defer teardown("users")
	err := s.CollectUser(model.User{
		ID:         "1",
		Name:       "test02",
		Surname:    "test03",
		Patronymic: "test04",
		Age:        1,
		Gender:     "test05",
		National:   "test06",
		Error:      nil,
	})
	assert.NoError(t, err)
	err = s.ChangeUser(model.User{
		ID:         "1",
		Name:       "T",
		Surname:    "E",
		Patronymic: "S",
		Age:        0,
		Gender:     "T",
		National:   "",
		Error:      nil,
	})
	assert.NoError(t, err)
}

func TestPGSStore_GetAll(t *testing.T) {
	s, teardown := TestPGStore(t, CFG)
	defer teardown("users")
	err := s.CollectUser(model.User{
		ID:         "2",
		Name:       "3",
		Surname:    "1",
		Patronymic: "2",
		Age:        2,
		Gender:     "",
		National:   "",
		Error:      nil,
	})
	assert.NoError(t, err)
	assert.NoError(t, err)
	var options filter.Options
	m, err := s.GetAll(options)
	assert.NotNil(t, m)
	assert.NoError(t, err)
}

func TestPGSStore_DeleteUser(t *testing.T) {
	s, teardown := TestPGStore(t, CFG)
	defer teardown("users")
	err := s.CollectUser(model.User{
		ID:         "2",
		Name:       "3",
		Surname:    "1",
		Patronymic: "2",
		Age:        2,
		Gender:     "",
		National:   "",
		Error:      nil,
	})
	assert.NoError(t, err)
	err = s.DeleteUser(model.User{
		ID:         "2",
		Name:       "3",
		Surname:    "1",
		Patronymic: "2",
		Age:        2,
		Gender:     "",
		National:   "",
		Error:      nil,
	})
	assert.NoError(t, err)
}

func TestPGSStore_PingClient(t *testing.T) {
	s, teardown := TestPGStore(t, CFG)
	defer teardown("users")
	err := s.PingClient()
	assert.NoError(t, err)
}
