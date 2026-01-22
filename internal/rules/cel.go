package rules

import (
	"regexp"

	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
	"github.com/kkrypt0nn/aegisbot/proto"
)

var urlRegex = regexp.MustCompile(`https?://[^\s]+`)

var isBotOverload = cel.MemberOverload("member_isBot_bool",
	[]*cel.Type{cel.ObjectType("proto.Member")},
	cel.BoolType,
	cel.UnaryBinding(func(m ref.Val) ref.Val {
		member, ok := m.Value().(*proto.Member)
		if !ok {
			return types.NewErr("unexpected type '%v', expected '*proto.Member'", member.ProtoReflect().Type())
		}
		return types.Bool(member.GetBot())
	}),
)

var messageHasLinks = cel.MemberOverload("message_hasLinks_bool",
	[]*cel.Type{cel.ObjectType("proto.Message")},
	cel.BoolType,
	cel.UnaryBinding(func(m ref.Val) ref.Val {
		message, ok := m.Value().(*proto.Message)
		if !ok {
			return types.NewErr("unexpected type '%v', expected '*proto.Mesage'", message.ProtoReflect().Type())
		}

		return types.Bool(urlRegex.MatchString(message.Content))
	}),
)

var messageGetLinks = cel.MemberOverload("message_getLinks_list",
	[]*cel.Type{cel.ObjectType("proto.Message")},
	cel.ListType(cel.StringType),
	cel.UnaryBinding(func(m ref.Val) ref.Val {
		message, ok := m.Value().(*proto.Message)
		if !ok {
			return types.NewErr("unexpected type '%v', expected '*proto.Mesage'", message.ProtoReflect().Type())
		}

		if message.Content == "" {
			return types.NewDynamicList(types.DefaultTypeAdapter, []any{})
		}

		matches := urlRegex.FindAllString(message.Content, -1)
		values := make([]any, len(matches))
		for i, s := range matches {
			values[i] = types.String(s)
		}
		return types.NewDynamicList(types.DefaultTypeAdapter, values)
	}),
)
