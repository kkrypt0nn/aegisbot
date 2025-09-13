package rules

import (
	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
	"github.com/kkrypt0nn/aegisbot/proto"
)

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
