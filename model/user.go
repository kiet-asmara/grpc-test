package model

type UserModel struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type UserAll struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

type UserNamePass struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}
