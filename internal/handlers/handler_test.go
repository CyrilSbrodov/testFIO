package handlers_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"testFIO/cmd/config"
	"testFIO/cmd/loggers"
	"testFIO/internal/handlers"
	"testFIO/internal/storage"
	"testFIO/internal/storage/model"
	"testFIO/internal/storage/repositories"
)

var (
	CFG config.ServerConfig
)

func TestMain(m *testing.M) {
	cfg := config.ServerConfigInit()
	CFG = *cfg
	os.Exit(m.Run())
}

func TestHandler_CollectUsers(t *testing.T) {
	type want struct {
		statusCode int
	}
	logger := loggers.NewLogger()
	repo, _ := repositories.NewRepository(&CFG, logger)
	service := storage.NewService(repo)

	type fields struct {
		Service storage.Service
	}
	tests := []struct {
		name    string
		fields  fields
		want    want
		request string
		req     model.User
	}{
		{
			name: "Test ok",
			fields: fields{
				*service,
			},
			request: "http://localhost:8080/update",
			want: want{
				200,
			},
			req: model.User{
				ID:      "test",
				Name:    "test01",
				Surname: "test02",
			},
		},
		{
			name: "Test ok",
			fields: fields{
				*service,
			},
			request: "http://localhost:8080/add/",
			want: want{
				200,
			},
			req: model.User{
				ID:      "test",
				Name:    "test01",
				Surname: "test02",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			usersJSON, err := json.Marshal(tt.req)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			request := httptest.NewRequest(http.MethodPost, tt.request, bytes.NewBuffer(usersJSON))
			w := httptest.NewRecorder()
			h := handlers.Handler{
				Service: *service,
			}
			h.CollectUsers().ServeHTTP(w, request)
			result := w.Result()
			defer result.Body.Close()
			assert.Equal(t, tt.want.statusCode, result.StatusCode)
		})
	}
}
