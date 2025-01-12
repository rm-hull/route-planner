package main

import (
	"fmt"
	"log"
	"os"

	"github.com/earthboundkid/versioninfo/v2"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"

	"github.com/rm-hull/route-planner/cmds"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	var rootCmd = &cobra.Command{
		Use:  "route-planner",
		Long: `HTTP server, DB migration and data import/export`,
	}

	var importGmlCmd = &cobra.Command{
		Use:   "gml [path]",
		Short: "Import GML data from specified path",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if err := cmds.ImportGmlData(args[0]); err != nil {
				log.Fatalf("failed to import GML: %v", err)
			}
		},
	}

	var importRefDataCmd = &cobra.Command{
		Use:   "refdata [table-name] [url]",
		Short: "Import reference data",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			if err := cmds.ImportRefData(args[0], args[1]); err != nil {
				log.Fatalf("failed to import reference data: %v", err)
			}
		},
	}

	var pingDbCmd = &cobra.Command{
		Use:   "ping",
		Short: "Ping Postgres database",
		Run: func(cmd *cobra.Command, args []string) {
			cmds.PingDatabase()
		},
	}

	var migrationCmd = &cobra.Command{
		Use:   "migration [up|down] <migrations_path>",
		Short: "Run DB migration",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			cmds.RunMigration(args[0], args[1])
		},
	}

	var versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Show version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(versioninfo.Short())
		},
	}

	rootCmd.AddCommand(importGmlCmd)
	rootCmd.AddCommand(importRefDataCmd)
	rootCmd.AddCommand(pingDbCmd)
	rootCmd.AddCommand(migrationCmd)
	rootCmd.AddCommand(versionCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
