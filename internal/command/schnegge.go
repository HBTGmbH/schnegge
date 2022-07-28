package command

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"schnegge/internal/base"
	"schnegge/internal/config"
	"strconv"
	"strings"
)

var listCmd *flag.FlagSet
var addCmd *flag.FlagSet
var setCmd *flag.FlagSet
var noCmd *flag.FlagSet

func ParseCommandLine(cfg config.Config, line []string) string {
	cmd, err := ParseCommandParameters(cfg, line[1:])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return cmd
}

func ParseCommandParameters(cfg config.Config, line []string) (string, error) {
	if len(line) < 1 {
		return "", nil
	}

	switch line[0] {
	case "set":
		setCmd = flag.NewFlagSet("set", flag.ContinueOnError)
		setCmd.Usage = func() {}
		err := parse(cfg, setCmd, line[1:])
		return "set", err
	case "add":
		addCmd = flag.NewFlagSet("add", flag.ContinueOnError)
		addCmd.Usage = func() {}
		err := parse(cfg, addCmd, line[1:])
		return "add", err
	case "list":
		listCmd = flag.NewFlagSet("list", flag.ContinueOnError)
		listCmd.Usage = func() {}
		err := parse(cfg, listCmd, line[1:])
		return "list", err
	case "exit":
		return "exit", nil
	default:
		if strings.HasPrefix(line[0], "-") {
			noCmd = flag.NewFlagSet("unknown", flag.ContinueOnError)
			noCmd.Usage = func() {}
			err := parse(cfg, noCmd, line[0:])
			return "", err
		} else {
			return "line[0]", errors.New("expected 'set', 'add', 'list' or 'exit' subcommand")
		}
	}
}

func parse(cfg config.Config, fs *flag.FlagSet, line []string) error {
	// init FlagSet
	// mitarbeiter := fs.String("mitarbeiter", "", "Der Mitarbeiter, für den das Command ausgeführt werden soll")
	dateVar := fs.String("datum", "", "Das Datum, für das das Command ausgeführt werden soll, special Values: heute, gestern, vorgestern, morgen, *Wochentag*, *DatumVar*")
	tokenIDVar := fs.String("tokenID", "", "Die TokenID aus Salat")
	tokenSecretVar := fs.String("tokenSecret", "", "Das TokenSecret aus Salat")
	verboseVar := fs.Bool("verbose", false, "Zeige Debug Informationen")
	minutesVar := fs.Int("minuten", 0, "Anzahl Minuten einer DailyReportData")
	hoursVar := fs.Int("stunden", 0, "Anzahl Stunden einer DailyReportData")
	trainingVar := fs.Bool("fortbildung", false, "Auf true setzen, wenn eine DailyReportData Fortbildungscharackter hat")
	orderVar := fs.String("auftrag", "", "Gibt einen Order/Suborder für eine DailyReportData an")
	serverVar := fs.String("server", "", "Gibt den Salat Server an, z.B. https://salat.hbt.de")
	sortByProjectVar := fs.Bool("nachAuftrag", false, "Gibt an, ob die Liste nach Auftrag statt nach Datum sortiert ausgegeben werden soll")
	noSplashVar := fs.Bool("noSplash", false, "Wenn true gesetzt, wird der Splashscreen nicht angezeigt")
	// parse
	err := fs.Parse(line)
	// copy Values to config
	addStringToConfig(cfg, config.Date, dateVar)
	addStringToConfig(cfg, config.TokenID, tokenIDVar)
	addStringToConfig(cfg, config.TokenSecret, tokenSecretVar)
	addBoolToConfig(cfg, config.Verbose, verboseVar)
	addIntToStringConfig(cfg, config.Minutes, minutesVar)
	addIntToStringConfig(cfg, config.Hours, hoursVar)
	addBoolToConfig(cfg, config.Training, trainingVar)
	addStringToConfig(cfg, config.Order, orderVar)
	addStringToConfig(cfg, config.Server, serverVar)
	addStringToConfig(cfg, config.Comment, parseArgs(fs))
	addBoolToConfig(cfg, config.SortByProject, sortByProjectVar)
	addBoolToConfig(cfg, config.NoSplash, noSplashVar)
	return err
}

func parseArgs(flagSet *flag.FlagSet) *string {
	values := flagSet.Args()
	result := strings.Join(values, " ")
	return &result
}

func addStringToConfig(cfg config.Config, key config.ConfigKey, value *string) {
	if *value != "" {
		cfg.AddValue(key, *value)
	}
}
func addBoolToConfig(cfg config.Config, key config.ConfigKey, value *bool) {
	if *value {
		cfg.AddValue(key, "true")
	}
}
func addIntToStringConfig(cfg config.Config, key config.ConfigKey, value *int) {
	if *value != 0 {
		cfg.AddValue(key, strconv.Itoa(*value))
	}
}

func PrintVerboseCommand(cfg config.Config) {
	base.Log.Println()
	base.Log.Println("=== COMMAND ===")
	date, _ := cfg.GetValue(config.Date)
	base.Log.Printf("Datum: %v\n", date)
	verbose, _ := cfg.GetValue(config.Verbose)
	base.Log.Printf("Verbose: %v\n", verbose)
	comment, _ := cfg.GetValue(config.Comment)
	base.Log.Printf("Kommentar: %v\n", comment)
}
