package user

import "errors"

type service struct {
	userRepo *UserRepo
}

type Service interface {
	Register(firstName, lastName, email, password string) (*User, error)
	Login(email, password string) (*User, error)
	Authorize(token string) (*User, error)
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

	if err := s.createOrUpdateUserToken(u); err != nil {
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

	if err := s.createOrUpdateUserToken(u); err != nil {
		return &User{}, err
	}

	return u, nil
}

func (s *service) Authorize(token string) (*User, error) {
	u, err := s.userRepo.userWithToken(token)
	if err != nil {
		return &User{}, err
	}

	_, err = u.ValidateToken(token)
	if err != nil {
		return &User{}, err
	}

	return u, nil
}

func (s *service) createOrUpdateUserToken(u *User) error {
	token, err := u.GenerateToken()
	if err != nil {
		return err
	}
	if token == "" {
		return errors.New("token cannot be an empty string")
	}

	u.Token = token

	if err := s.userRepo.update(u); err != nil {
		return err
	}

	return nil
}
