package main

/*
#include <stdlib.h>
*/
import "C"
import (
	"bytes"
	"fmt"
	"encoding/gob"
	"encoding/json"
	"reflect"
	"unsafe"
)

// go build -buildmode=c-shared -o gob.so
var allModels = make(map[string]interface{})

func gobEncode(o interface{}) []byte {
	b := new(bytes.Buffer)
	enc := gob.NewEncoder(b)
	if err := enc.Encode(&o); err != nil {
		return nil
	}
	return b.Bytes()
}

func gobDecodeValue(o interface{}, b []byte) interface{} {
	t := reflect.TypeOf(o)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	newSt := reflect.New(t)
	dec := gob.NewDecoder(bytes.NewBuffer(b))
	if err := dec.DecodeValue(newSt); err == nil {
		//fmt.Println(err)
		return newSt.Interface()
	}
	return nil
}

func gobDecodeAll(b []byte) map[string]interface{} {
	var os = make(map[string]interface{})
	for i := range allModels {
		o := gobDecodeValue(allModels[i], b)
		if o != nil {
			os[fmt.Sprintf("Struct-%s", i)] = o
		}
	}
	return os
}

func gobDecode(b []byte) interface{} {
	var o interface{}
	o = gobDecodeValue(&o, b)
	if o != nil {
		return o
	}
	return gobDecodeAll(b)
}

func gobDecodeWithName(n string, b []byte) interface{} {
	if v, ok := allModels[n]; ok {
		return gobDecodeValue(v, b)
	}
	return gobDecode(b)
}

func jsonFormat(o interface{}) []byte {
	b, _ := json.MarshalIndent(o, "", "    ")
	return b
}

func gobFormat(name string, b []byte) []byte {
	return jsonFormat(gobDecodeWithName(name, b))
}

func gobRegister(name string, o interface{}) {
	gob.Register(o)
	allModels[name] = o
}

//export CGobFormat
func CGobFormat(n *C.char, b []byte) *C.char {
	return C.CString(string(gobFormat(C.GoString(n), b)))
}

//export CGobFree
func CGobFree(cstr *C.char) {
	C.free(unsafe.Pointer(cstr))
}

func main() {}
