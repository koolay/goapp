package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/koolay/goapp/conf"
	cli "gopkg.in/urfave/cli.v2"
)

func NewInitCmd() *cli.Command {
	initCmd := &cli.Command{
		Name:  "init",
		Usage: "init config file",
		Action: func(c *cli.Context) error {
			defaultConfigFile := filepath.Join(".", defaultConfilgFile)
			configFile, err := conf.InitDeafultCfgFile(defaultConfigFile)
			if err != nil {
				return err
			}
			fmt.Println("Generated config file", configFile)
			return err
		},
	}
	return initCmd
}
