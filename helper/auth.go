package helper

import (
	"api/config"
	"api/models"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
)

func CreateToken(token *models.Auth) (err error) {
	token.AccessUUID = uuid.New().String()
	token.RefreshUUID = uuid.New().String()

	token.AccessExpire = time.Now().Add(time.Minute * 15).Unix()
	token.RefreshExpire = time.Now().Add(time.Hour * 24 * 7).Unix()

	access := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt.MapClaims{
		"authorized":  true,
		"access_uuid": token.AccessUUID,
		"user_id":     config.ENV.ID,
		"expire":      token.AccessExpire,
	})

	token.AccessToken, err = access.SignedString([]byte(config.ENV.AccessSecret))
	if err != nil {
		return
	}

	refresh := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"refresh_uuid": token.RefreshUUID,
		"user_id":      config.ENV.ID,
		"expire":       token.RefreshExpire,
	})

	token.RefreshToken, err = refresh.SignedString([]byte(config.ENV.RefreshSecret))
	return
}

func CheckToken(dto *models.TokenDto) (*models.Auth, error) {
	token, err := jwt.Parse(dto.RefreshToken, func(t *jwt.Token) (interface{}, error) {
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

	// var userUUID string
	// var refreshUUID string
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

	return &models.Auth{AccessUUID: claims["user_id"].(string), RefreshUUID: claims["refresh_uuid"].(string)}, nil
}

func ValidateStr(str1 string, str2 string) (equal bool) {
	var len1 = len(str1)
	var len2 = len(str2)

	var max = len2
	if len1 > len2 {
		max = len1
	}

	equal = true
	for i := 0; i < max; i++ {
		if i >= len1 || i >= len2 || str1[i] != str2[i] {
			equal = false
		}
	}

	return
}

func BotToken() (string, string) {
	hasher := md5.New()
	salt := uuid.New()
	hasher.Write([]byte(salt.String() + config.ENV.BotKey))
	return salt.String(), hex.EncodeToString(hasher.Sum(nil))
}

func GetToken() string {
	hasher := md5.New()
	salt := uuid.New()
	hasher.Write([]byte(salt.String()))
	return hex.EncodeToString(hasher.Sum(nil))
}

func HashSecret(secret string) string {
	hasher := md5.New()
	salt := uuid.New()
	hasher.Write([]byte(salt.String() + config.ENV.Pepper + secret))
	return salt.String() + "$" + hex.EncodeToString(hasher.Sum(nil))
}
