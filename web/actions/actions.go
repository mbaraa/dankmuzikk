package actions

type Actions struct {
	cache Cache
}

func New(cache Cache) *Actions {
	return &Actions{
		cache: cache,
	}
}
