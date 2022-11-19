package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"

	"github.com/burnb/duocard/internal/configs"
	"github.com/burnb/duocard/internal/importer"
)

func main() {
	if err := godotenv.Load(); err != nil {
		if _, ok := err.(*os.PathError); !ok {
			log.Fatal(err)
		}
	}

	cfg := &configs.App{}
	if err := cfg.Prepare(); err != nil {
		log.Fatal(err)
	}

	switch cfg.Command {
	case "import":
		importerSrv := importer.New(cfg)
		if err := importerSrv.Run(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	default:
		fmt.Printf("%q: no such command", cfg.Command)
	}
}
