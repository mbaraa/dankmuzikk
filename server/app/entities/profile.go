package entities

type Profile struct {
	Id       uint `json:"-"`
	Email    string
	Name     string
	PfpLink  string
	Username string
}
