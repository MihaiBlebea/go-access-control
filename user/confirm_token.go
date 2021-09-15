package user

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

type confirmTokenClaims struct {
	ConfirmSuccessURL string `json:"confirm_success_url"`
	ConfirmFailURL    string `json:"confirm_fail_url"`
	ConfirmWebhook    string `json:"confirm_webhook"`
	jwt.StandardClaims
}

type Payload struct {
	Success bool `json:"success"`
	ID      int  `json:"user_id"`
}

func SendWebhook(url string, payload *Payload) error {
	if url == "" {
		return errors.New("invalid url")
	}

	b, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(
		"POST",
		url,
		bytes.NewBuffer(b),
	)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return errors.New("request failed")
	}

	return nil
}

func generateConfirmToken(successURL, failURL, webhook string) (string, error) {
	expireAt := time.Now().Add(time.Minute * 30)

	claims := confirmTokenClaims{
		ConfirmSuccessURL: successURL,
		ConfirmFailURL:    failURL,
		ConfirmWebhook:    webhook,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireAt.Local().Unix(),
			Issuer:    "google.com",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return "", errors.New("secret string cannot be empty")
	}

	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}

func parseConfirmToken(token string) (*confirmTokenClaims, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return &confirmTokenClaims{}, errors.New("secret string cannot be empty")
	}

	t, err := jwt.ParseWithClaims(
		token,
		&confirmTokenClaims{},
		func(tkn *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		},
	)
	if err != nil {
		return &confirmTokenClaims{}, err
	}

	claims, ok := t.Claims.(*confirmTokenClaims)
	if !ok {
		return &confirmTokenClaims{}, errors.New("couldn't parse claims")
	}

	if claims.ExpiresAt < time.Now().UTC().Unix() {
		return &confirmTokenClaims{}, errors.New("jwt is expired")
	}

	return claims, nil
}
