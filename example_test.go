package config

import (
	"fmt"
	"os"
	"time"
)

func ExampleNewVar() {
	type Custom uint8
	const (
		customOne = iota + 1
		customTwo
		customThree
	)
	var parseCustom func(string) (Custom, error) = func(s string) (Custom, error) {
		switch s {
		case "one":
			return customOne, nil
		case "two":
			return customTwo, nil
		case "three":
			return customThree, nil
		case "undefined":
			return 0, nil
		default:
			return 0, fmt.Errorf("invalid Custom %s", s)
		}
	}

	_ = os.Setenv("ENV_STR", "hello, world")
	_ = os.Setenv("ENV_BOOL", "true")
	_ = os.Setenv("ENV_INT", "-1")
	_ = os.Setenv("ENV_UINT8", "123")
	_ = os.Setenv("ENV_TIME", "1999-12-31T23:59:59Z")
	_ = os.Setenv("ENV_DURATION", "8h47m")
	_ = os.Setenv("ENV_CUSTOM", "two")

	envStr, _ := NewVar(VarCfg[string]{EnvName: "ENV_STR"})
	envBool, _ := NewVar(VarCfg[bool]{EnvName: "ENV_BOOL"})
	envInt, _ := NewVar(VarCfg[int]{EnvName: "ENV_INT"})
	envUint8, _ := NewVar(VarCfg[int]{EnvName: "ENV_UINT8"})
	envTime, _ := NewVar(VarCfg[time.Time]{EnvName: "ENV_TIME"})
	envDuration, _ := NewVar(VarCfg[time.Duration]{EnvName: "ENV_DURATION"})
	envCustom, _ := NewVar(VarCfg[Custom]{
		EnvName:     "ENV_CUSTOM",
		Default:     customOne,
		ParseString: parseCustom,
	})
	_ = ReadAll()

	fmt.Printf("ENV_STR: %q\n", *envStr)
	fmt.Printf("ENV_BOOL: %t\n", *envBool)
	fmt.Printf("ENV_INT: %d\n", *envInt)
	fmt.Printf("ENV_UINT8: %d\n", *envUint8)
	fmt.Printf("ENV_TIME: %s\n", envTime.Format(time.DateTime))
	fmt.Printf("ENV_DURATION: %s\n", envDuration.String())
	fmt.Printf("ENV_CUSTOM: %d\n", *envCustom)

	// Output:
	// ENV_STR: "hello, world"
	// ENV_BOOL: true
	// ENV_INT: -1
	// ENV_UINT8: 123
	// ENV_TIME: 1999-12-31 23:59:59
	// ENV_DURATION: 8h47m0s
	// ENV_CUSTOM: 2
}
