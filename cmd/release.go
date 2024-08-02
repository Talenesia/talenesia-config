package cmd

import (
	"bufio"
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/exec"
	"time"

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

	for _, cmdLabel := range r.ReleaseCommandList() {
		slog.Info("start running command...", "cmd", cmdLabel)
		if cmdLabel != "sudo git fetch origin main" {
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

		stdout, err := cmd.StdoutPipe()
		if err != nil {
			slog.Error("error getting stdout pipe", "err", err)
			return
		}

		stderr, err := cmd.StderrPipe()
		if err != nil {
			slog.Error("error getting stderr pipe", "err", err)
			return
		}

		err = cmd.Start()
		if err != nil {
			slog.Error("error starting command", "err", err)
			return
		}

		// Function to check for passphrase prompt and respond
		go func() {
			defer stdin.Close()
			scanner := bufio.NewScanner(io.MultiReader(stdout, stderr))
			for scanner.Scan() {
				text := scanner.Text()
				fmt.Println(text) // Print command output
				slog.Info("Passphrase prompt detected, sending passphrase")

				time.Sleep(500 * time.Millisecond)
				_, err := io.WriteString(stdin, passphrase+"\n")
				if err != nil {
					slog.Error("Failed to write passphrase", "error", err)
				} else {
					slog.Info("Passphrase sent successfully")
				}
			}
		}()

		// Wait for the command to finish
		if err := cmd.Wait(); err != nil {
			slog.Error("command finished with error", "err", err)
			return
		}
	}
}

func (r *Root) ReleaseCommandList() []string {
	return []string{
		"set -e",
		`echo "🛫 Deploying application"`,
		"(sudo php artisan down) || true",
		"sudo git fetch origin main",
		`sudo git reset --hard origin/main`,
		`sudo composer install --no-interaction --prefer-dist --optimize-autoloader --no-dev`,
		`sudo php artisan migrate --force`,
		`sudo php artisan optimize`,
		`sudo php artisan config,clear`,
		`sudo php artisan cache,clear`,
		`sudo find . -type f -exec chmod 644 {} \;`,
		`sudo find . -type d -exec chmod 755 {} \;`,
		`sudo chown -R www-data,www-data .`,
		`sudo chgrp -R www-data ./storage ./bootstrap/cache`,
		`sudo chmod -R ug+rwx ./storage ./bootstrap/cache`,
		`sudo php artisan up`,
		`sudo systemctl restart laravel_worker.service`,
		`echo "🚀 Deployed application"`,
	}
}
