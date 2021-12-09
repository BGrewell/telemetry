package main

import (
	"flag"
)

type CommandLineFlag struct {
	Name     string
	Value    interface{}
	Usage    string
	AltNames *[]string
}

func AddFlagBool(name string, section interface{}, defaultValue bool, usage string, altNames *[]string) *bool {
	var v bool

	flag.BoolVar(&v, name, defaultValue, usage)
	if altNames != nil {
		for _, altName := range *altNames {
			flag.BoolVar(&v, altName, defaultValue, usage)
		}
	}

	AddToFlagMap(name, section, defaultValue, usage, altNames)

	return &v
}

func AddFlagString(name string, section interface{}, defaultValue string, usage string, altNames *[]string) *string {
	var v string

	flag.StringVar(&v, name, defaultValue, usage)
	if altNames != nil {
		for _, altName := range *altNames {
			flag.StringVar(&v, altName, defaultValue, usage)
		}
	}

	AddToFlagMap(name, section, defaultValue, usage, altNames)

	return &v
}

func AddFlagInt(name string, section interface{}, defaultValue int, usage string, altNames *[]string) *int {
	var v int

	flag.IntVar(&v, name, defaultValue, usage)
	if altNames != nil {
		for _, altName := range *altNames {
			flag.IntVar(&v, altName, defaultValue, usage)
		}
	}

	AddToFlagMap(name, section, defaultValue, usage, altNames)

	return &v
}

func AddFlagInt64(name string, section interface{}, defaultValue int64, usage string, altNames *[]string) *int64 {
	var v int64

	flag.Int64Var(&v, name, defaultValue, usage)
	if altNames != nil {
		for _, altName := range *altNames {
			flag.Int64Var(&v, altName, defaultValue, usage)
		}
	}

	AddToFlagMap(name, section, defaultValue, usage, altNames)

	return &v
}

func AddFlagUnint(name string, section interface{}, defaultValue uint, usage string, altNames *[]string) *uint {
	var v uint

	flag.UintVar(&v, name, defaultValue, usage)
	if altNames != nil {
		for _, altName := range *altNames {
			flag.UintVar(&v, altName, defaultValue, usage)
		}
	}

	AddToFlagMap(name, section, defaultValue, usage, altNames)

	return &v
}

func AddFlagUint64(name string, section interface{}, defaultValue uint64, usage string, altNames *[]string) *uint64 {
	var v uint64

	flag.Uint64Var(&v, name, defaultValue, usage)
	if altNames != nil {
		for _, altName := range *altNames {
			flag.Uint64Var(&v, altName, defaultValue, usage)
		}
	}

	AddToFlagMap(name, section, defaultValue, usage, altNames)

	return &v
}

func AddFlagFloat64(name string, section interface{}, defaultValue float64, usage string, altNames *[]string) *float64 {
	var v float64

	flag.Float64Var(&v, name, defaultValue, usage)
	if altNames != nil {
		for _, altName := range *altNames {
			flag.Float64Var(&v, altName, defaultValue, usage)
		}
	}

	AddToFlagMap(name, section, defaultValue, usage, altNames)

	return &v
}

func AddToFlagMap(name string, section interface{}, defaultValue interface{}, usage string, altNames *[]string) {

	var sections []string
	switch section.(type) {
	case string:
		sections = []string{section.(string)}
	case []string:
		sections = section.([]string)
	default:
		panic("Invalid section type")
	}

	for _, section := range sections {
		if _, ok := flagMap[section]; !ok {
			flagMap[section] = make([]*CommandLineFlag, 0)
			flagKeys = append(flagKeys, section)
		}
		clf := CommandLineFlag{
			Name:     name,
			Value:    defaultValue,
			Usage:    usage,
			AltNames: altNames,
		}
		flagMap[section] = append(flagMap[section], &clf)
	}

}