package api

import (
	"bytes"
	"context"
	"fmt"
	"iContext/internal/models"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/BurntSushi/toml"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

var (
	configPath = "../../configs/api.toml"
	config     = NewConfig()
	_, _       = toml.DecodeFile(configPath, config)
	redisAddr  = "localhost:6379"
	server     = New(config, redisAddr)
	_          = server.configureStorageField()
	_          = server.configureRedisField()
	ctx        = context.Background()
)

type TestCase struct {
	name                 string
	inputBody            string
	expectedStatusCode   int
	expectedResponseBody string
}

func TestAPI_SaveUser(t *testing.T) {
	// truncate a table and restart the primary key value
	server.storage.Exec("TRUNCATE users RESTART IDENTITY;")

	tests := []TestCase{
		{
			name:                 "Ok",
			inputBody:            `{"name": "Alex","age": 21}`,
			expectedStatusCode:   201,
			expectedResponseBody: fmt.Sprintf("{\n    \"id\": %d\n}", 1),
		},

		{
			name:                 "Wrong name field type",
			inputBody:            `{"name": 1234,"age": 21}`,
			expectedStatusCode:   400,
			expectedResponseBody: `{"message":"Request data is invalid"}`,
		},
		{
			name:                 "Wrong age field type",
			inputBody:            `{"name": "Alex","age": "21"}`,
			expectedStatusCode:   400,
			expectedResponseBody: `{"message":"Request data is invalid"}`,
		},
		{
			name:                 "Empty body",
			inputBody:            "",
			expectedStatusCode:   400,
			expectedResponseBody: `{"message":"Request data is invalid"}`,
		},

		{
			name:                 "Empty age field",
			inputBody:            `{"name": "Alex"}`,
			expectedStatusCode:   201,
			expectedResponseBody: fmt.Sprintf("{\n    \"id\": %d\n}", 2),
		},
		{
			name:                 "Empty name field",
			inputBody:            `{"age": 21}`,
			expectedStatusCode:   400,
			expectedResponseBody: `{"message":"Request data is invalid"}`,
		},
		{
			name:                 "Wrong body keys",
			inputBody:            `{"wrongkey": "Alex","thiswrong2": 21}`,
			expectedStatusCode:   400,
			expectedResponseBody: `{"message":"Request data is invalid"}`,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			runTest(t, http.MethodPost, test, "/postgres/users", server.SaveUser)
		})
	}
}

func TestAPI_redisIncr(t *testing.T) {
	setDB := models.RedisJSON{Key: "age", Value: 100}
	if err := server.client.Set(ctx, setDB); err != nil {
		t.Log("err:", err)
	}

	tests := []TestCase{
		{
			name:                 "Ok",
			inputBody:            `{"key": "age","value": 19}`,
			expectedStatusCode:   200,
			expectedResponseBody: "{\n    \"value\": 119\n}",
		},
		{
			name:                 "Key not exists",
			inputBody:            `{"key": "suchkey100%%doesnotexist","value": 25}`,
			expectedStatusCode:   400,
			expectedResponseBody: `{"message":"Request data is invalid"}`,
		},
		{
			name:                 "Empty body",
			inputBody:            "",
			expectedStatusCode:   400,
			expectedResponseBody: `{"message":"Request data is invalid"}`,
		},
		{
			name:                 "Missed value",
			inputBody:            `{"key": "age"}`,
			expectedStatusCode:   400,
			expectedResponseBody: `{"message":"Request data is invalid"}`,
		},
		{
			name:                 "Missed key",
			inputBody:            `{"value": 19}`,
			expectedStatusCode:   400,
			expectedResponseBody: `{"message":"Request data is invalid"}`,
		},
		{
			name:                 "Another value type",
			inputBody:            `{"key": "age","value": "19"}`,
			expectedStatusCode:   400,
			expectedResponseBody: `{"message":"Request data is invalid"}`,
		},
		{
			name:                 "Another key type",
			inputBody:            `{"key": 12345,"value": 19}`,
			expectedStatusCode:   400,
			expectedResponseBody: `{"message":"Request data is invalid"}`,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			runTest(t, http.MethodPost, test, "/redis/incr", server.redisIncr)
		})
	}
}

func TestAPI_cryptHmac(t *testing.T) {
	tests := []TestCase{
		{
			name:                 "Ok",
			inputBody:            `{"text": "test","key": "test123"}`,
			expectedStatusCode:   200,
			expectedResponseBody: `"b596e24739fd44d42ffd25f26ea367dad3a71f61c8c5fab6b6ee6ceeae5a7170b66445d6eaadfb49e6d4e968a2888726ff522e3bf065c966aa66a24153778382"`,
		},
		{
			name:                 "Missed Key",
			inputBody:            `{"text": "test"}`,
			expectedStatusCode:   400,
			expectedResponseBody: `{"message":"Request data is invalid"}`,
		},
		{
			name:                 "Missed Text",
			inputBody:            `{"key": "test123"}`,
			expectedStatusCode:   400,
			expectedResponseBody: `{"message":"Request data is invalid"}`,
		},
		{
			name:                 "Wrong key type",
			inputBody:            `{"text": "test","key": 123}`,
			expectedStatusCode:   400,
			expectedResponseBody: `{"message":"Request data is invalid"}`,
		},
		{
			name:                 "Wrong text type",
			inputBody:            `{"text": 12455,"key": "test123"}`,
			expectedStatusCode:   400,
			expectedResponseBody: `{"message":"Request data is invalid"}`,
		},
		{
			name:                 "Empty body",
			inputBody:            "",
			expectedStatusCode:   400,
			expectedResponseBody: `{"message":"Request data is invalid"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			runTest(t, http.MethodPost, test, "/sign/hmacsha512", server.cryptHmac)
		})
	}
}

func runTest(t *testing.T, method string, test TestCase, endpoint string, handlerFunc func(c *gin.Context)) {
	router := gin.Default()
	switch method {
	case "POST":
		router.POST(endpoint, handlerFunc)
	case "GET":
		router.GET(endpoint, handlerFunc)
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(method, endpoint, bytes.NewBufferString(test.inputBody))

	router.ServeHTTP(w, req)

	assert.Equal(t, w.Code, test.expectedStatusCode)
	assert.Equal(t, w.Body.String(), test.expectedResponseBody)
}
