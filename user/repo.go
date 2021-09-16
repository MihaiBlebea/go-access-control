package user

import (
	"errors"

	"gorm.io/gorm"
)

var (
	ErrNoRecord  error = errors.New("record not found")
	ErrNoRecords error = errors.New("records not found with filter")
)

type UserRepo struct {
	conn *gorm.DB
}

func NewRepo(conn *gorm.DB) *UserRepo {
	return &UserRepo{conn}
}

func (r *UserRepo) userWithAccessToken(token string) (*User, error) {
	user := User{}
	err := r.conn.Where("access_token = ?", token).Find(&user).Error
	if err != nil {
		return &user, err
	}

	if user.ID == 0 {
		return &user, ErrNoRecord
	}

	return &user, nil
}

func (r *UserRepo) userWithRefreshToken(token string) (*User, error) {
	user := User{}
	err := r.conn.Where("refresh_token = ?", token).Find(&user).Error
	if err != nil {
		return &user, err
	}

	if user.ID == 0 {
		return &user, ErrNoRecord
	}

	return &user, nil
}

func (r *UserRepo) userWithEmail(email string) (*User, error) {
	user := User{}
	err := r.conn.Where("email = ?", email).Find(&user).Error
	if err != nil {
		return &user, err
	}

	if user.ID == 0 {
		return &user, ErrNoRecord
	}

	return &user, nil
}

func (r *UserRepo) userWithConfirmToken(confirmToken string) (*User, error) {
	user := User{}
	err := r.conn.Where("confirm_token = ?", confirmToken).Find(&user).Error
	if err != nil {
		return &user, err
	}

	if user.ID == 0 {
		return &user, ErrNoRecord
	}

	return &user, err
}

func (r *UserRepo) usersWithProjectID(projectID int) ([]User, error) {
	users := []User{}
	err := r.conn.Where("project_id = ?", projectID).Find(&users).Error
	if err != nil {
		return users, err
	}

	return users, nil
}

func (r *UserRepo) store(user *User) error {
	return r.conn.Create(user).Error
}

func (r *UserRepo) update(user *User) error {
	return r.conn.Save(user).Error
}

func (r *UserRepo) delete(user *User) error {
	return r.conn.Delete(user).Error
}
