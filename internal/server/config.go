package server

type config struct {
	ClientPort int
	Version    string
	ConnLimit  int
}

func GetConfig() *config {
	return &config{
		ClientPort: 1337,
		Version:    "0.1",
		ConnLimit:  2,
	}
}
