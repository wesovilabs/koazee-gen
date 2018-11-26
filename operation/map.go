package operation

import (
	"fmt"
	"log"
	"os"
	"strings"
	"text/template"
	"time"
)

const mapDispatcherTemplateText = `
package maps

import "reflect"


type dispatchFunction func(items reflect.Value, fn interface{}) interface{}

var dispatcher = map[string]map[string]dispatchFunction{
	{{- range .Data }}
	"{{.Type}}": {
		{{- range .Output }}
		"{{.Type}}": {{.FuncMapValToVal}},
		{{- end }}
		{{- range .Output }}
		"*{{.Type}}": {{.FuncMapValToPtr}},
		{{- end }}
	},
	"*{{.Type}}": {
		{{- range .Output }}
		"{{.Type}}": {{.FuncMapPtrToVal}},
		{{- end }}
		{{- range .Output }}
		"*{{.Type}}": {{.FuncMapPtrToPtr}},
		{{- end }}
	},
	{{- end }}
}

func dispatch(items reflect.Value, function interface{}, info *mapInfo) (bool,interface{}) {
	input:=info.fnInputType.String()	
	output:=info.fnOutputType.String()
	if inputVal,ok:=dispatcher[input];ok{
		if outputVal,ok:=inputVal[output];ok{
			return true, outputVal(items,function)
		}
    }
	return false, nil
}

{{- range .Data }}

{{- range .Output }}

func {{.FuncMapValToVal}}(itemsValue reflect.Value, function interface{}) interface{}  {
	input := itemsValue.Interface().([]{{.ParentType}})
	output := make([]{{.Type}},len(input))
	fn := function.(func({{.ParentType}}) {{.Type}})
	for i := 0; i < len(input); i++ {
		output[i] = fn(input[i])
	}
	return output
}

func {{.FuncMapValToPtr}}(itemsValue reflect.Value, function interface{}) interface{}  {
	input := itemsValue.Interface().([]{{.ParentType}})
	output := make([]*{{.Type}},len(input))
	fn := function.(func({{.ParentType}}) *{{.Type}})
	for i := 0; i < len(input); i++ {
		output[i] = fn(input[i])
	}
	return output
}


func {{.FuncMapPtrToVal}}(itemsValue reflect.Value, function interface{}) interface{}  {
	input := itemsValue.Interface().([]*{{.ParentType}})
	output := make([]{{.Type}},len(input))
	fn := function.(func(*{{.ParentType}}) {{.Type}})
	for i := 0; i < len(input); i++ {
		output[i] = fn(input[i])
	}
	return output
}


func {{.FuncMapPtrToPtr}}(itemsValue reflect.Value, function interface{}) interface{}  {
	input := itemsValue.Interface().([]*{{.ParentType}})
	output := make([]*{{.Type}},len(input))
	fn := function.(func(*{{.ParentType}}) *{{.Type}})
	for i := 0; i < len(input); i++ {
		output[i] = fn(input[i])
	}
	return output
}
{{- end }}
{{- end }}
`



func (w *TypeOutput) FuncMapValToVal() string {
	return fmt.Sprintf("map%sTo%s", strings.Title(w.ParentType), strings.Title(w.Type))
}
func (w *TypeOutput) FuncMapValToPtr() string {
	return fmt.Sprintf("map%sToPtr%s", strings.Title(w.ParentType), strings.Title(w.Type))
}
func (w *TypeOutput) FuncMapPtrToVal() string {
	return fmt.Sprintf("mapPtr%sTo%s", strings.Title(w.ParentType), strings.Title(w.Type))
}

func (w *TypeOutput) FuncMapPtrToPtr() string {
	return fmt.Sprintf("mapPtr%sToPtr%s", strings.Title(w.ParentType), strings.Title(w.Type))
}

func GenerateMapDispatcher(outputPath string) {
	f, err := os.Create(fmt.Sprintf("%s/dispatcher.go", outputPath))
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	template.
		Must(template.New("").
			Parse(mapDispatcherTemplateText)).
		Execute(f, struct {
			Timestamp time.Time
			Data      []*TypeInput
		}{
			Timestamp: time.Now(),
			Data:      data,
		})

}
