// Code generated by "make cli"; DO NOT EDIT.
package accountscmd

import (
	"errors"
	"fmt"

	"github.com/hashicorp/boundary/api"
	"github.com/hashicorp/boundary/api/accounts"
	"github.com/hashicorp/boundary/internal/cmd/base"
	"github.com/hashicorp/boundary/internal/cmd/common"
	"github.com/hashicorp/go-secure-stdlib/strutil"
	"github.com/mitchellh/cli"
	"github.com/posener/complete"
)

func initPasswordFlags() {
	flagsOnce.Do(func() {
		extraFlags := extraPasswordActionsFlagsMapFunc()
		for k, v := range extraFlags {
			flagsPasswordMap[k] = append(flagsPasswordMap[k], v...)
		}
	})
}

var (
	_ cli.Command             = (*PasswordCommand)(nil)
	_ cli.CommandAutocomplete = (*PasswordCommand)(nil)
)

type PasswordCommand struct {
	*base.Command

	Func string

	plural string

	extraPasswordCmdVars
}

func (c *PasswordCommand) AutocompleteArgs() complete.Predictor {
	initPasswordFlags()
	return complete.PredictAnything
}

func (c *PasswordCommand) AutocompleteFlags() complete.Flags {
	initPasswordFlags()
	return c.Flags().Completions()
}

func (c *PasswordCommand) Synopsis() string {
	if extra := extraPasswordSynopsisFunc(c); extra != "" {
		return extra
	}

	synopsisStr := "account"

	synopsisStr = fmt.Sprintf("%s %s", "password-type", synopsisStr)

	return common.SynopsisFunc(c.Func, synopsisStr)
}

func (c *PasswordCommand) Help() string {
	initPasswordFlags()

	var helpStr string
	helpMap := common.HelpMap("account")

	switch c.Func {

	default:

		helpStr = c.extraPasswordHelpFunc(helpMap)

	}

	// Keep linter from complaining if we don't actually generate code using it
	_ = helpMap
	return helpStr
}

var flagsPasswordMap = map[string][]string{

	"create": {"auth-method-id", "name", "description"},

	"update": {"id", "name", "description", "version"},
}

func (c *PasswordCommand) Flags() *base.FlagSets {
	if len(flagsPasswordMap[c.Func]) == 0 {
		return c.FlagSet(base.FlagSetNone)
	}

	set := c.FlagSet(base.FlagSetHTTP | base.FlagSetClient | base.FlagSetOutputFormat)
	f := set.NewFlagSet("Command Options")
	common.PopulateCommonFlags(c.Command, f, "password-type account", flagsPasswordMap, c.Func)

	extraPasswordFlagsFunc(c, set, f)

	return set
}

func (c *PasswordCommand) Run(args []string) int {
	initPasswordFlags()

	switch c.Func {
	case "":
		return cli.RunResultHelp

	}

	c.plural = "password-type account"
	switch c.Func {
	case "list":
		c.plural = "password-type accounts"
	}

	f := c.Flags()

	if err := f.Parse(args); err != nil {
		c.PrintCliError(err)
		return base.CommandUserError
	}

	if strutil.StrListContains(flagsPasswordMap[c.Func], "id") && c.FlagId == "" {
		c.PrintCliError(errors.New("ID is required but not passed in via -id"))
		return base.CommandUserError
	}

	var opts []accounts.Option

	if strutil.StrListContains(flagsPasswordMap[c.Func], "auth-method-id") {
		switch c.Func {

		case "create":
			if c.FlagAuthMethodId == "" {
				c.PrintCliError(errors.New("AuthMethod ID must be passed in via -auth-method-id or BOUNDARY_AUTH_METHOD_ID"))
				return base.CommandUserError
			}

		}
	}

	client, err := c.Client()
	if err != nil {
		c.PrintCliError(fmt.Errorf("Error creating API client: %s", err.Error()))
		return base.CommandCliError
	}
	accountsClient := accounts.NewClient(client)

	switch c.FlagName {
	case "":
	case "null":
		opts = append(opts, accounts.DefaultName())
	default:
		opts = append(opts, accounts.WithName(c.FlagName))
	}

	switch c.FlagDescription {
	case "":
	case "null":
		opts = append(opts, accounts.DefaultDescription())
	default:
		opts = append(opts, accounts.WithDescription(c.FlagDescription))
	}

	if c.FlagFilter != "" {
		opts = append(opts, accounts.WithFilter(c.FlagFilter))
	}

	var version uint32

	switch c.Func {

	case "update":
		switch c.FlagVersion {
		case 0:
			opts = append(opts, accounts.WithAutomaticVersioning(true))
		default:
			version = uint32(c.FlagVersion)
		}

	}

	if ok := extraPasswordFlagsHandlingFunc(c, f, &opts); !ok {
		return base.CommandUserError
	}

	var result api.GenericResult

	switch c.Func {

	case "create":
		result, err = accountsClient.Create(c.Context, c.FlagAuthMethodId, opts...)

	case "update":
		result, err = accountsClient.Update(c.Context, c.FlagId, version, opts...)

	}

	result, err = executeExtraPasswordActions(c, result, err, accountsClient, version, opts)

	if err != nil {
		if apiErr := api.AsServerError(err); apiErr != nil {
			var opts []base.Option

			c.PrintApiError(apiErr, fmt.Sprintf("Error from controller when performing %s on %s", c.Func, c.plural), opts...)
			return base.CommandApiError
		}
		c.PrintCliError(fmt.Errorf("Error trying to %s %s: %s", c.Func, c.plural, err.Error()))
		return base.CommandCliError
	}

	output, err := printCustomPasswordActionOutput(c)
	if err != nil {
		c.PrintCliError(err)
		return base.CommandUserError
	}
	if output {
		return base.CommandSuccess
	}

	switch c.Func {

	}

	switch base.Format(c.UI) {
	case "table":
		c.UI.Output(printItemTable(result))

	case "json":
		if ok := c.PrintJsonItem(result); !ok {
			return base.CommandCliError
		}
	}

	return base.CommandSuccess
}

var (
	extraPasswordActionsFlagsMapFunc = func() map[string][]string { return nil }
	extraPasswordSynopsisFunc        = func(*PasswordCommand) string { return "" }
	extraPasswordFlagsFunc           = func(*PasswordCommand, *base.FlagSets, *base.FlagSet) {}
	extraPasswordFlagsHandlingFunc   = func(*PasswordCommand, *base.FlagSets, *[]accounts.Option) bool { return true }
	executeExtraPasswordActions      = func(_ *PasswordCommand, inResult api.GenericResult, inErr error, _ *accounts.Client, _ uint32, _ []accounts.Option) (api.GenericResult, error) {
		return inResult, inErr
	}
	printCustomPasswordActionOutput = func(*PasswordCommand) (bool, error) { return false, nil }
)
