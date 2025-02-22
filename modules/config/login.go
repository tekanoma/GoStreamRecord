package config

type Logins struct {
	Users []Login
}

type Login struct {
	Name string `json:"name"`
	Key  string `json:"key"`
}
