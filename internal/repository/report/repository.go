package report

type Repository struct {
	es Elastic
}

func New(es Elastic) *Repository {
	return &Repository{
		es: es,
	}
}
