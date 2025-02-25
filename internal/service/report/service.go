package report

type Service struct {
	r ReportRepo
	s SQLMapService
}

func New(r ReportRepo, s SQLMapService,
) *Service {
	return &Service{
		r: r,
		s: s,
	}
}
