package cmd

import (
	"log/slog"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

func (r *Root) ReleaseCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "release",
		Short: "Apply for new latest release",
		Long:  "Apply for new latest release",
		Run:   r.Release,
	}

	return cmd
}

func (r *Root) Release(cmd *cobra.Command, args []string) {
	for _, cmd := range r.ReleaseCommandList() {
		cmd := exec.Command("/bin/sh", "-c", cmd)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

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
		"git fetch origin main",
		`git reset --hard origin/main`,
		`sudo composer install --no-interaction --prefer-dist --optimize-autoloader --no-dev`,
		`sudo php artisan migrate --force`,
		`sudo php artisan optimize`,
		`sudo find . -type f -exec chmod 644 {} \;`,
		`sudo find . -type d -exec chmod 755 {} \;`,
		`sudo chown -R www-data:www-data .`,
		`sudo chgrp -R www-data ./storage ./bootstrap/cache`,
		`sudo chmod -R ug+rwx ./storage ./bootstrap/cache`,
		`sudo php artisan up`,
		`echo "🚀 Deployed application"`,
	}
}
