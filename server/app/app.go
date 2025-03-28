package app

type App struct {
	repo  Repository
	cache Cache
}

func New(repo Repository, cache Cache) *App {
	return &App{
		repo:  repo,
		cache: cache,
	}
}
