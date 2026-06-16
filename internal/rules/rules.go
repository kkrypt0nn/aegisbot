package rules

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/common/types"
	"github.com/kkrypt0nn/aegisbot/internal/actions"
	libtime "github.com/kkrypt0nn/aegisbot/internal/cel/libs/time"
	"github.com/kkrypt0nn/aegisbot/internal/cel/overloads"
	"github.com/kkrypt0nn/aegisbot/internal/event"
	"github.com/kkrypt0nn/aegisbot/proto"
	"sigs.k8s.io/yaml"
)

type YAMLRule struct {
	Rule struct {
		Name       string       `yaml:"name"`
		Disabled   bool         `yaml:"disabled"`
		Meta       RuleMeta     `yaml:"meta"`
		Strings    []RuleString `yaml:"strings"`
		Expression string       `yaml:"expression"`
		Actions    []RuleAction `yaml:"actions"`
	} `yaml:"rule"`
}

type RuleMeta struct {
	Event      event.EventType `yaml:"event"`
	IgnoreBots bool            `yaml:"ignoreBots"`
}

type RuleString struct {
	Name  string `yaml:"name"`
	Value string `yaml:"value"`
}

type RuleAction struct {
	Type actions.ActionType `yaml:"type"`

	// Optional fields depending on the action
	ChannelID string `yaml:"channelId,omitempty"`
	Message   string `yaml:"message,omitempty"`
	Duration  string `yaml:"duration,omitempty"`
	Reason    string `yaml:"reason,omitempty"`
}

type SimplifiedRule struct {
	Name       string
	Disabled   bool
	Event      event.EventType
	IgnoreBots bool

	Strings map[string]RuleString
	Program cel.Program

	Actions []RuleAction

	RawExpression string
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
		case event.EventMessageCreate, event.EventMessageUpdate:
			celVars = append(celVars,
				cel.Types(&proto.Member{}, &proto.Message{}),
				cel.Variable("author", cel.ObjectType("proto.Member")),
				cel.Variable("message", cel.ObjectType("proto.Message")),
			)
		case event.EventMemberJoin, event.EventMemberUpdate:
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
		celVars = append(celVars, overloads.MemberMethods()...)
		celVars = append(celVars, overloads.MessageMethods()...)

		// Add the libraries
		celVars = append(celVars, libtime.Lib())

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
			Disabled:   yamlRule.Rule.Disabled,
			Event:      yamlRule.Rule.Meta.Event,
			IgnoreBots: yamlRule.Rule.Meta.IgnoreBots,

			Strings: stringsMap,
			Program: program,
			Actions: yamlRule.Rule.Actions,

			RawExpression: yamlRule.Rule.Expression,
		}
		simplifiedRules = append(simplifiedRules, rule)
	}

	return simplifiedRules, nil
}

func Load(dir string) ([]*SimplifiedRule, map[string]*SimplifiedRule, error) {
	var allRules []*SimplifiedRule
	rulesByName := make(map[string]*SimplifiedRule)

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

		for _, r := range rules {
			if _, exists := rulesByName[r.Name]; exists {
				return fmt.Errorf("duplicate rule name: %s", r.Name)
			}
			rulesByName[r.Name] = r
		}

		allRules = append(allRules, rules...)
		return nil
	})

	if err != nil {
		return nil, nil, err
	}
	return allRules, rulesByName, nil
}

func (r *SimplifiedRule) Evaluate(ctx *event.Context) (bool, error) {
	if r.Disabled {
		return false, nil
	}
	if r.Event != ctx.Type {
		return false, nil
	}
	if r.IgnoreBots && ctx.Bot {
		return false, nil
	}

	input := map[string]any{}
	switch ctx.Type {
	case event.EventMessageCreate, event.EventMessageUpdate:
		input["author"] = ctx.Message.Author
		input["message"] = ctx.Message
	case event.EventMemberJoin, event.EventMemberUpdate:
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
