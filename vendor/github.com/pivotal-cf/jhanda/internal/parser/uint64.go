package parser

import (
	"flag"
	"fmt"
	"os"
	"reflect"
	"strconv"
)

func NewUint64(set *flag.FlagSet, field reflect.Value, tags reflect.StructTag) (*Flag, error) {
	var defaultValue uint64
	defaultStr, ok := tags.Lookup("default")
	if ok {
		var err error
		defaultValue, err = strconv.ParseUint(defaultStr, 0, 64)
		if err != nil {
			return &Flag{}, fmt.Errorf("could not parse uint64 default value %q: %s", defaultStr, err)
		}
	}

	var f Flag
	short, ok := tags.Lookup("short")
	if ok {
		set.Uint64Var(field.Addr().Interface().(*uint64), short, defaultValue, "")
		f.flags = append(f.flags, set.Lookup(short))
		f.name = fmt.Sprintf("-%s", short)
	}

	long, ok := tags.Lookup("long")
	if ok {
		set.Uint64Var(field.Addr().Interface().(*uint64), long, defaultValue, "")
		f.flags = append(f.flags, set.Lookup(long))
		f.name = fmt.Sprintf("--%s", long)
	}

	env, ok := tags.Lookup("env")
	if ok {
		envStr := os.Getenv(env)
		if envStr != "" {
			envValue, err := strconv.ParseUint(envStr, 0, 64)
			if err != nil {
				return &Flag{}, fmt.Errorf("could not parse uint64 environment variable %s value %q: %s", env, envStr, err)
			}

			field.SetUint(envValue)
			f.set = true
		}
	}

	_, f.required = tags.Lookup("required")

	return &f, nil
}
