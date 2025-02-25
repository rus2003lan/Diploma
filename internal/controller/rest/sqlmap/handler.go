package sqlmap

type Handler struct {
	s Service
}

func NewHandler(
	rs Service,
) *Handler {
	return &Handler{
		s: rs,
	}
}
