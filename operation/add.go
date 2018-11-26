package operation

import (
	"fmt"
	"log"
	"os"
	"strings"
	"text/template"
	"time"
)

const addDispatcherTemplateText = `
package add

import "reflect"


type dispatchFunction func(items reflect.Value, item interface{}) interface{}

var dispatcher = map[string]dispatchFunction{
	{{- range .Data }}
	"{{.Type}}":  {{.FuncAddVal}},
	"*{{.Type}}":  {{.FuncAddPtr}},
	{{- end }}
}

func dispatch(items reflect.Value, itemValue interface{}, info *addInfo) (bool,interface{}) {
	input:=(*info.itemType).String()
	if fnVal,ok:=dispatcher[input];ok{
		return true, fnVal(items,itemValue)
    }
	return false, nil
}

{{- range .Data }}

func {{.FuncAddVal}}(itemsValue reflect.Value, itemValue interface{}) interface{} {
	input := itemsValue.Interface().([]{{.Type}})
	item := itemValue.({{.Type}})
	return append(input,item)
}

func {{.FuncAddPtr}}(itemsValue reflect.Value, itemValue interface{}) interface{} {
	input := itemsValue.Interface().([]*{{.Type}})
	item := itemValue.(*{{.Type}})
	return append(input,item)
}

{{- end }}
`

func (w *TypeInput) FuncAddVal() string {
	return fmt.Sprintf("add%s", strings.Title(w.Type))
}
func (w *TypeInput) FuncAddPtr() string {
	return fmt.Sprintf("addPtr%s", strings.Title(w.Type))
}

func GenerateAddDispatcher(outputPath string) {
	f, err := os.Create(fmt.Sprintf("%s/dispatcher.go", outputPath))
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	template.
		Must(template.New("").
			Parse(addDispatcherTemplateText)).
		Execute(f, struct {
			Timestamp time.Time
			Data      []*TypeInput
		}{
			Timestamp: time.Now(),
			Data:      data,
		})

}
