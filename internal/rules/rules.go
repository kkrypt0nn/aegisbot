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
		Name       string       `yaml:"name"`
		Meta       RuleMeta     `yaml:"meta"`
		Strings    []RuleString `yaml:"strings"`
		Expression string       `yaml:"expression"`
		Action     RuleAction   `yaml:"action"`
	} `yaml:"rule"`
}

type RuleMeta struct {
	Event      string `yaml:"event"`
	IgnoreBots bool   `yaml:"ignoreBots"`
}

type RuleString struct {
	Name  string `yaml:"name"`
	Value string `yaml:"value"`
}

type RuleAction struct {
	Type event.EventType `yaml:"type"`

	// Optional fields depending on the action
	MessageTemplate string `yaml:"messageTemplate,omitempty"`
	Duration        string `yaml:"duration,omitempty"`
	Reason          string `yaml:"reason,omitempty"`
}

type SimplifiedRule struct {
	Name       string
	Event      string
	IgnoreBots bool

	Strings map[string]RuleString
	Program cel.Program

	Action RuleAction
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
		switch yamlRule.Rule.Meta.Event {
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
			return nil, fmt.Errorf("unsupported rule event: %s", yamlRule.Rule.Meta.Event)
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
			Event:      yamlRule.Rule.Meta.Event,
			IgnoreBots: yamlRule.Rule.Meta.IgnoreBots,

			Strings: stringsMap,
			Program: program,
			Action:  yamlRule.Rule.Action,
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
	if r.Event != string(ctx.Type) {
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
