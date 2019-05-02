package main

import (
	"fmt"
	"log"
	"os"

	"github.com/koolay/goapp/internal/http"

	"github.com/koolay/goapp/cmd"
)

var (
	version = "0.0.0"
	build   = "---"
)

func main() {
	v := fmt.Sprintf("version=%s date=%s", version, build)
	http.Version = v
	rootCmd := cmd.NewApp()
	rootCmd.Version = v
	log.SetOutput(os.Stdout)

	if err := rootCmd.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
