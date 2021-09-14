package user

import "errors"

type service struct {
	userRepo *UserRepo
}

type Service interface {
	Register(firstName, lastName, email, password string) (*User, error)
	Login(email, password string) (*User, error)
	Authorize(token string) (*User, error)
	RefreshToken(token string) (*User, error)
}

func NewService(userRepo *UserRepo) Service {
	return &service{userRepo}
}

func (s *service) Register(firstName, lastName, email, password string) (*User, error) {
	u, err := New(firstName, lastName, email, password)
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

	return u, nil
}

func (s *service) Login(email, password string) (*User, error) {
	u, err := s.userRepo.userWithEmail(email)
	if err != nil {
		return &User{}, err
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

func (s *service) Authorize(token string) (*User, error) {
	u, err := s.userRepo.userWithAccessToken(token)
	if err != nil {
		return &User{}, err
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
