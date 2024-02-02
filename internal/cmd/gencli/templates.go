package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"github.com/hashicorp/boundary/internal/types/resource"
	"github.com/hashicorp/go-secure-stdlib/strutil"
	"github.com/iancoleman/strcase"
)

func fillTemplates() {
	for pkg, structs := range inputStructs {
		for _, data := range structs {
			outBuf := new(bytes.Buffer)

			if err := cmdTemplate.Execute(outBuf, data); err != nil {
				fmt.Printf("error executing struct template for resource %s: %v\n", pkg, err)
				os.Exit(1)
			}

			fName := pkg
			if data.SubActionPrefix != "" {
				fName = fmt.Sprintf("%s_%s", strcase.ToKebab(data.SubActionPrefix), fName)
			}
			outFile, err := filepath.Abs(fmt.Sprintf("%s/%scmd/%s.gen.go", os.Getenv("CLI_GEN_BASEPATH"), pkg, fName))
			if err != nil {
				fmt.Printf("error opening file for package %s: %v\n", pkg, err)
				os.Exit(1)
			}
			outDir := filepath.Dir(outFile)
			if _, err := os.Stat(outDir); os.IsNotExist(err) {
				_ = os.Mkdir(outDir, os.ModePerm)
			}
			if err := os.WriteFile(outFile, outBuf.Bytes(), 0o644); err != nil {
				fmt.Printf("error writing file %q: %v\n", outFile, err)
				os.Exit(1)
			}
		}
	}
}

func camelCase(in string) string {
	return strcase.ToCamel(in)
}

func envCase(in string) string {
	return strcase.ToScreamingSnake(in)
}

func snakeCase(in string) string {
	return strcase.ToSnake(in)
}

func kebabCase(in string) string {
	return strcase.ToKebab(in)
}

func lowerSpaceCase(in string) string {
	return strcase.ToDelimited(in, ' ')
}

func hasAction(in []string, action string) bool {
	return strutil.StrListContains(in, action)
}

func pluralString(in string) string {
	return lowerSpaceCase(resource.Map[in].PluralString())
}

var cmdTemplate = template.Must(template.New("").Funcs(
	template.FuncMap{
		"camelCase":      camelCase,
		"envCase":        envCase,
		"snakeCase":      snakeCase,
		"kebabCase":      kebabCase,
		"lowerSpaceCase": lowerSpaceCase,
		"hasAction":      hasAction,
		"pluralString":   pluralString,
	},
).Parse(`{{ $input := . }}
// Code generated by "make cli"; DO NOT EDIT.
// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package {{ .Pkg }}cmd

import (
	"errors"
	"fmt"
	"net/http"
	"sync"

	"github.com/hashicorp/boundary/api"
	"github.com/hashicorp/boundary/api/{{ .Pkg }}"
	"github.com/hashicorp/boundary/internal/cmd/base"
	"github.com/hashicorp/boundary/internal/cmd/common"
	"github.com/hashicorp/go-secure-stdlib/strutil"
	"github.com/mitchellh/cli"
	"github.com/posener/complete"
)

func init{{ camelCase .SubActionPrefix }}Flags() {
	flagsOnce.Do(func() {
		extraFlags := extra{{ camelCase .SubActionPrefix }}ActionsFlagsMapFunc()
		for k, v := range extraFlags {
			flags{{ camelCase .SubActionPrefix }}Map[k] = append(flags{{ camelCase .SubActionPrefix }}Map[k], v...)
		}
	})
}

var (
	_ cli.Command             = (*{{ camelCase .SubActionPrefix }}Command)(nil)
	_ cli.CommandAutocomplete = (*{{ camelCase .SubActionPrefix }}Command)(nil)
)

type {{ camelCase .SubActionPrefix }}Command struct {
	*base.Command

	Func string

	plural string
	{{ if .HasExtraCommandVars }}
	extra{{ camelCase .SubActionPrefix }}CmdVars
	{{ end }}
}

func (c *{{ camelCase .SubActionPrefix }}Command) AutocompleteArgs() complete.Predictor {
	init{{ camelCase .SubActionPrefix }}Flags()
	return complete.PredictAnything
}

func (c *{{ camelCase .SubActionPrefix }}Command) AutocompleteFlags() complete.Flags {
	init{{ camelCase .SubActionPrefix }}Flags()
	return c.Flags().Completions()
}

func (c *{{ camelCase .SubActionPrefix }}Command) Synopsis() string {
	if extra := extra{{ camelCase .SubActionPrefix }}SynopsisFunc(c); extra != "" {
		return extra
	}

	synopsisStr := "{{ lowerSpaceCase .ResourceType }}"
	{{ if .SubActionPrefix }}
	synopsisStr = fmt.Sprintf("%s %s", "{{ kebabCase .SubActionPrefix }}-type", synopsisStr)
	{{ end }}
	return common.SynopsisFunc(c.Func, synopsisStr)
}

func (c *{{ camelCase .SubActionPrefix }}Command) Help() string {
	init{{ camelCase .SubActionPrefix }}Flags()
	
	var helpStr string
	helpMap := common.HelpMap("{{ lowerSpaceCase .ResourceType }}")

	switch c.Func {
	{{ if not .SkipNormalHelp }}
	{{ range $i, $action := .StdActions }}
	case "{{ $action }}":
		helpStr = helpMap[c.Func]() + c.Flags().Help()
	{{ end }}
	{{ end }}
	default:
	{{ if .HasExtraHelpFunc }}
		helpStr = c.extra{{ camelCase .SubActionPrefix }}HelpFunc(helpMap)
	{{ else }}
		helpStr = helpMap["base"]()
	{{ end }}
	}

	// Keep linter from complaining if we don't actually generate code using it
	_ = helpMap
	return helpStr
}

var flags{{ camelCase .SubActionPrefix }}Map = map[string][]string{
	{{ with $attrFlags := ", \"attributes\", \"attr\", \"string-attr\", \"bool-attr\", \"num-attr\"" }}
	{{ with $secretFlags := ", \"secrets\", \"secret\", \"string-secret\", \"bool-secret\", \"num-secret\"" }}
	{{ with $objectFlags := ", \"object\", \"kv\", \"string-kv\", \"bool-kv\", \"num-kv\"" }}
	{{ range $i, $action := $input.StdActions }}
	{{ if eq $action "create" }}
	"create": { "{{ kebabCase $input.Container }}-id", "name", "description" {{ if $input.IsPluginType }} , "plugin-id", "plugin-name" {{ end }} {{ if $input.HasGenericAttributes }} {{ $attrFlags }} {{ end }} {{ if $input.HasGenericSecrets }} {{ $secretFlags }} {{ end }} {{ if $input.HasJsonObject }} {{ $objectFlags }} {{ end }} },
	{{ end }}
	{{ if eq $action "read" }}
	"read": {"id"},
	{{ end }}
	{{ if eq $action "update" }}
	"update": {"id", "name", "description" {{ if hasAction $input.VersionedActions "update" }}, "version" {{ end }} {{ if $input.HasGenericAttributes }} {{ $attrFlags }} {{ end }} {{ if $input.HasGenericSecrets }} {{ $secretFlags }} {{ end }} {{ if $input.HasJsonObject }} {{ $objectFlags }} {{ end }} },
	{{ end }}
	{{ if eq $action "delete" }}
	"delete": {"id"},
	{{ end }}
	{{ if eq $action "list" }}
	"list": { "{{ kebabCase $input.Container }}-id"{{ if not $input.SkipFiltering }},  "filter" {{ end }} {{ if (eq $input.Container "Scope") }}, "recursive"{{ end }} },
	{{ end }}
	{{ end }}
	{{ end }}
	{{ end }}
	{{ end }}
}

func (c *{{ camelCase .SubActionPrefix }}Command) Flags() *base.FlagSets {
	if len(flags{{ camelCase .SubActionPrefix }}Map[c.Func]) == 0 {
		return c.FlagSet(base.FlagSetNone)
	}

	set := c.FlagSet(base.FlagSetHTTP | base.FlagSetClient | base.FlagSetOutputFormat)
	f := set.NewFlagSet("Command Options")
	common.PopulateCommonFlags(c.Command, f, "{{ if .SubActionPrefix }}{{ kebabCase .SubActionPrefix }}-type {{ end }}{{ lowerSpaceCase .ResourceType }}", flags{{ camelCase .SubActionPrefix }}Map, c.Func)

	{{ if .HasGenericAttributes }}
	f = set.NewFlagSet("Attribute Options")
		attrsInput := common.CombinedSliceFlagValuePopulationInput{
		FlagSet: f,
		FlagNames: flags{{ camelCase .SubActionPrefix }}Map[c.Func],
		FullPopulationFlag: &c.FlagAttributes,
		FullPopulationInputName: "attributes",
		PiecewisePopulationFlag: &c.FlagAttrs,
		PiecewisePopulationInputBaseName: "attr",
	}
	common.PopulateCombinedSliceFlagValue(attrsInput)
	{{ end }}

	{{ if .HasGenericSecrets }}
	f = set.NewFlagSet("Secrets Options")
	scrtsInput := common.CombinedSliceFlagValuePopulationInput{
		FlagSet: f,
		FlagNames: flags{{ camelCase .SubActionPrefix }}Map[c.Func],
		FullPopulationFlag: &c.FlagSecrets,
		FullPopulationInputName: "secrets",
		PiecewisePopulationFlag: &c.FlagScrts,
		PiecewisePopulationInputBaseName: "secret",
	}
	common.PopulateCombinedSliceFlagValue(scrtsInput)
	{{ end }}

	{{ if .HasJsonObject }}
	f = set.NewFlagSet("Object Options")
		objsInput := common.CombinedSliceFlagValuePopulationInput{
		FlagSet: f,
		FlagNames: flags{{ camelCase .SubActionPrefix }}Map[c.Func],
		FullPopulationFlag: &c.FlagObject,
		FullPopulationInputName: "object",
		PiecewisePopulationFlag: &c.FlagKv,
		PiecewisePopulationInputBaseName: "kv",
	}
	common.PopulateCombinedSliceFlagValue(objsInput)
	{{ end }}

	extra{{ camelCase .SubActionPrefix }}FlagsFunc(c, set, f)

	return set
}

func (c *{{ camelCase .SubActionPrefix }}Command) Run(args []string) int {
	init{{ camelCase .SubActionPrefix }}Flags()

	{{ if .HasExampleCliOutput }}
	if os.Getenv("BOUNDARY_EXAMPLE_CLI_OUTPUT") != "" {
		c.UI.Output(exampleOutput())
		return base.CommandSuccess
	}
	{{ end }}

	switch c.Func {
	case "":
		return cli.RunResultHelp
	{{ if (not (hasAction .StdActions "create" ) )}}
	case "create":
		return cli.RunResultHelp
	{{ end }}
	{{ if (not (hasAction .StdActions "update" ) )}}
	case "update":
		return cli.RunResultHelp
	{{ end }}
	}

	c.plural = "{{ if .SubActionPrefix }}{{ kebabCase .SubActionPrefix }}-type {{ end }}{{ lowerSpaceCase .ResourceType }}"
	switch c.Func {
	case "list":
		c.plural = "{{ if .SubActionPrefix }}{{ kebabCase .SubActionPrefix }}-type {{ end }}{{ pluralString .ResourceType }}"
	}

	f := c.Flags()

   {{ if .UsesAlias }}
		args = c.ExtractAliasFromArgs(args)

	if c.FlagAlias != ""{
		if c.{{ .AliasFieldFlag }} != "" {
			c.PrintCliError(errors.New("Cannot specify both an alias and {{ .AliasFieldFlagName }}; choose one or the other"))
			return base.CommandUserError
		}
		c.{{ .AliasFieldFlag }} = c.FlagAlias
   }
   {{ end }}

	if err := f.Parse(args); err != nil {
		c.PrintCliError(err)
		return base.CommandUserError
	}

	{{ if .HasId }}
	if strutil.StrListContains(flags{{ camelCase .SubActionPrefix }}Map[c.Func], "id") && c.FlagId == "" {
			c.PrintCliError(errors.New("ID is required but not passed in via -id"))
			return base.CommandUserError
	}
	{{ end }}

	{{ if .HasGenericAttributes }}
	if c.FlagAttributes != "" && len(c.FlagAttrs) > 0 {
		c.PrintCliError(errors.New("-attributes flag cannot be used along with the following flags: attr, bool-attr, num-attr, string-attr"))
		return base.CommandUserError
	}
	{{ end }}

	{{ if .HasGenericSecrets }}
	if c.FlagSecrets != "" && len(c.FlagScrts) > 0 {
		c.PrintCliError(errors.New("-secrets flag cannot be used along with the following flags: secret, bool-secret, num-secret, string-secret"))
		return base.CommandUserError
	}
	{{ end }}

	{{ if .HasJsonObject }}
	if c.FlagObject != "" && len(c.FlagKv) > 0 {
		c.PrintCliError(errors.New("-object flag cannot be used along with the following flags: kv, bool-kv, num-kv, string-kv"))
		return base.CommandUserError
	}
	{{ end }}

	var opts []{{ .Pkg }}.Option

	if strutil.StrListContains(flags{{ camelCase .SubActionPrefix }}Map[c.Func], "{{ kebabCase .Container }}-id") {
		switch c.Func {
		{{ if hasAction .StdActions "create" }}
		case "create":
			if c.Flag{{ .Container }}Id == "" {
				c.PrintCliError(errors.New("{{ .Container }} ID must be passed in via -{{ kebabCase .Container }}-id or BOUNDARY_{{ envCase .Container }}_ID"))
				return base.CommandUserError
			}
		{{ end }}
		{{ if hasAction .StdActions "list" }}
		case "list":
			if c.Flag{{ .Container }}Id == "" {
				c.PrintCliError(errors.New("{{ .Container }} ID must be passed in via -{{ kebabCase .Container }}-id or BOUNDARY_{{ envCase .Container }}_ID"))
				return base.CommandUserError
			}
		{{ end }}
		}
	}

	client, err := c.Client()
	if c.WrapperCleanupFunc != nil {
		defer func() {
			if err := c.WrapperCleanupFunc(); err != nil {
				c.PrintCliError(fmt.Errorf("Error cleaning kms wrapper: %w", err))
			}
		}()
	}
	if err != nil {
		c.PrintCliError(fmt.Errorf("Error creating API client: %w", err))
		return base.CommandCliError
	}
	{{ .Pkg }}Client := {{ .Pkg }}.NewClient(client)

	{{ if .HasName }}
	switch c.FlagName {
	case "":
	case "null":
		opts = append(opts, {{ .Pkg }}.DefaultName())
	default:
		opts = append(opts, {{ .Pkg }}.WithName(c.FlagName))
	}
	{{ end }}

	{{ if .HasDescription }}
	switch c.FlagDescription {
	case "":
	case "null":
		opts = append(opts, {{ .Pkg }}.DefaultDescription())
	default:
		opts = append(opts, {{ .Pkg }}.WithDescription(c.FlagDescription))
	}
	{{ end }}

	{{ if (eq .Container "Scope") }}
	switch c.FlagRecursive {
	case true:
		opts = append(opts, {{ .Pkg }}.WithRecursive(true))
	}
	{{ end }}

	{{ if not .SkipFiltering }}
	if c.FlagFilter != "" {
		opts = append(opts, {{ .Pkg }}.WithFilter(c.FlagFilter))
	}
	{{ end }}

	{{ if .HasScopeName }}
	switch c.FlagScopeName {
	case "":
	default:
		opts = append(opts, {{ .Pkg }}.WithScopeName(c.FlagScopeName))
	}
	{{ end }}
	
	{{ if .IsPluginType }}
	switch c.FlagPluginId {
	case "":
	default:
		opts = append(opts, {{ .Pkg }}.WithPluginId(c.FlagPluginId))
	}
	switch c.FlagPluginName {
	case "":
	default:
		opts = append(opts, {{ .Pkg }}.WithPluginName(c.FlagPluginName))
	}
	{{ end }}

	var version uint32
	{{ if .VersionedActions }}
	switch c.Func {
	{{ range $i, $action := .VersionedActions }}
	case "{{ $action }}":
		switch c.FlagVersion {
		case 0:
			opts = append(opts, {{ $input.Pkg }}.WithAutomaticVersioning(true))
		default:
			version = uint32(c.FlagVersion)
		}
	{{ end }}
	}
	{{ end }}

	{{ if .HasGenericAttributes }}
	if err := common.HandleAttributeFlags(
		c.Command,
		"attr",
		c.FlagAttributes,
		c.FlagAttrs,
		func() {
			opts = append(opts, {{ .Pkg }}.DefaultAttributes())
		},
		func(in map[string]interface{}) {
			opts = append(opts, {{ .Pkg }}.WithAttributes(in))
		}); err != nil {
		c.PrintCliError(fmt.Errorf("Error evaluating attribute flags to: %s", err.Error()))
		return base.CommandCliError
	}
	{{ end }}

	{{ if .HasGenericSecrets }}
	if err := common.HandleAttributeFlags(
		c.Command,
		"secret",
		c.FlagSecrets,
		c.FlagScrts,
		func() {
			opts = append(opts, {{ .Pkg }}.DefaultSecrets())
		},
		func(in map[string]interface{}) {
			opts = append(opts, {{ .Pkg }}.WithSecrets(in))
		}); err != nil {
		c.PrintCliError(fmt.Errorf("Error evaluating secret flags to: %s", err.Error()))
		return base.CommandCliError
	}
	{{ end }}

	{{ if .HasJsonObject }}
	if err := common.HandleAttributeFlags(
		c.Command,
		"kv",
		c.FlagObject,
		c.FlagKv,
		func() {},
		func(in map[string]interface{}) {
			opts = append(opts, {{ .Pkg }}.WithJsonCredentialObject(in))
		}); err != nil {
		c.PrintCliError(fmt.Errorf("Error evaluating kv flags to: %s", err.Error()))
		return base.CommandCliError
	}
	{{ end }}

	if ok := extra{{ camelCase .SubActionPrefix }}FlagsHandlingFunc(c, f, &opts, ); !ok {
		return base.CommandUserError
	}

	var resp *api.Response
	var item *{{ $input.Pkg }}.{{ camelCase $input.ResourceType }}
	{{ if hasAction .StdActions "list" }}
	var items []*{{ $input.Pkg }}.{{ camelCase $input.ResourceType }}
	{{ end }}
	{{ range $i, $action := $input.StdActions }}
	{{ if ( not ( hasAction $input.SkipClientCallActions $action) ) }}
	var {{ $action }}Result *{{ $input.Pkg }}.{{ camelCase $input.ResourceType }}{{ camelCase $action }}Result
	{{ end }}
	{{ end }}

	switch c.Func {
	{{ range $i, $action := $input.StdActions }}
	{{ if eq $action "create" }}
	{{ if ( not ( hasAction $input.SkipClientCallActions "create") ) }}
	case "create":
		{{ $action }}Result, err = {{ $input.Pkg }}Client.Create(c.Context, {{ if (and $input.SubActionPrefix $input.NeedsSubtypeInCreate) }}"{{ $input.SubActionPrefix }}",{{ end }} c.Flag{{ $input.Container }}Id, opts...)
		if exitCode := c.checkFuncError(err); exitCode > 0 {
			return exitCode
		}
		resp = {{ $action }}Result.GetResponse()
		item = {{ $action }}Result.GetItem()
	{{ end }}
	{{ end }}
	{{ if eq $action "read" }}
	case "read":
		{{ $action }}Result, err = {{ $input.Pkg }}Client.Read(c.Context, c.FlagId, opts...)
		if exitCode := c.checkFuncError(err); exitCode > 0 {
			return exitCode
		}
		resp = {{ $action }}Result.GetResponse()
		item = {{ $action }}Result.GetItem()
	{{ end }}
	{{ if eq $action "update" }}
	case "update":
		{{ $action }}Result, err = {{ $input.Pkg }}Client.Update(c.Context, c.FlagId, version, opts...)
		if exitCode := c.checkFuncError(err); exitCode > 0 {
			return exitCode
		}
		resp = {{ $action }}Result.GetResponse()
		item = {{ $action }}Result.GetItem()
	{{ end }}
	{{ if eq $action "delete" }}
	case "delete":
		{{ $action }}Result, err = {{ $input.Pkg }}Client.Delete(c.Context, c.FlagId, opts...)
		if exitCode := c.checkFuncError(err); exitCode > 0 {
			return exitCode
		}
		resp = {{ $action }}Result.GetResponse()
	{{ end }}
	{{ if eq $action "list" }}
	case "list":
		{{ $action }}Result, err = {{ $input.Pkg}}Client.List(c.Context, c.Flag{{ $input.Container }}Id, opts...)
		if exitCode := c.checkFuncError(err); exitCode > 0 {
			return exitCode
		}
		resp = {{ $action }}Result.GetResponse()
		items = {{ $action }}Result.GetItems()
	{{ end }}
	{{ end }}
	}

	resp, item, {{ if hasAction .StdActions "list" }}items, {{ end }}err = executeExtra{{ camelCase .SubActionPrefix }}Actions(c, resp, item,  {{ if hasAction .StdActions "list" }}items, {{ end }}err, {{ .Pkg }}Client, version, opts)
	if exitCode := c.checkFuncError(err); exitCode > 0 {
		return exitCode
	}

	output, err := printCustom{{ camelCase .SubActionPrefix }}ActionOutput(c)
	if err != nil {
		c.PrintCliError(err)
		return base.CommandUserError
	}
	if output {
		return base.CommandSuccess
	}

	switch c.Func {
	{{ range $i, $action := .StdActions }}
	{{ if eq $action "delete" }}
	case "delete":
		switch base.Format(c.UI) {
		case "json":
			if ok := c.PrintJsonItem(resp); !ok {
				return base.CommandCliError
			}

		case "table":
			c.UI.Output("The delete operation completed successfully.")
		}

		return base.CommandSuccess
	{{ end }}
	{{ if eq $action "list" }}
	case "list":
		switch base.Format(c.UI) {
		case "json":
			if ok := c.PrintJsonItems(resp); !ok {
				return base.CommandCliError
			}

		case "table":
			c.UI.Output(c.printListTable(items))
		}

		return base.CommandSuccess
	{{ end }}
	{{ end }}
	}

	switch base.Format(c.UI) {
	case "table":
		c.UI.Output(printItemTable(item, resp))

	case "json":
		if ok := c.PrintJsonItem(resp); !ok {
			return base.CommandCliError
		}
	}

	return base.CommandSuccess
}

func (c *{{ camelCase .SubActionPrefix }}Command) checkFuncError(err error) int {
	if err == nil {
		return 0
	}
	if apiErr := api.AsServerError(err); apiErr != nil {
		c.PrintApiError(apiErr, fmt.Sprintf("Error from controller when performing %s on %s", c.Func, c.plural))
		return base.CommandApiError
	}
	c.PrintCliError(fmt.Errorf("Error trying to %s %s: %s", c.Func, c.plural, err.Error()))
	return base.CommandCliError
}

var (
	{{ if not .SubActionPrefix }}
	flagsOnce = new(sync.Once)
	{{ end }}
	extra{{ camelCase .SubActionPrefix }}ActionsFlagsMapFunc = func() map[string][]string { return nil }
	extra{{ camelCase .SubActionPrefix }}SynopsisFunc = func(*{{ camelCase .SubActionPrefix }}Command) string { return "" }
	extra{{ camelCase .SubActionPrefix }}FlagsFunc = func(*{{ camelCase .SubActionPrefix }}Command, *base.FlagSets, *base.FlagSet) {}
	extra{{ camelCase .SubActionPrefix }}FlagsHandlingFunc = func(*{{ camelCase .SubActionPrefix }}Command, *base.FlagSets, *[]{{ .Pkg }}.Option) bool { return true }
	executeExtra{{ camelCase .SubActionPrefix }}Actions = func(_ *{{ camelCase .SubActionPrefix }}Command, inResp *api.Response, inItem *{{ $input.Pkg }}.{{ camelCase $input.ResourceType }}, {{ if hasAction .StdActions "list" }}inItems []*{{ $input.Pkg }}.{{ camelCase $input.ResourceType }}, {{ end }}inErr error, _ *{{ .Pkg }}.Client, _ uint32, _ []{{ .Pkg }}.Option) (*api.Response, *{{ $input.Pkg }}.{{ camelCase $input.ResourceType }}, {{ if hasAction .StdActions "list" }}[]*{{ $input.Pkg }}.{{ camelCase $input.ResourceType }}, {{ end }}error) {
		return inResp, inItem, {{ if hasAction .StdActions "list" }}inItems, {{ end }}inErr
	}
	printCustom{{ camelCase .SubActionPrefix }}ActionOutput = func(*{{ camelCase .SubActionPrefix }}Command) (bool, error) { return false, nil }
)
`))
