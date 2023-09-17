package model

import (
	"sync"
)

const (
	Name       = "firstname"
	Surname    = "lastname"
	Patronymic = "patronymic"
	Age        = "age"
	Gender     = "gender"
	National   = "national"
)

// User структура получаемой и хранимой сущности.
type User struct {
	ID         string
	Name       string `json:"name"`                 // имя
	Surname    string `json:"surname"`              // фамилия
	Patronymic string `json:"patronymic,omitempty"` // отчество (необязательное поле)
	Age        int    `json:"age,omitempty"`        // возраст
	Gender     string `json:"gender,omitempty"`     // пол
	National   string `json:"national,omitempty"`   // национальность
	Error      error  `json:"value,omitempty"`      // поле ошибки
}

// UsersStore инициализация временного хранилища.
var UsersStore = map[string]User{}

// CacheUsers структура юзеров для хранения.
type CacheUsers struct {
	Store map[string]User
	Sync  sync.Mutex
}

// NewCacheUsers инициализация нового временного хранилища.
func NewCacheUsers() *CacheUsers {
	return &CacheUsers{
		Store: make(map[string]User),
	}
}

type UserApiAge struct {
	Count int    `json:"count"`
	Name  string `json:"name"`
	Age   int    `json:"age"`
}

type UserApiGender struct {
	Count       int     `json:"count"`
	Name        string  `json:"name"`
	Gender      string  `json:"gender"`
	Probability float64 `json:"probability"`
}

type UserApiNational struct {
	Count   int       `json:"count"`
	Name    string    `json:"name"`
	Country []Country `json:"country"`
}
type Country struct {
	CountryID   string  `json:"country_id"`
	Probability float64 `json:"probability"`
}
