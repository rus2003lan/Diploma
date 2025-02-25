package sqlmap

type Service struct {
	r Repo
}

func New(r Repo,
) *Service {
	return &Service{
		r: r,
	}
}
