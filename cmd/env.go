package cmd

import (
	"log/slog"
	"os"
	"os/exec"
	"strings"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
	"github.com/talenesia/router/config"
)

func (r *Root) EnvCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "env",
		Short: "Apply for new or updated environment variables",
		Long:  "Apply for new or updated environment variables",
		Run:   r.ApplyEnv,
	}
	cmd.PersistentFlags().String("envars", "", `Environment Variable List (seperated by comma). APP_NAME="Talenesia Web",APP_DEBUG=prod`)
	cmd.PersistentFlags().String("src", "", "Environment Source Paths (.env)")
	cmd.PersistentFlags().String("dest", "", "Environment Destination Paths (.env)")

	return cmd
}

func (r *Root) ApplyEnv(cmd *cobra.Command, args []string) {
	envsFlag, err := cmd.Flags().GetString("envars")
	if err != nil {
		slog.Error("error read envs from flag", "err", err)
		return
	}

	envPath, err := cmd.Flags().GetString("src")
	if err != nil {
		slog.Error("error read src from flag", "err", err)
		return
	}

	envDest, err := cmd.Flags().GetString("dest")
	if err != nil {
		slog.Error("error read dest from flag", "err", err)
		return
	}

	envs := strings.Split(envsFlag, ",")

	filteredEnvs := make(map[string]string)
	for _, env := range envs {
		splitEnv := strings.Split(env, "=")
		if len(splitEnv) > 1 {
			filteredEnvs[splitEnv[0]] = splitEnv[1]
		}
	}

	envMap, err := godotenv.Read(envPath)
	if err != nil {
		slog.Error("error read envpath file", "err", err)
		return
	}

	for key, env := range filteredEnvs {
		envMap[key] = env
	}

	err = godotenv.Write(envMap, envDest)
	if err != nil {
		slog.Error("error write env file", "err", err)
		return
	}

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
