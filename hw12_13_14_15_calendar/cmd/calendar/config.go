package main

import (
	"log"
	"os"

	toml "github.com/naoina/toml"
)

// При желании конфигурацию можно вынести в internal/config.
// Организация конфига в main принуждает нас сужать API компонентов, использовать
// при их конструировании только необходимые параметры, а также уменьшает вероятность циклической зависимости.
type Config struct {
	Logger   LoggerConf
	Database DBConf
}

type DBConf struct {
	DBimplement string
	FilePath    string
	Server      string
	Port        int
	User        string
	Password    string
}

type LoggerConf struct {
	Level string
}

func CreateNewConfig(filename string) *Config {
	file, err := os.Open(filename)
	if err != nil {
		log.Println(err)
		return nil
	}
	defer file.Close()
	config := &Config{}
	if err := toml.NewDecoder(file).Decode(config); err != nil {
		log.Println(err)
		return nil
	}
	return config
}
