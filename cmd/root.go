package cmd

import (
	cli "gopkg.in/urfave/cli.v2"
)

// NewApp new root command
func NewApp() *cli.App {
	return &cli.App{
		Description: "goapp serve",
		Version:     "0.1.0",
		Commands: []*cli.Command{
			NewServeCmd(),
			NewInitCmd(),
		},
		Before: func(ctx *cli.Context) error {
			return nil
		},

		Action: func(c *cli.Context) error {
			return nil
		}}
}
