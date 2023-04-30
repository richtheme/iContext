package main

import (
	"flag"
	"fmt"

	"iContext/internal/api"

	"github.com/BurntSushi/toml"
	log "github.com/sirupsen/logrus"
)

func main() {
	configPath := "configs/api.toml"
	config := api.NewConfig()

	_, err := toml.DecodeFile(configPath, config)
	if err != nil {
		log.Println("Can not find config file. using default values:", err)
	}

	redisHost := flag.String("host", "", "Redis host. (Required)")
	redisPort := flag.String("port", "", "Redis port. (Required)")
	flag.Parse()
	redisAddr := fmt.Sprintf("%s:%s", *redisHost, *redisPort)
	fmt.Println(redisAddr)

	server := api.New(config, redisAddr)
	fmt.Println(server)
	log.Fatal(server.Start())
}
