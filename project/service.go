package proj

type Service interface {
	CreateProject(name, host string) (*Project, error)
	RegenApiKey(apiKey string) (*Project, error)
	GetProject(apiKey string) (*Project, error)
}

type service struct {
	projecRepo *ProjectRepo
}

func NewService(projecRepo *ProjectRepo) Service {
	return &service{projecRepo}
}

func (s *service) CreateProject(name, host string) (*Project, error) {

	p := New(name, host)

	if err := s.projecRepo.store(p); err != nil {
		return &Project{}, err
	}

	return p, nil
}

func (s *service) RegenApiKey(apiKey string) (*Project, error) {
	p, err := s.projecRepo.projectWithApiKey(apiKey)
	if err != nil {
		return &Project{}, err
	}

	p.RegenApiKey()

	if err := s.projecRepo.update(p); err != nil {
		return &Project{}, err
	}

	return p, nil
}

func (s *service) GetProject(apiKey string) (*Project, error) {
	p, err := s.projecRepo.projectWithApiKey(apiKey)
	if err != nil {
		return &Project{}, err
	}

	return p, nil
}
