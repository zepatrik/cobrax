package cobrax

import (
	"bytes"
	"strings"
	"text/template"

	"github.com/spf13/cobra"
)

var usageTemplateFuncs = template.FuncMap{}

// AddUsageTemplateFunc adds a template function to the usage template.
func AddUsageTemplateFunc(name string, f interface{}) {
	usageTemplateFuncs[name] = f
}

const (
	helpTemplate = `{{insertTemplate . (or .Long .Short) | trimTrailingWhitespaces}}

{{if or .Runnable .HasSubCommands}}{{.UsageString}}{{end}}`
	usageTemplate = `Usage:{{if .Runnable}}
  {{.UseLine}}{{end}}{{if .HasAvailableSubCommands}}
  {{.CommandPath}} [command]{{end}}{{if gt (len .Aliases) 0}}

Aliases:
  {{.NameAndAliases}}{{end}}{{if .HasExample}}

Examples:
{{insertTemplate . .Example}}{{end}}{{if .HasAvailableSubCommands}}

Available Commands:{{range .Commands}}{{if (or .IsAvailableCommand (eq .Name "help"))}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableLocalFlags}}

Flags:
{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasAvailableInheritedFlags}}

Global Flags:
{{.InheritedFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasHelpSubCommands}}

Additional help topics:{{range .Commands}}{{if .IsAdditionalHelpTopicCommand}}
  {{rpad .CommandPath .CommandPathPadding}} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableSubCommands}}

Use "{{.CommandPath}} [command] --help" for more information about a command.{{end}}
`
)

// EnableUsageTemplating enables gotemplates for usage strings, i.e. cmd.Short, cmd.Long, and cmd.Example.
// The data for the template is the command itself. Especially useful are `.Root.Name` and `.CommandPath`.
// This will be inherited by all subcommands, so enabling it on the root command is sufficient.
func EnableUsageTemplating(cmds ...*cobra.Command) {
	cobra.AddTemplateFunc("insertTemplate", func(cmd *cobra.Command, tmpl string) (string, error) {
		t := template.New("")
		t.Funcs(usageTemplateFuncs)
		t, err := t.Parse(tmpl)
		if err != nil {
			return "", err
		}
		var out bytes.Buffer
		if err := t.Execute(&out, cmd); err != nil {
			return "", err
		}
		return out.String(), nil
	})
	for _, cmd := range cmds {
		cmd.SetHelpTemplate(helpTemplate)
		cmd.SetUsageTemplate(usageTemplate)
	}
}

// DisableUsageTemplating resets the commands usage template to the default.
// This can be used to undo the effects of EnableUsageTemplating, specifically for a subcommand.
func DisableUsageTemplating(cmds ...*cobra.Command) {
	defaultCmd := new(cobra.Command)
	for _, cmd := range cmds {
		cmd.SetHelpTemplate(defaultCmd.HelpTemplate())
		cmd.SetUsageTemplate(defaultCmd.UsageTemplate())
	}
}

// AssertUsageTemplates asserts that the usage string of the commands are properly templated.
func AssertUsageTemplates(t TestingT, cmd *cobra.Command) {
	var usage, help string
	func() {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("Usage template for %s panicked: %v", cmd.CommandPath(), r)
				t.FailNow()
			}
		}()
		usage = cmd.UsageString()

		out, err := cmd.OutOrStdout(), cmd.ErrOrStderr()
		bb := new(bytes.Buffer)

		cmd.SetOut(bb)
		cmd.SetErr(bb)
		if err := cmd.Help(); err != nil {
			t.Errorf("Help template for %s failed: %v", cmd.CommandPath(), err)
			t.FailNow()
		}
		help = bb.String()

		cmd.SetOut(out)
		cmd.SetErr(err)
	}()
	if strings.Contains(usage, "{{") {
		t.Errorf("Usage template for %s not properly templated: %s", cmd.CommandPath(), usage)
	}
	if strings.Contains(help, "{{") {
		t.Errorf("Help template for %s not properly templated: %s", cmd.CommandPath(), help)
	}
	for _, child := range cmd.Commands() {
		AssertUsageTemplates(t, child)
	}
}
