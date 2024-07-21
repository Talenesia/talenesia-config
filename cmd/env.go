package cmd

import (
	"log/slog"
	"strings"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
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
		if val, exist := envMap[key]; exist {
			slog.Info("Update Env", "Before", val, "After", env)
			envMap[key] = env
		}
	}

	err = godotenv.Write(envMap, envDest)
	if err != nil {
		slog.Error("error write env file", "err", err)
		return
	}
}
