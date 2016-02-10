package qcli

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

// Flag retrieves a value from a given flag name. If it is unable to dereference a
// value the function will return nil
func Flag(name string) interface{} {
	val := flagvars.m["_"][name]
	switch val.(type) {
	case *bool:
		return *(val.(*bool))
		break
	case *float64:
		return *(val.(*float64))
		break
	case *string:
		return *(val.(*string))
		break
	}
	return nil
}

//FlagSet is just an empty accessor
type FlagSet struct {
	m map[string]interface{}
}

// Flagset retrieves a FlagSet given the name as provided in the required JSON file
// Why the weird difference in case? Because Godocs complains at public members of private structures!
// Thus bloating up your IDEs and godocs with useless stuff. Yes this is a complaint
func Flagset(name string) FlagSet {
	return FlagSet{flagvars.m[name]}
}

// Flag retrieves a flag value that is associated with a flag name on a parent FlagSet
func (f FlagSet) Flag(name string) interface{} {
	val := f.m[name]
	switch val.(type) {
	case *bool:
		return *(val.(*bool))
		break
	case *float64:
		return *(val.(*float64))
		break
	case *string:
		return *(val.(*string))
		break
	}
	return nil
}

var flagvars struct {
	m map[string]map[string]interface{}
}

func getArgIndex(name string) int {
	for i, a := range os.Args {
		if a == name {
			return i
		}
	}
	return -1
}

func initFlagSets(_flagsets interface{}) {
	flagsets := _flagsets.([]interface{})
	for i, _flagset := range flagsets {
		flagset := _flagset.(map[string]interface{})

		errorHandling := flagset["errorHandling"].(float64)
		name := flagset["name"].(string)

		if errorHandling < 0 || errorHandling > 2 {
			log.Fatalf("Unable to parse flagset at index \"%d\"\n", i)
		}
		if len(name) == 0 {
			log.Fatalf("Unable to parse flagset at index %d\n", i)
		}

		newFlagSet := flag.NewFlagSet(name, flag.ErrorHandling(errorHandling))
		flagvars.m[name] = make(map[string]interface{})
		err := initFlags(flagset["flags"], name, newFlagSet)
		if err != nil {
			log.Fatalf("Unable to parse flagset at index %d: %s\n", i, err)
		}
		if index := getArgIndex(name); index > -1 {
			newFlagSet.Parse(os.Args[index:])
		}
	}
}

func initFlags(_flags interface{}, flagsetName string, flagset ...*flag.FlagSet) error {
	var err error
	var isFlagset bool
	if len(flagset) > 0 {
		isFlagset = true
	}
	flags := _flags.([]interface{})
	for i, _flag := range flags {
		flg := _flag.(map[string]interface{})
		name := flg["name"].(string)
		_type, _ := flg["type"].(string)
		usage, _ := flg["usage"].(string)
		def := flg["default"]
		if len(name) == 0 {
			return fmt.Errorf("Unable to parse flag at index %d: No name specified!", i)
		}

		if len(_type) == 0 {
			//do type switch
			switch def.(type) {
			case bool:
				var b bool
				if isFlagset {
					flagvars.m[flagsetName][name] = &b
					flagset[0].BoolVar(&b, name, def.(bool), usage)
				} else {
					flag.BoolVar(&b, name, def.(bool), usage)
				}
				break
			case float64:
				var f float64
				if isFlagset {
					flagvars.m[flagsetName][name] = &f
					flagset[0].Float64Var(&f, name, def.(float64), usage)
				} else {
					flag.Float64Var(&f, name, def.(float64), usage)
				}
				break
			case string:
				var s string
				if isFlagset {
					flagvars.m[flagsetName][name] = &s
					flagset[0].StringVar(&s, name, def.(string), usage)
				} else {
					flag.StringVar(&s, name, def.(string), usage)
				}
				break
			default:
				return fmt.Errorf("Unable to parse flag at index %d: Unable to determine flag type!", i)
			}
		} else {
			//do value switch
			switch _type {
			case "bool":
				if def == nil {
					def = false
				}
				var b bool
				if isFlagset {
					flagvars.m[flagsetName][name] = &b
					flagset[0].BoolVar(&b, name, def.(bool), usage)
				} else {
					flag.BoolVar(&b, name, def.(bool), usage)
				}
				break
			case "float64":
				if def == nil {
					def = 0
				}
				var f float64
				if isFlagset {
					flagvars.m[flagsetName][name] = &f
					flagset[0].Float64Var(&f, name, def.(float64), usage)
				} else {
					flag.Float64Var(&f, name, def.(float64), usage)
				}
				break
			case "string":
				if def == nil {
					def = ""
				}
				var s string
				if isFlagset {
					flagvars.m[flagsetName][name] = &s
					flagset[0].StringVar(&s, name, def.(string), usage)
				} else {
					flag.StringVar(&s, name, def.(string), usage)
				}
				break
			default:
				return fmt.Errorf("Unable to parse flag at index %d: Unable to determine flag type!", i)
			}
		}
	}

	return err
}

func init() {
	flagvars.m = make(map[string]map[string]interface{})
	jsonFile, err := os.Open("./flags.json")
	if err != nil {
		return
	}
	jsonData, _ := ioutil.ReadAll(jsonFile)
	var flagMap map[string]interface{}
	err = json.Unmarshal(jsonData, &flagMap)
	if err != nil {
		log.Fatalln(err)
	}

	//handle top-level flags
	flags := flagMap["flags"]
	if flags != nil {
		initFlags(flagMap["flags"], "")
	}

	//handle flagsets
	initFlagSets(flagMap["flagsets"])

	flag.Parse()
}
