package libtime

import (
	"time"

	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
)

func Lib() cel.EnvOption {
	return cel.Lib(&libTime{})
}

type libTime struct{}

func (*libTime) LibraryName() string {
	return "core.time"
}

func (*libTime) CompileOptions() []cel.EnvOption {
	return []cel.EnvOption{
		func(env *cel.Env) (*cel.Env, error) {
			nowFunc := cel.Function(
				"time.now",
				cel.Overload(
					"time_now_timestamp",
					[]*cel.Type{},
					types.TimestampType,
					cel.FunctionBinding(func(args ...ref.Val) ref.Val {
						return env.CELTypeAdapter().NativeToValue(types.Timestamp{Time: time.Now()})
					}),
				),
			)
			sinceFunc := cel.Function(
				"time.since",
				cel.Overload(
					"time_since_duration",
					[]*cel.Type{types.TimestampType},
					types.DurationType,
					cel.UnaryBinding(func(arg ref.Val) ref.Val {
						ts, ok := arg.(types.Timestamp)
						if !ok {
							return types.NewErr("time.since: expected timestamp")
						}
						return types.Duration{Duration: time.Since(ts.Time)}
					}),
				),
			)
			return env.Extend(nowFunc, sinceFunc)
		},
	}
}

func (*libTime) ProgramOptions() []cel.ProgramOption {
	return nil
}
