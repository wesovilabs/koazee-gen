package operation

import (
	"fmt"
	"log"
	"os"
	"strings"
	"text/template"
	"time"
)

const containsDispatcherTemplateText = `
package contains

import "reflect"


type dispatchFunction func(reflect.Value, interface{}) bool

var dispatcher = map[string]dispatchFunction{
	{{- range .Data }}
	"{{.Type}}":  {{.FuncContainsVal}},
	"*{{.Type}}":  {{.FuncContainsPtr}},
	{{- end }}
}

func dispatch(items reflect.Value, val interface{},itemsType reflect.Type) (bool,bool) {
	input:=itemsType.String()
	if fnVal,ok:=dispatcher[input];ok{
		return true, fnVal(items,val)
    }
	return false, false
}

{{- range .Data }}

func {{.FuncContainsVal}}(itemsValue reflect.Value, val interface{}) bool  {
	element := val.({{.Type}})
	array := itemsValue.Interface().([]{{.Type}})
	for _, val := range array {
		if val == element {
			return true
		}
	}
	return false
}

func {{.FuncContainsPtr}}(itemsValue reflect.Value, val interface{}) bool  {
	element := val.(*{{.Type}})
	array := itemsValue.Interface().([]*{{.Type}})
	for _, val := range array {
		if val == element {
			return true
		}
	}
	return false
}

{{- end }}
`

type Contains struct {

}

func (w *TypeInput) FuncContainsVal() string {
	return fmt.Sprintf("contains%s", strings.Title(w.Type))
}
func (w *TypeInput) FuncContainsPtr() string {
	return fmt.Sprintf("containsPtr%s", strings.Title(w.Type))
}



func GenerateContainsDispatcher(outputPath string) {
	f, err := os.Create(fmt.Sprintf("%s/dispatcher.go", outputPath))
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	fmt.Println(fmt.Sprintf("%s/dispatcher.go", outputPath))
	template.
		Must(template.New("").
			Parse(containsDispatcherTemplateText)).
		Execute(f, struct {
			Timestamp time.Time
			Data      []*TypeInput
		}{
			Timestamp: time.Now(),
			Data:      data,
		})

}