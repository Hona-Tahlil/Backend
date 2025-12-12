package domainpostgres

type QueryModifier interface {
	Apply(query interface{}) interface{}
}
