package main

import (
	"log"
	"log/slog"

	app "hotel/app"
	"hotel/internal/config"
	"hotel/internal/db"
	"hotel/internal/db/seed"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.Load()
		if err != nil {
			log.Fatalf("config: %v", err)
		}

		application, err := app.New(cfg)
		if err != nil {
			log.Fatalf("initialize app: %v", err)
		}

		if err := application.Run(); err != nil {
			log.Fatalf("run app: %v", err)
		}
	},
}

var seedCmd = &cobra.Command{
	Use:   "seed",
	Short: "Seed the database",
	Run: func(cmd *cobra.Command, args []string) {
		slog.Info("Seeding started")
		cfg, err := config.Load()
		if err != nil {
			log.Fatalf("config: %v", err)
		}

		database, err := db.Open(cfg.DBPath)
		if err != nil {
			panic("Couldn't connect to the database")
		}
		seed.Seed(database, cfg)
		slog.Info("Seeding Complete")
	},
}

func main() {
	rootCmd.AddCommand(seedCmd)
	rootCmd.Execute()
}
