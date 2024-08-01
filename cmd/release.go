package cmd

import (
	"log/slog"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
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
	stage, err := cmd.Flags().GetString("stage")
	if err != nil {
		slog.Error("error reading stage", "err", err)
		return
	}

	var passphrase string
	if stage == "production" {
		passphrase, err = cmd.Flags().GetString("passphrase")
		if err != nil {
			slog.Error("error reading passphrase", "err", err)
			return
		}
	}

	for _, cmd := range r.ReleaseCommandList() {
		slog.Info("start running command...", "cmd", cmd)

		cmd := exec.Command("/bin/sh", "-c", cmd)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if passphrase != "" {
			cmd.Stdin = strings.NewReader(passphrase)
		}
		err := cmd.Run()
		if err != nil {
			slog.Error("error executing the command", "cmd", cmd, "error", err)
			return
		}
	}
}

func (r *Root) ReleaseCommandList() []string {
	return []string{
		"set -e",
		`echo "🛫 Deploying application"`,
		"(sudo php artisan down) || true",
		"sudo -S git fetch origin main",
		`sudo -S git reset --hard origin/main`,
		`sudo composer install --no-interaction --prefer-dist --optimize-autoloader --no-dev`,
		`sudo php artisan migrate --force`,
		`sudo php artisan optimize`,
		`sudo php artisan config:clear`,
		`sudo php artisan cache:clear`,
		`sudo find . -type f -exec chmod 644 {} \;`,
		`sudo find . -type d -exec chmod 755 {} \;`,
		`sudo chown -R www-data:www-data .`,
		`sudo chgrp -R www-data ./storage ./bootstrap/cache`,
		`sudo chmod -R ug+rwx ./storage ./bootstrap/cache`,
		`sudo php artisan up`,
		`sudo systemctl restart laravel_worker.service`,
		`echo "🚀 Deployed application"`,
	}
}
