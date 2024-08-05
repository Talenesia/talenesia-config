package cmd

import (
	"log/slog"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
	"github.com/talenesia/router/config"
)

func (r *Root) ReleaseCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "release",
		Short: "Apply for new latest release",
		Long:  "Apply for new latest release",
		Run:   r.Release,
	}

	cmd.PersistentFlags().String("passphrase", "", "Passphrase only for production")
	cmd.PersistentFlags().String("stage", "", "Stage production or staging")

	return cmd
}

func (r *Root) Release(cmd *cobra.Command, args []string) {
	labels, exist := config.Conf.Commands[cmd.Use]
	if !exist {
		slog.Error("command does not exist")
		return
	}

	for _, cmdLabel := range labels {
		slog.Info("start running command...", "cmd", cmdLabel)
		cmd := exec.Command("/bin/sh", "-c", cmdLabel)

		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		err := cmd.Run()
		if err != nil {
			slog.Error("error executing the command", "cmd", cmd, "error", err)
			return
		}
	}
}
