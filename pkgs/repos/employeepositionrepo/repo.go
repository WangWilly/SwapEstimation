package employeepositionrepo

type repo struct{}

func New() *repo {
	return &repo{}
}
