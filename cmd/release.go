package cmd

import (
	"io"
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

	for cmdLabel, usingPass := range r.ReleaseCommandList() {
		slog.Info("start running command...", "cmd", cmdLabel)
		if !usingPass {
			cmd := exec.Command("/bin/sh", "-c", cmdLabel)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr

			err := cmd.Run()
			if err != nil {
				slog.Error("error executing the command", "cmd", cmd, "error", err)
				return
			}

			continue
		}

		cmd := exec.Command("/bin/sh", "-c", cmdLabel)

		stdin, err := cmd.StdinPipe()
		if err != nil {
			slog.Error("error executing the command", "cmd", cmd, "error", err)
			return
		}

		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if passphrase != "" && usingPass {
			go func() {
				stdin.Close()
				io.WriteString(stdin, passphrase+"\n")
			}()
		}

		if output, err := cmd.CombinedOutput(); err != nil {
			slog.Error("error executing the command", "cmd", cmdLabel, "error", err)
			return
		} else {
			slog.Info("result", "output", output)
		}

	}
}

func (r *Root) ReleaseCommandList() map[string]bool {
	return map[string]bool{
		"set -e":                            false,
		`echo "🛫 Deploying application"`:    false,
		"(sudo php artisan down) || true":   false,
		"sudo git fetch origin main":        true,
		`sudo git reset --hard origin/main`: false,
		`sudo composer install --no-interaction --prefer-dist --optimize-autoloader --no-dev`: false,
		`sudo php artisan migrate --force`:                   false,
		`sudo php artisan optimize`:                          false,
		`sudo php artisan config:clear`:                      false,
		`sudo php artisan cache:clear`:                       false,
		`sudo find . -type f -exec chmod 644 {} \;`:          false,
		`sudo find . -type d -exec chmod 755 {} \;`:          false,
		`sudo chown -R www-data:www-data .`:                  false,
		`sudo chgrp -R www-data ./storage ./bootstrap/cache`: false,
		`sudo chmod -R ug+rwx ./storage ./bootstrap/cache`:   false,
		`sudo php artisan up`:                                false,
		`sudo systemctl restart laravel_worker.service`:      false,
		`echo "🚀 Deployed application"`:                      false,
	}
}
