package user

import (
	"errors"
	"os"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID           int       `json:"id"`
	FirstName    string    `json:"first_name"`
	LastName     string    `json:"last_name"`
	Email        string    `json:"email" gorm:"unique"`
	Password     string    `json:"password"`
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type customClaims struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
	jwt.StandardClaims
}

func New(firstName, lastName, email, password string) (*User, error) {
	hash, err := hashPassword(password)
	if err != nil {
		return &User{}, err
	}
	u := User{
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		Password:  hash,
	}

	refreshToken, err := u.generateRefreshToken()
	if err != nil {
		return &User{}, err
	}

	u.RefreshToken = refreshToken

	return &u, nil
}

func (u *User) generateAccessToken() (string, error) {
	if u.ID == 0 || u.Email == "" {
		return "", errors.New("user id or user email is not valid")
	}

	expireAt := time.Now().Add(time.Minute * 3)

	claims := customClaims{
		ID:    u.ID,
		Email: u.Email,
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

func (u *User) validateAccessToken(token string) (bool, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return false, errors.New("secret string cannot be empty")
	}

	t, err := jwt.ParseWithClaims(
		token,
		&customClaims{},
		func(tkn *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		},
	)
	if err != nil {
		return false, err
	}

	claims, ok := t.Claims.(*customClaims)
	if !ok {
		return false, errors.New("couldn't parse claims")
	}

	if claims.ExpiresAt < time.Now().UTC().Unix() {
		return false, errors.New("jwt is expired")
	}

	if claims.ID != u.ID {
		return false, errors.New("invalid user id")
	}

	return true, nil
}

func (u *User) generateRefreshToken() (string, error) {
	expireAt := time.Now().Add(time.Hour * 24 * 28 * 3)

	claims := jwt.StandardClaims{
		ExpiresAt: expireAt.Local().Unix(),
		Issuer:    "google.com",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return "", errors.New("secret string cannot be empty")
	}

	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}

func (u *User) validateRefreshToken(token string) (bool, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return false, errors.New("secret string cannot be empty")
	}

	t, err := jwt.ParseWithClaims(
		token,
		&jwt.StandardClaims{},
		func(tkn *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		},
	)
	if err != nil {
		return false, err
	}

	claims, ok := t.Claims.(*jwt.StandardClaims)
	if !ok {
		return false, errors.New("couldn't parse claims")
	}

	if claims.ExpiresAt < time.Now().UTC().Unix() {
		return false, errors.New("jwt is expired")
	}

	return true, nil
}

func (u *User) validatePasswordHash(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))

	return err == nil
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)

	return string(bytes), err
}
