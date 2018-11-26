package operation

import (
	"fmt"
	"log"
	"os"
	"strings"
	"text/template"
	"time"
)

const filterDispatcherTemplateText = `
package filter

import "reflect"


type dispatchFunction func(items reflect.Value, fn interface{}) interface{}

var dispatcher = map[string]dispatchFunction{
	{{- range .Data }}
	"{{.Type}}":  {{.FuncFilterVal}},
	"*{{.Type}}":  {{.FuncFilterPtr}},
	{{- end }}
}

func dispatch(items reflect.Value, function interface{}, info *filterInfo) (bool,interface{}) {
	input:=info.fnInputType.String()
	if fnVal,ok:=dispatcher[input];ok{
		return true, fnVal(items,function)
    }
	return false, nil
}

{{- range .Data }}

func {{.FuncFilterVal}}(itemsValue reflect.Value, function interface{}) interface{}  {
	input := itemsValue.Interface().([]{{.Type}})
	output := make([]{{.Type}},0)
	fn := function.(func({{.Type}}) bool)
	for i := 0; i < len(input); i++ {
		if fn(input[i]) {
			output = append(output, input[i]) 
		}
	}
	return output
}

func {{.FuncFilterPtr}}(itemsValue reflect.Value, function interface{}) interface{}  {
	input := itemsValue.Interface().([]*{{.Type}})
	output := make([]*{{.Type}},0)
	fn := function.(func(*{{.Type}}) bool)
	for i := 0; i < len(input); i++ {
		if fn(input[i]) {
			output = append(output, input[i])
		}
	}
	return output
}

{{- end }}
`
func (w *TypeOutput) FuncFilterVal() string {
	return fmt.Sprintf("filter%s", strings.Title(w.Type))
}
func (w *TypeOutput) FuncFilterPtr() string {
	return fmt.Sprintf("filterPtr%s", strings.Title(w.Type))
}



func GenerateFilterDispatcher(outputPath string) {
	f, err := os.Create(fmt.Sprintf("%s/dispatcher.go", outputPath))
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	template.
		Must(template.New("").
			Parse(filterDispatcherTemplateText)).
		Execute(f, struct {
			Timestamp time.Time
			Data      []*TypeInput
		}{
			Timestamp: time.Now(),
			Data:      data,
		})

}