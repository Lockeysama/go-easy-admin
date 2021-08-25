package config

type Config struct {
	Port int16 `json:"port" default:"8080"`
	Name string
	A    bool    `default:"true"`
	B    float32 `default:"8080.123"`
	C    string  `default:"8080xxx"`
}
