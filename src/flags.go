package qcli

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

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
		err := initFlags(flagset["flags"], newFlagSet)
		if err != nil {
			log.Fatalf("Unable to parse flagset at index %d: %s\n", i, err)
		}
		newFlagSet.Parse([]string{"force"})
	}
}

func initFlags(_flags interface{}, flagset ...*flag.FlagSet) error {
	var err error = nil
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
			return errors.New(fmt.Sprintf("Unable to parse flag at index %d: No name specified!", i))
		}

		if len(_type) == 0 {
			//do type switch
			switch def.(type) {
			case bool:
				if isFlagset {
					flagset[0].Bool(name, def.(bool), usage)
				} else {
					flag.Bool(name, def.(bool), usage)
				}
				break
			case float64:
				if isFlagset {
					flagset[0].Float64(name, def.(float64), usage)
				} else {
					flag.Float64(name, def.(float64), usage)
				}
				break
			case string:
				if isFlagset {
					flagset[0].String(name, def.(string), usage)
				} else {
					flag.String(name, def.(string), usage)
				}
				break
			default:
				return errors.New(fmt.Sprintf("Unable to parse flag at index %d: Unable to determine flag type!", i))
			}
		} else {
			//do value switch
			switch _type {
			case "bool":
				if isFlagset {
					flagset[0].Bool(name, def.(bool), usage)
				} else {
					flag.Bool(name, def.(bool), usage)
				}
				break
			case "float64":
				if isFlagset {
					flagset[0].Float64(name, def.(float64), usage)
				} else {
					flag.Float64(name, def.(float64), usage)
				}
				break
			case "string":
				if isFlagset {
					flagset[0].String(name, def.(string), usage)
				} else {
					flag.String(name, def.(string), usage)
				}
				break
			default:
				return errors.New(fmt.Sprintf("Unable to parse flag at index %d: Unable to determine flag type!", i))
			}
		}
	}

	return err
}

func init() {
	jsonFile, _ := os.Open("./flags.json")
	jsonData, _ := ioutil.ReadAll(jsonFile)
	var flagMap map[string]interface{}
	err := json.Unmarshal(jsonData, &flagMap)
	if err != nil {
		log.Fatalln(err)
	}

	//handle top-level flags
	initFlags(flagMap["flags"])

	//handle flagsets
	initFlagSets(flagMap["flagsets"])

	flag.Parse()
}
