package helper

import (
	"api/config"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
)

func CreateToken() (map[string]interface{}, error) {
	var err error
	token := make(map[string]interface{})

	token["access_uuid"] = uuid.New().String()
	token["refresh_uuid"] = uuid.New().String()

	token["access_expire"] = time.Now().Add(time.Minute * 15).Unix()
	token["refresh_expire"] = time.Now().Add(time.Hour * 24 * 7).Unix()

	access := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt.MapClaims{
		"authorized":  true,
		"access_uuid": token["access_uuid"],
		"user_id":     config.ENV.ID,
		"expire":      token["access_expire"],
	})

	token["access_token"], err = access.SignedString([]byte(config.ENV.AccessSecret))
	if err != nil {
		return nil, err
	}

	refresh := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"refresh_uuid": token["refresh_uuid"],
		"user_id":      config.ENV.ID,
		"expire":       token["refresh_expire"],
	})

	token["refresh_token"], err = refresh.SignedString([]byte(config.ENV.RefreshSecret))
	return token, err
}

func CheckToken(refreshToken string) (map[string]interface{}, error) {
	token, err := jwt.Parse(refreshToken, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("invalid signing method")
		}

		if _, ok := t.Claims.(jwt.Claims); !ok && !t.Valid {
			return nil, fmt.Errorf("expired token")
		}

		return []byte(config.ENV.RefreshSecret), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("Unauthorized token")
	}

	if _, ok = claims["refresh_uuid"].(string); !ok {
		return nil, fmt.Errorf("Invalid token inforamation")
	}

	if _, ok = claims["user_id"].(string); !ok {
		return nil, fmt.Errorf("Invalid token inforamation")
	}

	return map[string]interface{}{"access_uuid": claims["user_id"].(string), "refresh_uuid": claims["refresh_uuid"].(string)}, nil
}

func BotToken() (string, string) {
	hasher := md5.New()
	salt := uuid.New()
	hasher.Write([]byte(salt.String() + config.ENV.BotKey))
	return salt.String(), hex.EncodeToString(hasher.Sum(nil))
}

func GetToken() string {
	hasher := md5.New()
	hasher.Write([]byte(uuid.New().String()))
	return hex.EncodeToString(hasher.Sum(nil))
}

func HashSecret(secret string) string {
	hasher := md5.New()
	salt := uuid.New().String()
	hasher.Write([]byte(salt + config.ENV.Pepper + secret))
	return salt + "$" + hex.EncodeToString(hasher.Sum(nil))
}
