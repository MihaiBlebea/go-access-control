package proj

import (
	"math/rand"
	"strings"
	"time"
)

const apiKeyLength = 20

type Project struct {
	ID        int       `json:"id"`
	Name      string    `json:"name" gorm:"unique"`
	Slug      string    `json:"slug" gorm:"unique"`
	Host      string    `json:"host" gorm:"unique"`
	ApiKey    string    `json:"api_key"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func New(name, host string) *Project {
	return &Project{
		Name:   name,
		Slug:   toSlug(name),
		Host:   host,
		ApiKey: genApiKey(apiKeyLength),
	}
}

func (p *Project) RegenApiKey() {
	p.ApiKey = genApiKey(apiKeyLength)
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

func toSlug(name string) string {
	return strings.ToLower(strings.ReplaceAll(name, " ", "-"))
}
