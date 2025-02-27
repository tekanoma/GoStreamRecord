package file

type API_secrets struct {
	Keys []secrets `json:keys`
}

type secrets struct {
	User string `json:user`
	Name string `json:name`
	Key  string `json:secret`
}

func (a API_secrets) NewKey() secrets {
	return secrets{}
}
