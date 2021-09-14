package user

import (
	"errors"
	"os"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        int       `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email" gorm:"unique"`
	Password  string    `json:"password"`
	Token     string    `json:"token"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type customClaims struct {
	ID    int    `json:"id"`
	Email string `json:"username"`
	jwt.StandardClaims
}

func New(firstName, lastName, email, password string) (*User, error) {
	hash, err := hashPassword(password)
	if err != nil {
		return &User{}, err
	}

	return &User{
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		Password:  hash,
	}, nil
}

func (u *User) GenerateToken() (string, error) {
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

func (u *User) ValidateToken(token string) (bool, error) {
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

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)

	return string(bytes), err
}

func (u *User) validatePasswordHash(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))

	return err == nil
}
