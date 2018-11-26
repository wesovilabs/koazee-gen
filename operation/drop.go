package operation

import (
	"fmt"
	"log"
	"os"
	"strings"
	"text/template"
	"time"
)

const dropDispatcherTemplateText = `
package drop

import "reflect"


type dispatchFunction func(items reflect.Value, item interface{}) interface{}

func reverseIndexes(indexes []int) []int{
	for i := len(indexes)/2 - 1; i >= 0; i-- {
		opp := len(indexes) - 1 - i
		indexes[i], indexes[opp] = indexes[opp], indexes[i]
	}
	return indexes
}

var dispatcher = map[string]dispatchFunction{
	{{- range .Data }}
	"{{.Type}}":  {{.FuncDropVal}},
	"*{{.Type}}":  {{.FuncDropPtr}},
	{{- end }}
}

func dispatch(items reflect.Value, itemValue interface{}, info *dropInfo) (bool,interface{}) {
	input:=(*info.itemType).String()
	if fnVal,ok:=dispatcher[input];ok{
		return true, fnVal(items,itemValue)
    }
	return false, nil
}

{{- range .Data }}

func {{.FuncDropVal}}(itemsValue reflect.Value, itemValue interface{}) interface{} {
	input := itemsValue.Interface().([]{{.Type}})
	indexes := make([]int, 0)
	item := itemValue.({{.Type}})
	for i := 0; i < len(input); i++ {
		if input[i] == item {
			indexes = append(indexes, i)
		}
	}
	indexes = reverseIndexes(indexes)
	for _, index := range indexes {
		if index > 0 {
			input = append(input[:index], input[index+1:]...)
			continue
		}
		input = input[index+1:]
	}
	return input
}

func {{.FuncDropPtr}}(itemsValue reflect.Value, itemValue interface{}) interface{} {
	input := itemsValue.Interface().([]*{{.Type}})
	indexes := make([]int, 0)
	item := itemValue.(*{{.Type}})
	for i := 0; i < len(input); i++ {
		if input[i] == item {
			indexes = append(indexes, i)
		}
	}
	indexes = reverseIndexes(indexes)
	for _, index := range indexes {
		input = append(input[:index], input[index+1:]...)
	}
	return input
}

{{- end }}
`

func (w *TypeInput) FuncDropVal() string {
	return fmt.Sprintf("drop%s", strings.Title(w.Type))
}
func (w *TypeInput) FuncDropPtr() string {
	return fmt.Sprintf("dropPtr%s", strings.Title(w.Type))
}

func GenerateDropDispatcher(outputPath string) {
	f, err := os.Create(fmt.Sprintf("%s/dispatcher.go", outputPath))
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	template.
		Must(template.New("").
			Parse(dropDispatcherTemplateText)).
		Execute(f, struct {
			Timestamp time.Time
			Data      []*TypeInput
		}{
			Timestamp: time.Now(),
			Data:      data,
		})

}
