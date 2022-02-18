package info_test

import (
	"api/config"
	"api/interfaces"
	"api/models"
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"gopkg.in/stretchr/testify.v1/require"
)

const (
	SALT = "54"
)

var token string
var client = &http.Client{}

var dto = models.InfoDto{
	Countries: "KZ",
	Views:     19,
	Visitors:  2,
}

func TestInfoCreate(t *testing.T) {
	body, err := json.Marshal(dto)
	require.NoError(t, err)

	var req *http.Request
	req, err = http.NewRequest("POST", fmt.Sprintf("http://127.0.0.1:%s/api/info", config.ENV.Port), bytes.NewBuffer(body))
	require.NoError(t, err)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bear %s", token))

	var res *http.Response
	res, err = client.Do(req)

	require.NoError(t, err)
	defer res.Body.Close()

	require.Equal(t, res.StatusCode, http.StatusCreated)

	var result models.Success
	err = json.NewDecoder(res.Body).Decode(&result)
	require.NoError(t, err)

	require.Equal(t, result.Status, "OK")
	require.NotEmpty(t, result.Result)

	model := result.Result.([]interface{})[0].(map[string]interface{})

	fmt.Println(model)
	require.Equal(t, model["countries"], dto.Countries)
	require.Equal(t, uint16(model["media"].(float64)), uint16(0))
	require.Equal(t, uint16(model["clicks"].(float64)), uint16(0))
	require.Equal(t, uint16(model["views"].(float64)), dto.Views)
	require.Equal(t, uint16(model["visitors"].(float64)), dto.Visitors)
}

func TestInfoGetAll(t *testing.T) {
	res, err := http.Get(fmt.Sprintf("http://127.0.0.1:%s/api/info?countries=UK,KZ", config.ENV.Port))
	require.NoError(t, err)

	defer res.Body.Close()

	var result models.Success
	err = json.NewDecoder(res.Body).Decode(&result)
	require.NoError(t, err)

	require.Equal(t, result.Status, "OK")
	require.NotEmpty(t, result.Result)

	model := result.Result.([]interface{})[0].(map[string]interface{})

	fmt.Println(model)
	require.Equal(t, model["countries"], dto.Countries)
	require.Equal(t, uint16(model["media"].(float64)), uint16(0))
	require.Equal(t, uint16(model["clicks"].(float64)), uint16(0))
	require.Equal(t, uint16(model["views"].(float64)), dto.Views)
	require.Equal(t, uint16(model["visitors"].(float64)), dto.Visitors)
}

func TestInfoDelete(t *testing.T) {
	// TODO:
}

func init() {
	config.NewConfig([]func() interfaces.Config{
		config.NewEnvConfig("./"),
	}).Init()

	hasher := md5.New()
	hasher.Write([]byte(SALT + config.ENV.Pepper + config.ENV.Pass))
	body, err := json.Marshal(&models.LoginDto{
		User: config.ENV.User,
		Pass: fmt.Sprintf("%s$%s", SALT, hex.EncodeToString(hasher.Sum(nil))),
	})

	if err != nil {
		panic(err)
	}

	var res *http.Response
	res, err = http.Post(fmt.Sprintf("http://127.0.0.1:%s/api/login", config.ENV.Port), "application/json", bytes.NewBuffer(body))
	if err != nil {
		panic(err)
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		panic(fmt.Errorf("Auth FAILED"))
	}

	var entity = models.TokenEntity{}
	if err = json.NewDecoder(res.Body).Decode(&entity); err != nil {
		panic(err)
	}

	token = entity.AccessToken
}
