package proj

import (
	"math/rand"
	"time"
)

type Project struct {
	ID        int       `json:"id"`
	Name      string    `json:"name" gorm:"unique"`
	Host      string    `json:"host" gorm:"unique"`
	ApiKey    string    `json:"api_key"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func New(name, host string) *Project {
	return &Project{
		Name:   name,
		Host:   host,
		ApiKey: genApiKey(10),
	}
}

func (p *Project) RegenApiKey() {
	p.ApiKey = genApiKey(10)
}

func genApiKey(n int) string {
	letterBytes := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	rand.Seed(time.Now().UnixNano())

	b := make([]byte, n)
	for i, _ := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}

	return string(b)
}
