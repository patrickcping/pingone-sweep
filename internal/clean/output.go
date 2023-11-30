package clean

import (
	"fmt"

	"github.com/fatih/color"
)

var (
	test map[string]string
)

func handleOutput(configKey string, output CleanOutput, dryRun bool) {

	configKeyFormat := color.New(color.FgBlue, color.Bold).SprintFunc()
	printString := configKeyFormat(configKey)

	if dryRun {
		dryRunFormat := color.New(color.FgMagenta).SprintFunc()
		printString = fmt.Sprintf("%s %s%s%s", printString, configKeyFormat("("), dryRunFormat("DRY RUN"), configKeyFormat(")"))
	}

	if output.ConfigItem.IdentifierToEvaluate != "" {
		printString = fmt.Sprintf("%s - %s (%s)", printString, output.ConfigItem.IdentifierToEvaluate, output.ConfigItem.Id)
	} else {
		printString = fmt.Sprintf("%s - %s", printString, output.ConfigItem.Id)
	}

	printString = fmt.Sprintf("%s with action %s", printString, output.Action)

	switch output.Result {
	case ENUMCLEANOUTPUTRESULT_SUCCESS:
		printString = fmt.Sprintf("%s - %s", printString, color.GreenString("Success"))
	case ENUMCLEANOUTPUTRESULT_NOACTION_OK:
		printString = fmt.Sprintf("%s - %s", printString, color.GreenString("No action taken"))
	case ENUMCLEANOUTPUTRESULT_NOACTION_WARN:
		printString = fmt.Sprintf("%s - %s", printString, color.YellowString("No action taken (needs review)"))
	case ENUMCLEANOUTPUTRESULT_FAILURE:
		printString = fmt.Sprintf("%s - %s", printString, color.RedString("Request Failure"))
	}

	if output.Message != nil && *output.Message != "" {
		printString = fmt.Sprintf("%s: %s", printString, *output.Message)
	}

	fmt.Printf("%s\n", printString)
}

func GetJSON() map[string]string {
	return test
}

func Init() {
	test = make(map[string]string)
}
