package main

import (
	"fmt"
	"os"
	"schnegge/internal/base"
	"schnegge/internal/cli"
	"schnegge/internal/client"
	"schnegge/internal/command"
	"schnegge/internal/config"
)

func main() {
	configCfg, err := config.ReadConfig()
	if err != nil {
		base.Log.Panic(err)
	}
	parameterCfg := config.NewHierarchyConfig(configCfg)
	cmd := command.ParseCommandLine(parameterCfg, os.Args)
	if b, _ := parameterCfg.GetValue(config.NoSplash); b != "true" {
		cli.PrintSplashScreen()
	}
	checkVerbose(parameterCfg)
	if cmd == "" {
		for {
			cmd, inputCfg := cli.Input(parameterCfg)
			checkVerbose(inputCfg)
			config.ReadConfig()
			excecuteCommand(inputCfg, cmd)
		}
	}
	excecuteCommand(parameterCfg, cmd)
}

func excecuteCommand(cfg config.Config, cmd string) {
	switch cmd {
	case "set":
		if err := config.WriteConfigFile(cfg); err != nil {
			fmt.Println("Fehler beim Schreiben der Konfiguration", err)
		}
	case "add":
		order, _ := cfg.GetValue(config.Order)
		record := client.NewDailyReportData(cfg, cli.GetEmployeeId(cfg, order))
		base.Log.Println("DailyReportData:", record)
		client.Record(cfg, record)
	case "list":
		cli.PrintOvertime(client.ReadOvertime(cfg))
		sortByProject, _ := cfg.GetValue(config.SortByProject)
		if sortByProject == "true" {
			cli.PrintByProject(client.ReadDailyReports(cfg))
		} else {
			cli.PrintDailyReports(client.ReadDailyReports(cfg))
		}
	case "exit":
		os.Exit(0)
	default:
		fmt.Println("Unknown command:", cmd)
	}
}

func checkVerbose(cfg config.Config) {
	if verbose, _ := cfg.GetValue(config.Verbose); verbose == "true" {
		command.PrintVerboseCommand(cfg)
		config.PrintVerboseConfig()
		base.EnableLogs()
	} else {
		base.DisableLogs()
	}
}
