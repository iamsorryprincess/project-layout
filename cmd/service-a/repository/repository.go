package repository

type Repository struct {
}

func New() *Repository {
	return &Repository{}
}

func (r *Repository) GetData() ([]string, error) {
	return []string{"test1", "test2", "test3"}, nil
}

func (r *Repository) SaveData(data []string) error {
	return nil
}
