package ports

type UnitOfWork interface {
	Factory() RepositoryFactory
	WithTransaction(fn func(RepositoryFactory) error) error
}
