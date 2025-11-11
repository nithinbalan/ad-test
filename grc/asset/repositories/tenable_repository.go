package repositories

type TenableRepositoryInterface interface {
}

type TenableRepository struct{}

func NewTenableRepository() TenableRepositoryInterface {
	return &TenableRepository{}
}
