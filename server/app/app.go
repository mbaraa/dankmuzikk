package app

type App struct {
	repo Repository
}

func New(repository Repository) *App {
	return &App{
		repo: repository,
	}
}
