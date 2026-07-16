package repositories

type QualysRepositoryInterface interface {
}

type QualysRepository struct{}

func NewQualysRepository() QualysRepositoryInterface {
	return &QualysRepository{}
}
