package user

import (
	"errors"

	"github.com/MihaiBlebea/go-access-control/email"
)

type service struct {
	userRepo     *UserRepo
	emailService email.Service
}

type Service interface {
	Register(projectID int, firstName, lastName, email, password string) (*User, error)
	Login(email, password string) (*User, error)
	Authorize(projectID int, token string) (*User, error)
	RefreshToken(token string) (*User, error)
	RemoveUser(projectID int, token string) (int, error)
	ConfirmUser(confirmToken string) (*User, error)
}

func NewService(userRepo *UserRepo, emailService email.Service) Service {
	return &service{userRepo, emailService}
}

func (s *service) Register(
	projectID int,
	firstName, lastName, email, password string) (*User, error) {

	u, err := New(projectID, firstName, lastName, email, password)
	if err != nil {
		return &User{}, err
	}

	if err := s.userRepo.store(u); err != nil {
		return &User{}, err
	}

	if err := s.updateAccessToken(u); err != nil {
		return &User{}, err
	}

	if err := s.userRepo.update(u); err != nil {
		return &User{}, err
	}

	if err := s.emailService.ConfirmEmail(u.Email, u.ConfirmToken); err != nil {
		return &User{}, err
	}

	return u, nil
}

func (s *service) Login(email, password string) (*User, error) {
	u, err := s.userRepo.userWithEmail(email)
	if err != nil {
		return &User{}, err
	}

	if !u.Confirmed {
		return &User{}, errors.New("user is not confirmed")
	}

	if ok := u.validatePasswordHash(password); !ok {
		return &User{}, errors.New("password or email is invalid")
	}

	if err := s.updateAccessToken(u); err != nil {
		return &User{}, err
	}

	if err := s.updateRefreshToken(u); err != nil {
		return &User{}, err
	}

	if err := s.userRepo.update(u); err != nil {
		return &User{}, err
	}

	return u, nil
}

func (s *service) Authorize(projectID int, token string) (*User, error) {
	u, err := s.userRepo.userWithAccessToken(token)
	if err != nil {
		return &User{}, err
	}

	if !u.Confirmed {
		return &User{}, errors.New("user is not confirmed")
	}

	ok, err := u.validateAccessToken(token)
	if err != nil {
		return &User{}, err
	}

	if !ok {
		return &User{}, errors.New("user access_token is invalid")
	}

	return u, nil
}

func (s *service) RefreshToken(token string) (*User, error) {
	u, err := s.userRepo.userWithRefreshToken(token)
	if err != nil {
		return &User{}, err
	}

	if !u.Confirmed {
		return &User{}, errors.New("user is not confirmed")
	}

	ok, err := u.validateRefreshToken(token)
	if err != nil {
		return &User{}, err
	}

	if !ok {
		return &User{}, errors.New("user refresh token is invalid")
	}

	if err := s.updateAccessToken(u); err != nil {
		return &User{}, err
	}

	if err := s.userRepo.update(u); err != nil {
		return &User{}, err
	}

	return u, nil
}

func (s *service) RemoveUser(projectID int, token string) (int, error) {
	u, err := s.Authorize(projectID, token)
	if err != nil {
		return 0, err
	}

	if err := s.userRepo.delete(u); err != nil {
		return 0, err
	}

	return u.ID, nil
}

func (s *service) ConfirmUser(confirmToken string) (*User, error) {
	u, err := s.userRepo.userWithConfirmToken(confirmToken)
	if err != nil {
		return &User{}, err
	}

	u.confirm()

	if err := s.userRepo.update(u); err != nil {
		return &User{}, err
	}

	return u, nil
}

func (s *service) updateAccessToken(u *User) error {
	token, err := u.generateAccessToken()
	if err != nil {
		return err
	}
	if token == "" {
		return errors.New("token cannot be an empty string")
	}

	u.AccessToken = token

	return nil
}

func (s *service) updateRefreshToken(u *User) error {
	token, err := u.generateRefreshToken()
	if err != nil {
		return err
	}
	if token == "" {
		return errors.New("token cannot be an empty string")
	}

	u.RefreshToken = token

	return nil
}
