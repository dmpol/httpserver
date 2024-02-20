package handlers

type PersonAuth struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Person struct {
	PersonAuth
	Email string `json:"email"`
}

type PersonWithId struct {
	Id int `json:"id"`
	Person
}
