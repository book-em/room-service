package internal

type Service interface {
}

type service struct {
	repo Repository
}

func NewService(r Repository) Service {
	return &service{r}
}
