package serve

import (
	"github.com/spf13/cobra"
	"golang-monorepo-boilerplate/core"
	"golang-monorepo-boilerplate/core/config"
)

func newServeCmd() (cmd *cobra.Command) {
	cmd = &cobra.Command{
		Use:   "serve",
		Short: "Start the api server and ui server",
		Run: func(cmd *cobra.Command, args []string) {
			d, _ := core.New(cmd.Context(), cmd)
			err := d.Init(cmd.Context())
			if err != nil {
				d.Logger().WithError(err).Errorf("Unable to initialize app core.")
				return
			}
			d.Serve(cmd.Context())
		},
	}

	return cmd
}

func RegisterCommandRecursive(parent *cobra.Command) {
	c := newServeCmd()
	config.RegisterServeFlags(c.PersistentFlags())
	parent.AddCommand(c)
}
