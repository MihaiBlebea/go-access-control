package proj

import (
	"errors"

	"gorm.io/gorm"
)

var (
	ErrNoRecord  error = errors.New("record not found")
	ErrNoRecords error = errors.New("records not found with filter")
)

type ProjectRepo struct {
	conn *gorm.DB
}

func NewRepo(conn *gorm.DB) *ProjectRepo {
	return &ProjectRepo{conn}
}

func (r *ProjectRepo) projectWithApiKey(apiKey string) (*Project, error) {
	project := Project{}
	err := r.conn.Where("api_key = ?", apiKey).Find(&project).Error
	if err != nil {
		return &project, err
	}

	if project.ID == 0 {
		return &project, ErrNoRecord
	}

	return &project, nil
}

func (r *ProjectRepo) projectWithSlug(slug string) (*Project, error) {
	project := Project{}
	err := r.conn.Where("slug = ?", slug).Find(&project).Error
	if err != nil {
		return &project, err
	}

	if project.ID == 0 {
		return &project, ErrNoRecord
	}

	return &project, nil
}

func (r *ProjectRepo) store(user *Project) error {
	return r.conn.Create(user).Error
}

func (r *ProjectRepo) update(user *Project) error {
	return r.conn.Save(user).Error
}

func (r *ProjectRepo) delete(user *Project) error {
	return r.conn.Delete(user).Error
}
