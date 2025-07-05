package actions

type Actions struct {
	requests Requests
	cache    Cache
}

func New(requests Requests, cache Cache) *Actions {
	return &Actions{
		requests: requests,
		cache:    cache,
	}
}
