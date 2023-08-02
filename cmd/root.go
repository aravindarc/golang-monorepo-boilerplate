package cmd

import (
	"github.com/spf13/cobra"
	"golang-monorepo-boilerplate/cmd/migrate"
	"golang-monorepo-boilerplate/cmd/serve"
)

func NewRootCmd() (cmd *cobra.Command) {
	cmd = &cobra.Command{
		Use: "gmb",
	}

	return cmd
}

func Execute() {
	c := NewRootCmd()

	serve.RegisterCommandRecursive(c)
	migrate.RegisterCommandRecursive(c)

	err := c.Execute()
	if err != nil {
		panic(err)
	}
}
