package actions

type Actions struct {
	requests Requests
}

func New(requests Requests) *Actions {
	return &Actions{
		requests: requests,
	}
}

type RequestContext struct {
	SessionToken string
}
