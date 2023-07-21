package log

import (
	"errors"
	"fmt"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"reflect"
	"strings"
)

func checkConfigValid(config logFileConfig)(error){
	validate := validator.New()
	err := validate.Struct(config)
	return buildError(err)
}

func buildError(err error) error{
	if err == nil{
		return nil
	}

	keys := filterErrKeys(err)
	msg := filterErrMsg(keys)
	return errors.New(msg)
}

func filterErrKeys(err error)map[string]int{
	e := en.New()
	uni := ut.New(e)
	trans, _ := uni.GetTranslator("en")

	keys := map[string]int{}
	errMap := err.(validator.ValidationErrors).Translate(trans)
	for field, _ := range errMap{
		key := field[strings.Index(field, ".")+1:]
		keys[key] = 0
	}

	return keys
}

func filterErrMsg(keys map[string]int) string{
	msg := ""

	vt := reflect.TypeOf(logFileConfig{})
	for i:= 0; i < vt.NumField(); i++{
		name := vt.Field(i).Name
		if _,ok := keys[name];!ok{
			continue
		}
		tagContent := vt.Field(i).Tag.Get("err")
		msg += fmt.Sprintf("%s:%s\n",name,tagContent)
	}

	return msg
}