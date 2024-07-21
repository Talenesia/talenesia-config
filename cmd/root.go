package cmd

import (
	"log"
	"os"

	"github.com/spf13/cobra"
)

type Root struct {
	Cmd *cobra.Command
}

func New() *Root {
	r := &Root{
		Cmd: &cobra.Command{
			Use:     "talenesia",
			Short:   "Internal Command Line for Deployment Management",
			Example: "talenesia env",
			RunE: func(cmd *cobra.Command, args []string) error {
				return cmd.Help()
			},
		},
	}

	r.Cmd.AddCommand(r.EnvCmd())
	r.Cmd.AddCommand(r.ReleaseCmd())
	return r
}

func (r *Root) Execute() {
	if err := r.Cmd.Execute(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
