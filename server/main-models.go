package main

type Config struct {
	Database struct {
		Host     string `json:"host"`
		Port     int    `json:"port"`
		Username string `json:"username"`
		Password string `json:"password"`
		Schema   string `json:"schema"`
		Database string `json:"database"`
	} `json:"database"`
	User struct {
		Secret     string `json:"secret"`
		PrivateKey string `json:"private_key"`
		PublicKey  string `json:"public_key"`
	} `json:"user"`
}

type DatabaseConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	Schema   string `json:"schema"`
	Database string `json:"database"`
}
