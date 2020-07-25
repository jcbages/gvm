package main

import (
	"io"
)

// ClassFile represents a .class JVM file
type ClassFile struct {
	Magic             U4
	MinorVersion      U2
	MajorVersion      U2
	ConstantPoolCount U2
	ConstantPool      *ConstantPool
	AccessFlags       U2
	ThisClass         U2
	SuperClass        U2
	InterfacesCount   U2
	Interfaces        []U2
	FieldsCount       U2
	Fields            []*FieldInfo
	MethodsCount      U2
	Methods           []*MethodInfo
	AttributesCount   U2
	Attributes        []*AttributeInfo
}

// NewClassFile creates a new ClassFile based on the given io.Reader
func NewClassFile(r io.Reader) *ClassFile {
	var c ClassFile

	c.Magic = ReadU4(r)
	c.MinorVersion = ReadU2(r)
	c.MajorVersion = ReadU2(r)

	c.ConstantPoolCount = ReadU2(r)
	c.ConstantPool = NewConstantPool(r, c.ConstantPoolCount-1)

	c.AccessFlags = ReadU2(r)
	c.ThisClass = ReadU2(r)
	c.SuperClass = ReadU2(r)

	c.InterfacesCount = ReadU2(r)
	c.Interfaces = readInterfaces(r, c.InterfacesCount)

	c.FieldsCount = ReadU2(r)
	c.Fields = readFields(r, c.FieldsCount)

	c.MethodsCount = ReadU2(r)
	c.Methods = readMethods(r, c.MethodsCount)

	c.AttributesCount = ReadU2(r)
	c.Attributes = readAttributes(r, c.AttributesCount)

	return &c
}

func readInterfaces(r io.Reader, size U2) []U2 {
	var result []U2
	for i := U2(0); i < size; i++ {
		result = append(result, ReadU2(r))
	}
	return result
}

func readFields(r io.Reader, size U2) []*FieldInfo {
	var result []*FieldInfo
	for i := U2(0); i < size; i++ {
		result = append(result, NewFieldInfo(r))
	}
	return result
}

func readAttributes(r io.Reader, size U2) []*AttributeInfo {
	var result []*AttributeInfo
	for i := U2(0); i < size; i++ {
		result = append(result, NewAttributeInfo(r))
	}
	return result
}

func readMethods(r io.Reader, size U2) []*MethodInfo {
	var result []*MethodInfo
	for i := U2(0); i < size; i++ {
		result = append(result, NewMethodInfo(r))
	}
	return result
}
