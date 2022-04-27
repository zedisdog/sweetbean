package tools

import (
	"testing"

	"github.com/golang-jwt/jwt/v4"
	"github.com/uniplaces/carbon"
)

func TestExpire(t *testing.T) {
	key := []byte("123")
	claims := jwt.RegisteredClaims{
		ExpiresAt: &jwt.NumericDate{Time: carbon.Now().AddMinute().Time},
	}
	token, err := GenerateToken(key, claims)
	if err != nil {
		t.Fatal(err)
	}
	_, err = Parse(token, key)
	if err != nil {
		t.Fatal(err)
	}

	claims = jwt.RegisteredClaims{
		ExpiresAt: &jwt.NumericDate{Time: carbon.Now().SubSecond().Time},
	}
	token, err = GenerateToken(key, claims)
	if err != nil {
		t.Fatal(err)
	}
	_, err = Parse(token, key)
	if err == nil {
		t.Fatal("token is still active")
	}
}
