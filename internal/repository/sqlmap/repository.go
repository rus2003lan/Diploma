package sqlmap

type Repository struct {
	ceph   Ceph
	bucket string
}

func New(c Ceph, b string) *Repository {
	return &Repository{
		bucket: b,
		ceph:   c,
	}
}
