package main

import (
	"encoding/json"
	"github.com/FDUTCH/obsidian/proxy"
	"log"
	"os"
)

/*
	Universal solution  for basically all your proxy needs
*/

func main() {
	configs, err := loadConfig()
	if err != nil {
		configs, err = exampleConfigs()
	}

	if err != nil {
		log.Fatal(err)
	}

	for _, conf := range configs {
		go log.Fatal(conf.Run())
	}

	select {}

}

func exampleConfigs() ([]proxy.Options, error) {
	var configs = []proxy.Options{
		{
			Port:    80,
			Network: "http",
			Servers: []string{"www.wikipedia.org", "bedrock.dev", "pornhub.com", "example.com"},
		},
		{
			Port:          19132,
			Network:       "udp",
			RemoteAddress: "mco.cubecraft.net",
		},
		{
			Port:          25565,
			Network:       "tcp",
			RemoteAddress: "play.hypixel.net",
		},
	}
	file, err := os.OpenFile(configName, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	enc := json.NewEncoder(file)

	enc.SetIndent("", "    ")

	err = enc.Encode(configs)

	return configs, err
}

func loadConfig() ([]proxy.Options, error) {
	var configs []proxy.Options

	file, err := os.Open(configName)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	err = json.NewDecoder(file).Decode(&configs)

	return configs, err
}

const configName = "config.json"
