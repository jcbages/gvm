package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run *.go [.class filename]")
		os.Exit(1)
	}

	file, err := os.Open(os.Args[1])
	if err != nil {
		panic(err)
	}

	cf := NewClassFile(file)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Magic Number: %x\n", cf.Magic)
	fmt.Printf("Version: %v.%v\n", cf.MajorVersion, cf.MinorVersion)

	fmt.Printf("Constant Pool Count: %v\n", cf.ConstantPoolCount)
	// for i, v := range cf.ConstantPool.Info {
	// 	fmt.Printf("ConstantPool [%v]: %v\n", i+1, v)
	// }

	fmt.Printf("Access Flags: 0x%04x\n", cf.AccessFlags)
	fmt.Printf("This Class: %v\n", cf.ConstantPool.At(cf.ThisClass))
	fmt.Printf("Super Class: %v\n", cf.ConstantPool.At(cf.SuperClass))

	fmt.Printf("Interfaces Count: %v\n", cf.InterfacesCount)
	for i, v := range cf.Interfaces {
		fmt.Printf("Interface [%v]: %v\n", i+1, v)
	}

	fmt.Printf("Fields Count: %v\n", cf.FieldsCount)
	for i, v := range cf.Fields {
		fmt.Printf("Field [%v]: %v\n", i+1, v)
		for j, w := range v.Attributes {
			fmt.Printf("    Attribute [%v]: %v\n", j+1, w)
		}
	}

	fmt.Printf("Methods Count: %v\n", cf.MethodsCount)
	for i, v := range cf.Methods {
		fmt.Printf("Method [%v]: %v\n", i+1, v)
		for j, w := range v.Attributes {
			fmt.Printf("    Attribute [%v]: %v\n", j+1, w)
		}
	}

	fmt.Printf("Attributes Count: %v\n", cf.AttributesCount)
	for i, v := range cf.Attributes {
		fmt.Printf("Attribute [%v]: %v\n", i+1, v)
	}
}
