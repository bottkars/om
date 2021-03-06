package parser

import (
	"flag"
	"fmt"
	"os"
	"reflect"
	"strconv"
)

func NewBool(set *flag.FlagSet, field reflect.Value, tags reflect.StructTag) (*Flag, error) {
	var defaultValue bool
	defaultStr, ok := tags.Lookup("default")
	if ok {
		var err error
		defaultValue, err = strconv.ParseBool(defaultStr)
		if err != nil {
			return &Flag{}, fmt.Errorf("could not parse bool default value %q: %s", defaultStr, err)
		}
	}

	var f Flag
	short, ok := tags.Lookup("short")
	if ok {
		set.BoolVar(field.Addr().Interface().(*bool), short, defaultValue, "")
		f.flags = append(f.flags, set.Lookup(short))
		f.name = fmt.Sprintf("-%s", short)
	}

	long, ok := tags.Lookup("long")
	if ok {
		set.BoolVar(field.Addr().Interface().(*bool), long, defaultValue, "")
		f.flags = append(f.flags, set.Lookup(long))
		f.name = fmt.Sprintf("--%s", long)
	}

	env, ok := tags.Lookup("env")
	if ok {
		envStr := os.Getenv(env)
		if envStr != "" {
			envValue, err := strconv.ParseBool(envStr)
			if err != nil {
				return &Flag{}, fmt.Errorf("could not parse bool environment variable %s value %q: %s", env, envStr, err)
			}

			field.SetBool(envValue)
			f.set = true
		}
	}

	_, f.required = tags.Lookup("required")

	return &f, nil
}
