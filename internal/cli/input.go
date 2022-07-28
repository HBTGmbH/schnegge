/*
Package cli (the interactive part) was inspired by https://github.com/c-bata/kube-prompt
*/
package cli

import (
	"fmt"
	"github.com/c-bata/go-prompt"
	"os"
	"schnegge/internal/command"
	"schnegge/internal/config"
	"strings"
)

var possibleCommands = []string{"set", "list", "add"}

var commands = []prompt.Suggest{
	{Text: "set", Description: "Setzt einen Konfigurationswert dauerhaft in der Datei ~/.schnegge"},
	{Text: "list", Description: "Zeigt deine aktuellen Reports"},
	{Text: "add", Description: "Fügt eine Buchung hinzu"},
	{Text: "exit", Description: "Beendet das Programm"},
}

var addOptions = []prompt.Suggest{
	{Text: "-auftrag", Description: "Der Auftrag auf den gebucht werden soll"},
	{Text: "-datum", Description: "Das Datum an welchem gebucht werden soll, z.B. 06.05.1981, heute, gestern, vorgestern, morgen, Montag ..."},
	{Text: "-stunden", Description: "Die Anzahl Stunden die gebucht werden soll, Default ist 0"},
	{Text: "-minuten", Description: "Die Anzahl Minuten die gebucht werden soll, Default ist 0"},
	{Text: "-fortbildung", Description: "Gibt an, ob eine Buchung Fortbildungscharackter hat, Default ist false"},
}

var listOptions = []prompt.Suggest{
	{Text: "-datum", Description: "Das Datum oder der Bereich welcher angezeigt werden soll, z.B. 06.05.1981, heute, gestern, vorgestern, morgen, Donnerstag, Januar, Woche"},
	{Text: "-nachAuftrag", Description: "Wenn angegeben wird die Ausgabe nach Aufträgen statt nach Datum sortiert"},
}

var helpOptions = []prompt.Suggest{
	{Text: "-h", Description: "Hier gibt es vielleicht Hilfe"},
	{Text: "-help", Description: "Hier gibt es vielleicht Hilfe"},
}

var setOptions = []prompt.Suggest{
	{Text: "-tokenID", Description: "Die TokenID aus Salat"},
	{Text: "-tokenSecret", Description: "Die TokenSecret aus Salat"},
	{Text: "-server", Description: "Der Salat Server"},
	{Text: "-noSplash", Description: "Wenn true, wird der Splashscreen nicht angezeigt"},
	{Text: "-auftrag", Description: "Der Auftrag auf den gebucht werden soll"},
}

var globalOptions = []prompt.Suggest{}

var datumValuesList = []prompt.Suggest{
	{Text: "Januar"},
	{Text: "Februar"},
	{Text: "März"},
	{Text: "April"},
	{Text: "Mai"},
	{Text: "Juni"},
	{Text: "Juli"},
	{Text: "August"},
	{Text: "September"},
	{Text: "Oktober"},
	{Text: "November"},
	{Text: "Dezember"},
	{Text: "Monat"},
	{Text: "Heute"},
	{Text: "Gestern"},
	{Text: "Vorgestern"},
	{Text: "Morgen"},
	{Text: "Montag"},
	{Text: "Dienstag"},
	{Text: "Mittwoch"},
	{Text: "Donnerstag"},
	{Text: "Freitag"},
	{Text: "Samstag"},
	{Text: "Sonntag"},
	{Text: "Woche"},
	{Text: "Vorwoche"},
	{Text: "-7"},
	{Text: "-14"},
}

var datumValuesAdd = []prompt.Suggest{
	{Text: "Montag"},
	{Text: "Dienstag"},
	{Text: "Mittwoch"},
	{Text: "Donnerstag"},
	{Text: "Freitag"},
	{Text: "Samstag"},
	{Text: "Sonntag"},
	{Text: "Heute"},
	{Text: "Gestern"},
	{Text: "Vorgestern"},
	{Text: "Morgen"},
	{Text: "-7"},
	{Text: "-14"},
}

// local copy of the config for all method in this file
var inputConfig *config.HierarchyConfig

// for history function
var oldLines []string

func InitInputConfig(cfg config.Config) {
	inputConfig = config.NewHierarchyConfig(cfg)
}

func inputCompleter(d prompt.Document) []prompt.Suggest {
	if d.TextBeforeCursor() == "" {
		return prompt.FilterHasPrefix(commands, "", true)
		//return []prompt.Suggest{}
	}
	args := splitLine(d.TextBeforeCursor())

	reparseLine(d, args)
	w := d.GetWordBeforeCursor()

	// If word before the cursor starts with "-", returns CLI flag options.
	if strings.HasPrefix(w, "-") {
		return optionCompleter(inputConfig, args, strings.HasPrefix(w, "--"))
	}

	// Return suggestions for option
	if suggests, found := completeOptionArguments(inputConfig, d); found {
		return suggests
	}
	if len(args) <= 1 {
		return prompt.FilterHasPrefix(commands, args[0], true)
	}
	return []prompt.Suggest{}
}

func reparseLine(d prompt.Document, args []string) {
	defer quiet()()
	if strings.HasSuffix(d.TextBeforeCursor(), " ") {
		InitInputConfig(inputConfig.BaseConfig)
		command.ParseCommandParameters(inputConfig, args)
	}
}

func splitLine(line string) []string {
	sequences := strings.Split(line, "\"")
	var result []string
	for n, s := range sequences {
		if n%2 == 1 {
			result = append(result, s)
		} else {
			result = append(result, strings.Split(strings.TrimSpace(s), " ")...)
		}
	}
	return result
}

func Input(baseCfg config.Config) (string, config.Config) {
	InitInputConfig(baseCfg)
	prompt.OptionShowCompletionAtStart()
	fmt.Println()
	t := prompt.Input("> ", inputCompleter,
		prompt.OptionHistory(oldLines),
		prompt.OptionTitle("Schnegge"),
		prompt.OptionAddKeyBind(prompt.KeyBind{Key: prompt.ControlC, Fn: func(buffer *prompt.Buffer) {
			os.Exit(0)
		}}))
	oldLines = append(oldLines, t)
	args := splitLine(t)
	cmd, _ := command.ParseCommandParameters(inputConfig, args)
	return cmd, inputConfig
}

func getPreviousOption(d prompt.Document) (cmd, option string, found bool) {
	args := strings.Split(d.TextBeforeCursor(), " ")
	l := len(args)
	if l >= 2 {
		option = args[l-2]
	}
	if strings.HasPrefix(option, "-") {
		return args[0], option, true
	}
	return "", "", false
}

func optionCompleter(cfg config.Config, args []string, long bool) []prompt.Suggest {
	l := len(args)
	if l <= 1 {
		if long {
			return prompt.FilterHasPrefix(helpOptions, "--", false)
		}
		return helpOptions
	}

	var suggests []prompt.Suggest
	commandArgs, _ := excludeOptions(args)
	switch commandArgs[0] {
	case "add":
		fetchOrders(cfg)
		suggests = addOptions
	case "list":
		suggests = listOptions
	case "set":
		suggests = setOptions
	default:
		suggests = helpOptions
	}

	suggests = append(suggests, globalOptions...)
	if long {
		return prompt.FilterContains(
			prompt.FilterHasPrefix(suggests, "--", false),
			strings.TrimLeft(args[l-1], "--"),
			true,
		)
	}
	return prompt.FilterContains(suggests, strings.TrimLeft(args[l-1], "-"), true)
}

// excludeOptions return all commands, i.e. if you use subcommands
func excludeOptions(args []string) ([]string, bool) {
	l := len(args)
	if l == 0 {
		return nil, false
	}
	cmd := args[0]
	filtered := make([]string, 0, l)

	var skipNextArg bool
	for i := 0; i < len(args); i++ {
		if skipNextArg {
			skipNextArg = false
			continue
		}

		if cmd == "logs" && args[i] == "-f" {
			continue
		}
		for _, s := range []string{
			"-auftrag",
			"-datum",
			"-minuten",
			"-stunden",
			"-tokenID",
			"-tokenSecret",
			"-server",
		} {
			if strings.HasPrefix(args[i], s) {
				if strings.Contains(args[i], "=") {
					// we can specify option value like '-o=json'
					skipNextArg = false
				} else {
					skipNextArg = true
				}
				continue
			}
		}
		if strings.HasPrefix(args[i], "-") {
			continue
		}

		filtered = append(filtered, args[i])
	}
	return filtered, skipNextArg
}

func completeOptionArguments(cfg config.Config, d prompt.Document) ([]prompt.Suggest, bool) {
	cmd, option, found := getPreviousOption(d)
	if !found {
		return []prompt.Suggest{}, false
	}
	switch cmd {
	case "list":
		if option == "-datum" || option == "-d" || option == "--datum" {
			return prompt.FilterHasPrefix(
				datumValuesList,
				d.GetWordBeforeCursor(),
				true,
			), true
		}
	case "add", "set":
		if option == "-datum" || option == "-d" || option == "--datum" {
			return prompt.FilterHasPrefix(
				datumValuesAdd,
				d.GetWordBeforeCursor(),
				true,
			), true
		} else if option == "-auftrag" || option == "-a" || option == "--auftrag" {
			return prompt.FilterContains(completeAuftraege(cfg), d.GetWordBeforeCursor(), true), true
		}
	}
	return []prompt.Suggest{}, false
}

func quiet() func() {
	null, _ := os.Open(os.DevNull)
	sout := os.Stdout
	serr := os.Stderr
	os.Stdout = null
	os.Stderr = null
	return func() {
		defer null.Close()
		os.Stdout = sout
		os.Stderr = serr
	}
}
