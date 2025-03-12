package entities

type Profile struct {
	Id       uint   `json:"-"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	PfpLink  string `json:"pfp_link"`
	Username string `json:"username"`
}
