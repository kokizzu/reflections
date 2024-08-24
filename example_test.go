package reflections_test

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"

	"github.com/oleiade/reflections"
)

type MyStruct struct {
	MyEmbeddedStruct
	FirstField  string `matched:"first tag"`
	SecondField int    `matched:"second tag"`
	ThirdField  string `unmatched:"third tag"`
}

type MyEmbeddedStruct struct {
	EmbeddedField string
}

func ExampleGetField() {
	s := MyStruct{
		FirstField:  "first value",
		SecondField: 2,
		ThirdField:  "third value",
	}

	fieldsToExtract := []string{"FirstField", "ThirdField"}

	for _, fieldName := range fieldsToExtract {
		value, err := reflections.GetField(s, fieldName)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(value)

		// output:
		// first value
		// third value
	}
}

func ExampleGetFieldKind() {
	s := MyStruct{
		FirstField:  "first value",
		SecondField: 2,
		ThirdField:  "third value",
	}

	var firstFieldKind reflect.Kind
	var secondFieldKind reflect.Kind
	var err error

	// GetFieldKind will return reflect.String
	firstFieldKind, err = reflections.GetFieldKind(s, "FirstField")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(firstFieldKind)

	// GetFieldKind will return reflect.Int
	secondFieldKind, err = reflections.GetFieldKind(s, "SecondField")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(secondFieldKind)

	// output:
	// string
	// int
}

func ExampleGetFieldType() {
	s := MyStruct{
		FirstField:  "first value",
		SecondField: 2,
		ThirdField:  "third value",
	}

	var firstFieldType string
	var secondFieldType string
	var err error

	// GetFieldType will return reflect.String
	firstFieldType, err = reflections.GetFieldType(s, "FirstField")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(firstFieldType)

	// GetFieldType will return reflect.Int
	secondFieldType, err = reflections.GetFieldType(s, "SecondField")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(secondFieldType)

	// output:
	// string
	// int
}

func ExampleGetFieldTag() {
	s := MyStruct{}

	tag, err := reflections.GetFieldTag(s, "FirstField", "matched")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(tag)

	tag, err = reflections.GetFieldTag(s, "ThirdField", "unmatched")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(tag)

	// output:
	// first tag
	// third tag
}

func ExampleHasField() {
	s := MyStruct{
		FirstField:  "first value",
		SecondField: 2,
		ThirdField:  "third value",
	}

	// has == true
	has, _ := reflections.HasField(s, "FirstField")
	fmt.Println(has)

	// has == false
	has, _ = reflections.HasField(s, "FourthField")
	fmt.Println(has)

	// output:
	// true
	// false
}

func ExampleFields() {
	s := MyStruct{
		FirstField:  "first value",
		SecondField: 2,
		ThirdField:  "third value",
	}

	var fields []string

	// Fields will list every structure exportable fields.
	// Here, it's content would be equal to:
	// []string{"FirstField", "SecondField", "ThirdField"}
	fields, _ = reflections.Fields(s)
	fmt.Println(fields)

	// output:
	// [MyEmbeddedStruct FirstField SecondField ThirdField]
}

func ExampleItems() {
	s := MyStruct{
		FirstField:  "first value",
		SecondField: 2,
		ThirdField:  "third value",
	}

	var structItems map[string]interface{}

	// Items will return a field name to
	// field value map
	structItems, _ = reflections.Items(s)
	fmt.Println(structItems)

	// output:
	// map[FirstField:first value MyEmbeddedStruct:{} SecondField:2 ThirdField:third value]
}

func ExampleItemsDeep() {
	s := MyStruct{
		FirstField:  "first value",
		SecondField: 2,
		ThirdField:  "third value",
		MyEmbeddedStruct: MyEmbeddedStruct{
			EmbeddedField: "embedded value",
		},
	}

	var structItems map[string]interface{}

	// ItemsDeep will return a field name to
	// field value map, including fields from
	// anonymous embedded structs
	structItems, _ = reflections.ItemsDeep(s)
	fmt.Println(structItems)

	// output:
	// map[EmbeddedField:embedded value FirstField:first value SecondField:2 ThirdField:third value]
}

func ExampleTags() {
	s := MyStruct{
		FirstField:  "first value",
		SecondField: 2,
		ThirdField:  "third value",
	}

	var structTags map[string]string

	// Tags will return a field name to tag content
	// map. Nota that only field with the tag name
	// you've provided which will be matched.
	// Here structTags will contain:
	// {
	//     "FirstField": "first tag",
	//     "SecondField": "second tag",
	// }
	structTags, _ = reflections.Tags(s, "matched")
	fmt.Println(structTags)

	// output:
	// map[FirstField:first tag MyEmbeddedStruct: SecondField:second tag ThirdField:]
}

func ExampleSetField() {
	s := MyStruct{
		FirstField:  "first value",
		SecondField: 2,
		ThirdField:  "third value",
	}

	// In order to be able to set the structure's values,
	// a pointer to it has to be passed to it.
	err := reflections.SetField(&s, "FirstField", "new value")
	if err != nil {
		log.Fatal(err)
	}

	// Note that if you try to set a field's value using the wrong type,
	// an error will be returned
	_ = reflections.SetField(&s, "FirstField", 123) // err != nil

	// output:
}

func ExampleGetFieldNameByTagValue() {
	type Order struct {
		Step     string `json:"order_step"`
		ID       string `json:"id"`
		Category string `json:"category"`
	}
	type Condition struct {
		Field string `json:"field"`
		Value string `json:"value"`
		Next  string `json:"next"`
	}

	// JSON data from external source
	orderJSON := `{
		"order_step": "cooking",
		"id": "45457-fv54f54",
		"category": "Pizzas"
	}`

	conditionJSON := `{
		"field": "order_step", 
		"value": "cooking",
		"next": "serve"
	}`

	// Storing JSON in corresponding Variables
	var order Order
	err := json.Unmarshal([]byte(orderJSON), &order)
	if err != nil {
		log.Fatal(err)
	}

	var condition Condition
	err = json.Unmarshal([]byte(conditionJSON), &condition)
	if err != nil {
		log.Fatal(err)
	}

	fieldName, _ := reflections.GetFieldNameByTagValue(order, "json", condition.Field)
	fmt.Println(fieldName)
	fieldValue, _ := reflections.GetField(order, fieldName)
	fmt.Println(fieldValue)

	// Output:
	// Step
	// cooking
}
