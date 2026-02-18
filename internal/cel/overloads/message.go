package overloads

import (
	"regexp"

	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
	"github.com/kkrypt0nn/aegisbot/proto"
)

var urlRegex = regexp.MustCompile(`https?://[^\s]+`)

func MessageMethods() []cel.EnvOption {
	return []cel.EnvOption{
		cel.Function("isDM",
			cel.MemberOverload("message_isDM_bool",
				[]*cel.Type{cel.ObjectType("proto.Message")},
				cel.BoolType,
				cel.UnaryBinding(func(m ref.Val) ref.Val {
					msg, ok := m.Value().(*proto.Message)
					if !ok {
						return types.NewErr("unexpected type '%v', expected '*proto.Message'", m.Type())
					}
					return types.Bool((msg.Channel.GetType() == proto.ChannelType_CHANNEL_TYPE_DM) || (msg.Channel.GetType() == proto.ChannelType_CHANNEL_TYPE_GROUP_DM))
				}),
			),
		),

		cel.Function("getMentions",
			cel.MemberOverload("message_getMentions_list",
				[]*cel.Type{cel.ObjectType("proto.Message")},
				cel.ListType(cel.StringType),
				cel.UnaryBinding(func(m ref.Val) ref.Val {
					msg, ok := m.Value().(*proto.Message)
					if !ok {
						return types.NewErr("unexpected type '%v', expected '*proto.Message'", m.Type())
					}
					if msg.Mentions == nil {
						return types.NewDynamicList(types.DefaultTypeAdapter, []any{})
					}
					res := make([]any, len(msg.Mentions))
					for i, member := range msg.Mentions {
						res[i] = member
					}
					return types.NewDynamicList(types.DefaultTypeAdapter, res)
				}),
			),
		),

		cel.Function("hasLinks",
			cel.MemberOverload("message_hasLinks_bool",
				[]*cel.Type{cel.ObjectType("proto.Message")},
				cel.BoolType,
				cel.UnaryBinding(func(m ref.Val) ref.Val {
					msg, ok := m.Value().(*proto.Message)
					if !ok {
						return types.NewErr("unexpected type '%v', expected '*proto.Message'", m.Type())
					}
					return types.Bool(urlRegex.MatchString(msg.Content))
				}),
			),
		),

		cel.Function("getLinks",
			cel.MemberOverload("message_getLinks_list",
				[]*cel.Type{cel.ObjectType("proto.Message")},
				cel.ListType(cel.StringType),
				cel.UnaryBinding(func(m ref.Val) ref.Val {
					msg, ok := m.Value().(*proto.Message)
					if !ok {
						return types.NewErr("unexpected type '%v', expected '*proto.Message'", m.Type())
					}
					if msg.Content == "" {
						return types.NewDynamicList(types.DefaultTypeAdapter, []any{})
					}
					matches := urlRegex.FindAllString(msg.Content, -1)
					values := make([]any, len(matches))
					for i, s := range matches {
						values[i] = types.String(s)
					}
					return types.NewDynamicList(types.DefaultTypeAdapter, values)
				}),
			),
		),
	}
}
