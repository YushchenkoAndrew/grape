package helper

func CreateToken() (map[string]interface{}, error) {
	// var err error
	// token := make(map[string]interface{})

	// return token, err

	return map[string]interface{}{}, nil
}

func CheckToken(refreshToken string) (map[string]interface{}, error) {
	// token, err := jwt.Parse(refreshToken, func(t *jwt.Token) (interface{}, error) {
	// 	if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
	// 		return nil, fmt.Errorf("invalid signing method")
	// 	}

	// 	if _, ok := t.Claims.(jwt.Claims); !ok && !t.Valid {
	// 		return nil, fmt.Errorf("expired token")
	// 	}

	// 	return []byte(config.ENV.RefreshSecret), nil
	// })

	// if err != nil {
	// 	return nil, err
	// }

	// claims, ok := token.Claims.(jwt.MapClaims)
	// if !ok || !token.Valid {
	// 	return nil, fmt.Errorf("Unauthorized token")
	// }

	// if _, ok = claims["refresh_uuid"].(string); !ok {
	// 	return nil, fmt.Errorf("Invalid token inforamation")
	// }

	// if _, ok = claims["user_id"].(string); !ok {
	// 	return nil, fmt.Errorf("Invalid token inforamation")
	// }

	// return map[string]interface{}{"access_uuid": claims["user_id"].(string), "refresh_uuid": claims["refresh_uuid"].(string)}, nil
	return map[string]interface{}{}, nil
}

// func BotToken() (string, string) {
// 	hasher := md5.New()
// 	salt := uuid.New()
// 	hasher.Write([]byte(salt.String() + config.ENV.BotKey))
// 	return salt.String(), hex.EncodeToString(hasher.Sum(nil))
// }

// func GetToken() string {
// 	hasher := md5.New()
// 	hasher.Write([]byte(uuid.New().String()))
// 	return hex.EncodeToString(hasher.Sum(nil))
// }

// func HashSecret(secret string) string {
// 	hasher := md5.New()
// 	salt := uuid.New().String()
// 	hasher.Write([]byte(salt + config.ENV.Pepper + secret))
// 	return salt + "$" + hex.EncodeToString(hasher.Sum(nil))
// }
