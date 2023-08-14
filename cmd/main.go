package main

import (
	"Vectory/api/handlers"
	"Vectory/db"
	"Vectory/gen/api/restapi"
	"Vectory/gen/api/restapi/operations"
	"flag"
	"github.com/go-openapi/loads"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
)

type Config struct {
	FilesPath  string `yaml:"files_path"`
	ListenPort int    `yaml:"listen_port"`
}

// main is invoked when deploying Vectory on the cloud.
func main() {
	var cfgPath string

	flag.StringVar(&cfgPath, "config", "", "config path for Vectory")
	flag.Parse()

	cfg, err := readConfig(cfgPath)
	if err != nil {
		log.Fatalf("startup: %v", err)
	}

	vectoryDB, err := db.Open(cfg.FilesPath)
	if err != nil {
		log.Fatalf("startup: %v", err)
	}

	apiSpec, err := loads.Spec("./api/spec.yaml")
	if err != nil {
		log.Fatalf("startup: %v", err)
	}

	api := operations.NewVectoryAPI(apiSpec)
	handlers.InitHandlers(api, vectoryDB)

	server := restapi.NewServer(api)
	server.Port = cfg.ListenPort

	defer server.Shutdown()

	server.ConfigureAPI()

	err = server.Serve()
	if err != nil {
		log.Fatal(err)
	}
}

func readConfig(path string) (*Config, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config

	err = yaml.Unmarshal(b, &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
