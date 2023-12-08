package lingquan

import (
	_ "embed"
	"errors"
	"github.com/xeipuuv/gojsonschema"
)

//go:embed schema/model.json
var schemaString string

var modelSchema *gojsonschema.Schema

func init() {
	var err error
	modelSchema, err = gojsonschema.NewSchema(gojsonschema.NewStringLoader(schemaString))
	if err != nil {
		panic(err)
	}
}

func ValidateJson(modelJson string) (bool, error) {
	result, err := modelSchema.Validate(gojsonschema.NewStringLoader(modelJson))
	if err != nil {
		return false, err
	}
	if result.Valid() {
		return true, nil
	} else {
		info := ""
		for _, e := range result.Errors() {
			info += e.String() + "\n"
		}

		return false, errors.New(info)
	}
}
