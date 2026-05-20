package config

import "os"

type Config struct {
    Name      string
    Contact   string
    Certs     string
    Education []string
	Port      string
}

func Load() *Config {
	return &Config{
        Name:    os.Getenv("NAME"),
        Contact: os.Getenv("CONTACT"),
        Certs:   os.Getenv("CERTS"),
        Education: []string{
            os.Getenv("EDUCATION_1"),
            os.Getenv("EDUCATION_2"),
        },
		Port: os.Getenv("PORT"),
    }
}
