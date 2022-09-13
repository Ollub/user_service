package users

type User struct {
	ID        uint32 `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Ver       int    `json:"-"`
	PassHash  string `json:"-"`
}

type UserIn struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

type UserUpdate struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}
