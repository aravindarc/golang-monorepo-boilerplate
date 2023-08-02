package migrate

import (
	"github.com/spf13/cobra"
	"golang-monorepo-boilerplate/core"
	"golang-monorepo-boilerplate/core/config"
)

func newMigrateCmd() (cmd *cobra.Command) {
	const cmdDesc = `Supported arguments are:
- up            - runs all available migrations
- down [number] - reverts the last [number] applied migrations
`

	cmd = &cobra.Command{
		Use:   "migrate",
		Short: "Migrate the database",
		Long:  cmdDesc,
		Run: func(cmd *cobra.Command, args []string) {
			// todo validate args
			d, _ := core.New(cmd.Context(), cmd)
			runner, err := d.Persister().MigrateRunner()
			if err != nil {
				panic(err)
			}
			autoConfirm, _ := cmd.Flags().GetBool("auto-confirm")
			err = runner.Run(autoConfirm, args...)
			if err != nil {
				panic(err)
			}
		},
	}

	cmd.PersistentFlags().Bool("auto-confirm", false, "confirm all the migration without prompt")

	return cmd
}

func RegisterCommandRecursive(parent *cobra.Command) {
	c := newMigrateCmd()
	config.RegisterMigrateFlags(c.PersistentFlags())
	parent.AddCommand(c)
}
