package api

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"iContext/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (api *API) redisIncr(c *gin.Context) {
	var redisRequest models.RedisJSON
	if err := c.BindJSON(&redisRequest); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "Request data is invalid")
		return
	}

	if redisRequest.Key == "" || redisRequest.Value == 0 {
		newErrorResponse(c, http.StatusBadRequest, "Request data is invalid")
		return
	}

	redisInDB, err := api.client.Get(c, redisRequest.Key)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "Request data is invalid")
		return
	}

	redisInDB.Value += redisRequest.Value
	api.client.Set(c, redisInDB)

	c.IndentedJSON(http.StatusOK, gin.H{"value": redisInDB.Value})
}

func (api *API) cryptHmac(c *gin.Context) {
	var cryptData models.CryptPair
	if err := c.BindJSON(&cryptData); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "Request data is invalid")
		return
	}
	if cryptData.Text == "" || cryptData.Key == "" {
		newErrorResponse(c, http.StatusBadRequest, "Request data is invalid")
		return
	}
	sigStr := SignHmacsha512(cryptData)
	c.IndentedJSON(http.StatusOK, sigStr)
}

func SignHmacsha512(cryptPair models.CryptPair) string {
	sighash := hmac.New(sha512.New, []byte(cryptPair.Key))
	sighash.Write([]byte(cryptPair.Text))
	sigstr := hex.EncodeToString(sighash.Sum(nil))
	return sigstr
}

func (api *API) SaveUser(c *gin.Context) {
	var user models.User
	if err := c.BindJSON(&user); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "Request data is invalid")
		return
	}
	if user.Name == "" {
		newErrorResponse(c, http.StatusBadRequest, "Request data is invalid")
		return
	}

	id, err := api.storage.CreateUser(&user)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.IndentedJSON(http.StatusCreated, gin.H{"id": id})
}
