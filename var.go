package config

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"
)

// MARK: Types

type Signed interface {
	~int | ~int64 | ~int32 | ~int16 | ~int8
}
type Unsigned interface {
	~uint | ~uint64 | ~uint32 | ~uint16 | ~uint8
}
type Integer interface{ Signed | Unsigned }
type Float interface{ ~float64 | ~float32 }

type Var interface {
	set(v any) error
	EnvName() string
	FlagName() string
	Value() any
	String() string
	Default() any
	Reset()
	Persist() error
	Description() string
}

type VarCfg[T any] struct {
	EnvName     string
	FlagName    string
	Default     T
	Description string
	ParseString func(string) (T, error)
}

type cfgVar[T any] struct {
	envName        string
	flagName       string
	dVal           T
	value          *T
	customStrParse func(string) (T, error)
	description    string
}

// MARK: Public Functions

func NewVar[T any](cfg VarCfg[T]) (*T, Var) {
	v := &cfgVar[T]{
		envName:     cfg.EnvName,
		flagName:    cfg.FlagName,
		dVal:        cfg.Default,
		description: cfg.Description,
	}
	v.value = new(T)
	*v.value = cfg.Default
	if cfg.ParseString != nil {
		v.customStrParse = cfg.ParseString
	}
	v.registerFlag()
	global.Register(v)
	return v.value, v
}

func (v *cfgVar[T]) EnvName() string {
	return v.envName
}

func (v *cfgVar[T]) FlagName() string {
	return v.flagName
}

func (v *cfgVar[T]) Value() any {
	return v.value
}

func (v *cfgVar[T]) Default() any {
	return v.dVal
}

func (v *cfgVar[T]) Description() string {
	return v.description
}

func (v *cfgVar[T]) String() string {
	switch t := any(v.Value).(type) {
	case nil:
		return "null"
	case fmt.Stringer:
		return t.String()
	case *fmt.Stringer:
		if t == nil {
			return "null"
		}
		return (*t).String()
	case string:
		return t
	case time.Duration:
		return t.String()
	case bool:
		return strconv.FormatBool(t)
	case int:
		return strconv.Itoa(t)
	case int64:
		return strconv.Itoa(int(t))
	case int32:
		return strconv.Itoa(int(t))
	case int16:
		return strconv.Itoa(int(t))
	case int8:
		return strconv.Itoa(int(t))
	case uint:
		return strconv.FormatUint(uint64(t), 10)
	case uint64:
		return strconv.FormatUint(t, 10)
	case uint32:
		return strconv.FormatUint(uint64(t), 10)
	case uint16:
		return strconv.FormatUint(uint64(t), 10)
	case uint8:
		return strconv.FormatUint(uint64(t), 10)
	case float64:
		return strconv.FormatFloat(t, 'g', -1, 64)
	case float32:
		return strconv.FormatFloat(float64(t), 'g', -1, 32)
	default:
		return fmt.Sprintf(`%v`, t)
	}
}

func (v *cfgVar[T]) Reset() {
	if v.value == nil {
		v.value = new(T)
	}
	*v.value = v.dVal
}

func (v *cfgVar[T]) Persist() error {
	if v.value == nil {
		v.dVal = *new(T)
		return os.Unsetenv(v.envName)
	}
	v.dVal = *v.value
	if v.envName != "" {
		return os.Setenv(v.envName, v.String())
	}
	return nil
}

// MARK: Private Functions

func (v *cfgVar[T]) defaultStrParse(s string) (err error) {
	switch v := any(v.value).(type) {
	case *string:
		*v = s
		return
	case *bool:
		*v, err = strconv.ParseBool(s)
		return
	case *time.Duration:
		*v, err = time.ParseDuration(s)
		return
	case *time.Time:
		*v, err = time.Parse(time.RFC3339, s)
		return
	case *int:
		return parseInt(v, s)
	case *int64:
		return parseInt64(v, s)
	case *int32:
		return parseInt32(v, s)
	case *int16:
		return parseInt16(v, s)
	case *int8:
		return parseInt8(v, s)
	case *uint:
		return parseUint(v, s)
	case *uint64:
		return parseUint64(v, s)
	case *uint32:
		return parseUint32(v, s)
	case *uint16:
		return parseUint16(v, s)
	case *uint8:
		return parseUint8(v, s)
	default:
		err = fmt.Errorf("can't parse string %[1]q to Var of type %[2]T", s, v)
		return
	}
}

func (v *cfgVar[T]) set(a any) (err error) {
	switch a := a.(type) {
	case string:
		if v.customStrParse != nil {
			*v.value, err = v.customStrParse(a)
			return err
		}
		return v.defaultStrParse(a)
	case T:
		*v.value = a
		return nil
	default:
		return fmt.Errorf("can't assign (%[1]T)%[1]v to Var of type %[2]T", a, v.dVal)
	}
}

func (v *cfgVar[T]) registerFlag() {
	if v == nil || v.flagName == "" {
		return
	}
	if v.customStrParse != nil {
		flag.Func(v.flagName, v.description, func(s string) (err error) {
			*v.value, err = v.customStrParse(s)
			return err
		})
	}

	switch t := any(v.value).(type) {
	case *bool:
		flag.BoolVar(t, v.flagName, any(v.dVal).(bool), v.description)
	case *string:
		flag.StringVar(t, v.flagName, any(v.dVal).(string), v.description)
	case *int:
		flag.IntVar(t, v.flagName, any(v.dVal).(int), v.description)
	case *uint:
		flag.UintVar(t, v.flagName, any(v.dVal).(uint), v.description)
	case *float64:
		flag.Float64Var(t, v.flagName, any(v.dVal).(float64), v.description)
	case *time.Duration:
		flag.DurationVar(t, v.flagName, any(v.dVal).(time.Duration), v.description)
	default:
		panic(fmt.Errorf("must provide a custom parse function for type %T", t))
	}
}
