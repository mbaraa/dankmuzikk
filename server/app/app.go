package app

type App struct {
	repo        Repository
	cache       Cache
	playerCache PlayerCache
}

func New(repo Repository, cache Cache, playerCache PlayerCache) *App {
	return &App{
		repo:        repo,
		cache:       cache,
		playerCache: playerCache,
	}
}
