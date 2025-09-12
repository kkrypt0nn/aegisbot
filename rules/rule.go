package rules

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/disgoorg/disgo/events"
	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/common/types"
	"github.com/kkrypt0nn/aegisbot/actions"
	"github.com/kkrypt0nn/aegisbot/log"
	"github.com/kkrypt0nn/aegisbot/proto"
	"sigs.k8s.io/yaml"
)

type YAMLRule struct {
	Rule struct {
		Name       string       `yaml:"name"`
		Meta       RuleMeta     `yaml:"meta"`
		Strings    []RuleString `yaml:"strings"`
		Expression string       `yaml:"expression"`
	} `yaml:"rule"`
}

type RuleMeta struct {
	Action  string `yaml:"action"`
	Context string `yaml:"context"`
}

type RuleString struct {
	Name  string `yaml:"name"`
	Value string `yaml:"value"`
}

type SimplifiedRule struct {
	Name       string
	Action     string
	Context    string
	Strings    map[string]RuleString
	Expression string
	Program    cel.Program
}

type Context struct {
	Message *proto.Message
	Member  *proto.Member
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
			stringsMap[s.Name] = RuleString{
				Name:  s.Name,
				Value: s.Value,
			}
		}

		var celVars []cel.EnvOption
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

		for _, s := range yamlRule.Rule.Strings {
			celVars = append(celVars, cel.Constant(s.Name, cel.StringType, types.String(s.Value)))
		}

		env, err := cel.NewEnv(celVars...)
		if err != nil {
			return nil, err
		}
		condition := strings.TrimSpace(yamlRule.Rule.Expression)
		ast, iss := env.Compile(condition)
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
			Strings:    stringsMap,
			Expression: condition,
			Program:    program,
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

func (r *SimplifiedRule) Evaluate(ctx *Context) (bool, error) {
	input := map[string]any{}
	switch r.Context {
	case "message":
		if ctx.Message == nil {
			return false, fmt.Errorf("no message in context")
		}
		input["message"] = ctx.Message
	case "member":
		if ctx.Member == nil {
			return false, fmt.Errorf("no member in context")
		}
		input["member"] = ctx.Member
	default:
		return false, fmt.Errorf("unknown rule context: %s", r.Context)
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

func (r *SimplifiedRule) EvaluateMessage(e *events.MessageCreate) bool {
	if r.Context != "message" {
		return false
	}
	if e.Message.Author.Bot {
		return false
	}

	ok, err := r.Evaluate(&Context{
		Message: &proto.Message{
			Content: e.Message.Content,
		},
	})
	if err != nil || !ok {
		return false
	}

	log.Info(fmt.Sprintf("Rule %s matched for user %s", r.Name, e.Message.Author.Username))

	actions.Execute(r.Action, e.Client().Rest(), &actions.Context{
		GuildID:   e.GuildID.String(),
		ChannelID: e.ChannelID.String(),
		UserID:    e.Message.Author.ID.String(),
		MessageID: e.MessageID.String(),
	})

	return true
}

func (r *SimplifiedRule) EvaluateMember(e *events.GuildMemberUpdate) bool {
	if r.Context != "member" {
		return false
	}

	ok, err := r.Evaluate(&Context{
		Member: &proto.Member{
			Name: *e.Member.Nick,
		},
	})
	if err != nil || !ok {
		return false
	}

	log.Info(fmt.Sprintf("Rule %s matched for user %s", r.Name, e.Member.User.Username))

	actions.Execute(r.Action, e.Client().Rest(), &actions.Context{
		GuildID:   e.GuildID.String(),
		ChannelID: "1397277714192531707",
		UserID:    e.Member.User.ID.String(),
		MessageID: "",
	})

	return true
}
