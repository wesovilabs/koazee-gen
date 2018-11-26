package operation

import (
	"fmt"
	"log"
	"os"
	"strings"
	"text/template"
	"time"
)

const reduceDispatcherTemplateText = `
package reduce

import "reflect"


type dispatchFunction func(items reflect.Value, fn interface{}) interface{}

var dispatcher = map[string]map[string]dispatchFunction{
	{{- range .Data }}
	"{{.Type}}": {
		{{- range .Output }}
		"{{.Type}}": {{.FuncReduceValToVal}},
		{{- end }}
		{{- range .Output }}
		"*{{.Type}}": {{.FuncReduceValToPtr}},
		{{- end }}
	},
	"*{{.Type}}": {
		{{- range .Output }}
		"{{.Type}}": {{.FuncReducePtrToVal}},
		{{- end }}
		{{- range .Output }}
		"*{{.Type}}": {{.FuncReducePtrToPtr}},
		{{- end }}
	},
	{{- end }}
}

func dispatch(items reflect.Value, function interface{}, info *reduceInfo) (bool,interface{}) {
	output:=info.fnIn1Type.String()
	input:=info.fnIn2Type.String()
	if inputVal,ok:=dispatcher[input];ok{
		if outputVal,ok:=inputVal[output];ok{
			return true, outputVal(items,function)
		}
    }
	return false, nil
}

{{- range .Data }}

{{- range .Output }}

func {{.FuncReduceValToVal}}(itemsValue reflect.Value, function interface{}) interface{} {
	items := itemsValue.Interface().([]{{.ParentType}})
	fn := function.(func({{.Type}}, {{.ParentType}}) {{.Type}})
	acc := {{.Value}}
	for _, item := range items {
		acc = fn(acc, item)
	}
	return acc
}

func {{.FuncReduceValToPtr}}(itemsValue reflect.Value, function interface{}) interface{} {
	items := itemsValue.Interface().([]{{.ParentType}})
	fn := function.(func(*{{.Type}}, {{.ParentType}}) *{{.Type}})
	acc := {{.Value}}
	accPtr := &acc
	for _, item := range items {
		accPtr = fn(accPtr, item)
	}
	return accPtr
}

func {{.FuncReducePtrToVal}}(itemsValue reflect.Value, function interface{}) interface{} {
	items := itemsValue.Interface().([]*{{.ParentType}})
	fn := function.(func(*{{.Type}}, *{{.ParentType}}) *{{.Type}})
	acc := {{.Value}}
	accPtr := &acc
	for _, item := range items {
		accPtr = fn(accPtr, item)
	}
	return accPtr
}


func {{.FuncReducePtrToPtr}}(itemsValue reflect.Value, function interface{}) interface{} {
	items := itemsValue.Interface().([]*{{.ParentType}})
	fn := function.(func({{.Type}}, *{{.ParentType}}) {{.Type}})
	acc := {{.Value}}
	for _, item := range items {
		acc = fn(acc, item)
	}
	return acc
}
{{- end }}
{{- end }}
`
func (w *TypeOutput) FuncReduceValToVal() string {
	return fmt.Sprintf("reduce%sTo%s", strings.Title(w.ParentType), strings.Title(w.Type))
}
func (w *TypeOutput) FuncReduceValToPtr() string {
	return fmt.Sprintf("reduce%sToPtr%s", strings.Title(w.ParentType), strings.Title(w.Type))
}
func (w *TypeOutput) FuncReducePtrToVal() string {
	return fmt.Sprintf("reducePtr%sTo%s", strings.Title(w.ParentType), strings.Title(w.Type))
}

func (w *TypeOutput) FuncReducePtrToPtr() string {
	return fmt.Sprintf("reducePtr%sToPtr%s", strings.Title(w.ParentType), strings.Title(w.Type))
}




func GenerateReduceDispatcher(outputPath string) {
	f, err := os.Create(fmt.Sprintf("%s/dispatcher.go", outputPath))
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	template.
		Must(template.New("").
			Parse(reduceDispatcherTemplateText)).
		Execute(f, struct {
			Timestamp time.Time
			Data      []*TypeInput
		}{
			Timestamp: time.Now(),
			Data:      data,
		})

}