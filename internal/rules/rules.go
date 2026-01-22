package rules

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/common/types"
	"github.com/kkrypt0nn/aegisbot/internal/event"
	"github.com/kkrypt0nn/aegisbot/proto"
	"sigs.k8s.io/yaml"
)

type YAMLRule struct {
	Rule struct {
		Name          string       `yaml:"name"`
		Meta          RuleMeta     `yaml:"meta"`
		Strings       []RuleString `yaml:"strings"`
		Expression    string       `yaml:"expression"`
		AlertTemplate string       `yaml:"alertTemplate"`
		BanTemplate   string       `yaml:"banTemplate"`
		KickTemplate  string       `yaml:"kickTemplate"`
	} `yaml:"rule"`
}

type RuleMeta struct {
	Action     string `yaml:"action"`
	Context    string `yaml:"context"`
	IgnoreBots bool   `yaml:"ignoreBots"`
}

type RuleString struct {
	Name  string `yaml:"name"`
	Value string `yaml:"value"`
}

type SimplifiedRule struct {
	Name       string
	Action     string
	Context    string
	IgnoreBots bool

	Strings map[string]RuleString
	Program cel.Program

	AlertTemplate string
	BanTemplate   string
	KickTemplate  string
}

func Parse(filePath string) ([]*SimplifiedRule, error) {
	file, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	// TODO: yaml.JSONToYAML then yaml.Unmarshal
	var yamlRules []YAMLRule
	if err := yaml.Unmarshal(file, &yamlRules); err != nil {
		return nil, err
	}

	var simplifiedRules []*SimplifiedRule
	for _, yamlRule := range yamlRules {
		stringsMap := make(map[string]RuleString)
		for _, s := range yamlRule.Rule.Strings {
			stringsMap[s.Name] = s
		}

		var celVars []cel.EnvOption

		// Add the variables depending on the context
		switch yamlRule.Rule.Meta.Context {
		case "message":
			celVars = append(celVars,
				cel.Types(&proto.Message{}),
				cel.Variable("message", cel.ObjectType("proto.Message")),
			)
		case "member":
			celVars = append(celVars,
				cel.Types(&proto.Member{}),
				cel.Variable("member", cel.ObjectType("proto.Member")),
			)
		default:
			return nil, fmt.Errorf("unsupported rule context: %s", yamlRule.Rule.Meta.Context)
		}

		// Add the strings defined in the rule
		for _, s := range yamlRule.Rule.Strings {
			celVars = append(celVars, cel.Constant(s.Name, cel.StringType, types.String(s.Value)))
		}

		// Add the custom functions for either the types or as a global function
		celVars = append(celVars, cel.Function("isBot", isBotOverload))
		celVars = append(celVars, cel.Function("hasLinks", messageHasLinks))
		celVars = append(celVars, cel.Function("getLinks", messageGetLinks))

		env, err := cel.NewEnv(celVars...)
		if err != nil {
			return nil, err
		}

		ast, iss := env.Compile(strings.TrimSpace(yamlRule.Rule.Expression))
		if iss.Err() != nil {
			return nil, iss.Err()
		}

		program, err := env.Program(ast)
		if err != nil {
			return nil, err
		}

		rule := &SimplifiedRule{
			Name:       yamlRule.Rule.Name,
			Action:     yamlRule.Rule.Meta.Action,
			Context:    yamlRule.Rule.Meta.Context,
			IgnoreBots: yamlRule.Rule.Meta.IgnoreBots,

			Strings: stringsMap,
			Program: program,

			AlertTemplate: yamlRule.Rule.AlertTemplate,
			BanTemplate:   yamlRule.Rule.BanTemplate,
			KickTemplate:  yamlRule.Rule.KickTemplate,
		}
		simplifiedRules = append(simplifiedRules, rule)
	}

	return simplifiedRules, nil
}

func Load(dir string) ([]*SimplifiedRule, error) {
	var allRules []*SimplifiedRule
	err := filepath.Walk(dir, func(path string, fileInfo os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if fileInfo.IsDir() || (!strings.HasSuffix(fileInfo.Name(), ".yaml") && !strings.HasSuffix(fileInfo.Name(), ".yml")) {
			return nil
		}

		rules, err := Parse(path)
		if err != nil {
			return err
		}
		allRules = append(allRules, rules...)

		return nil
	})

	if err != nil {
		return nil, err
	}
	return allRules, nil
}

func (r *SimplifiedRule) Evaluate(ctx *event.Context) (bool, error) {
	if r.Context != string(ctx.Type) {
		return false, nil
	}
	if r.IgnoreBots && ctx.Member != nil && ctx.Member.Bot {
		return false, nil
	}

	input := map[string]any{}
	switch ctx.Type {
	case event.EventMessage:
		input["message"] = ctx.Message
		input["member"] = ctx.Member
	case event.EventMember:
		input["member"] = ctx.Member
	}

	out, _, err := r.Program.Eval(input)
	if err != nil {
		return false, err
	}

	result, ok := out.Value().(bool)
	if !ok {
		return false, fmt.Errorf("not a boolean result")
	}

	return result, nil
}
