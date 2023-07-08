package jwt

import (
	"log"
	"reflect"

	"github.com/bytedance/sonic"
	"github.com/golang-jwt/jwt/v5"
)

type User struct {
	Id        string `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Verified  bool   `json:"verified"`
	Role      string `json:"role"`
	CreatedAt int64  `json:"createdAt"`
}

type TokenData struct {
	User      *User `json:"user"`
	ExpiresAt int64 `json:"exp"`
}

func (t TokenData) MapClaims() jwt.MapClaims {
	claims := make(jwt.MapClaims)

	tokenValue := reflect.ValueOf(t)
	tokenType := reflect.TypeOf(t)

	for i := 0; i < tokenValue.NumField(); i++ {
		field := tokenValue.Field(i)
		tag := tokenType.Field(i).Tag.Get("json")
		claims[tag] = field.Interface()
	}
	return claims
}

func (t *TokenData) FromClaims(claims map[string]interface{}) {
	jsonData, err := sonic.Marshal(claims)
	if err != nil {
		log.Fatal(err)
	}

	err = sonic.Unmarshal(jsonData, t)
	if err != nil {
		log.Fatal(err)
	}
}

func (t *TokenData) GenerateToken(secret string) (string, error) {
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, t.MapClaims()).SignedString([]byte(secret))

	if err != nil {
		return "", err
	}
	return token, nil
}

func (t *TokenData) ParseToken(token string, secret string) (bool, error) {
	tokenData, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if tokenData == nil {
		t = &TokenData{}
		return false, err
	}

	claims, ok := tokenData.Claims.(jwt.MapClaims)
	t.FromClaims(claims)
	return ok && tokenData.Valid, err
}
