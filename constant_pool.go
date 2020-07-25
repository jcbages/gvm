package main

import (
	"fmt"
	"io"
	"math"
)

// CpTag is the type of the CpInfo struct
type CpTag int

// CpInfo tag types
// All in [0..18] except [0, 2, 13, 14]
const (
	Class              CpTag = 7
	FieldRef           CpTag = 9
	MethodRef          CpTag = 10
	InterfaceMethodRef CpTag = 11
	String             CpTag = 8
	Integer            CpTag = 3
	Float              CpTag = 4
	Long               CpTag = 5
	Double             CpTag = 6
	NameAndType        CpTag = 12
	Utf8               CpTag = 1
	MethodHandle       CpTag = 15
	MethodType         CpTag = 16
	InvokeDynamic      CpTag = 18
)

var tagName = map[CpTag]string{
	Class:              "Class",
	FieldRef:           "FieldRef",
	MethodRef:          "MethodRef",
	InterfaceMethodRef: "InterfaceMethodRef",
	String:             "String",
	Integer:            "Integer",
	Float:              "Float",
	Long:               "Long",
	Double:             "Double",
	NameAndType:        "NameAndType",
	Utf8:               "Utf8",
	MethodHandle:       "MethodHandle",
	MethodType:         "MethodType",
	InvokeDynamic:      "InvokeDynamic",
}

// ConstantPool represents an array of CpInfo entries
type ConstantPool struct {
	Info []*CpInfo
}

// This is the "singleton" instance of the constant pool
var pool ConstantPool

// CpInfo represents a constant pool entry
type CpInfo struct {
	Tag                      CpTag
	NameIndex                U2
	ClassIndex               U2
	NameAndTypeIndex         U2
	StringIndex              U2
	Bytes                    U4
	HighBytes                U4
	LowBytes                 U4
	DescriptorIndex          U2
	Length                   U2
	ReferenceKind            U1
	ReferenceIndex           U2
	BootstrapMethodAttrIndex U2
	Utf8Bytes                []U1
}

// NewConstantPool creates a new Constant Pool
func NewConstantPool(r io.Reader, size U2) *ConstantPool {
	for i := U2(0); i < size; i++ {
		pool.Info = append(pool.Info, NewCpInfo(r))
	}
	return &pool
}

// At fetch the CPInfo at the given index (one based)
func (pool *ConstantPool) At(i U2) *CpInfo {
	return pool.Info[i-1]
}

// NewCpInfo creates a new poolInfo based on the given io.Reader
func NewCpInfo(r io.Reader) *CpInfo {
	var c CpInfo

	c.Tag = CpTag(ReadU1(r))

	switch c.Tag {
	case Class:
		c.NameIndex = ReadU2(r)
	case FieldRef, MethodRef, InterfaceMethodRef:
		c.ClassIndex = ReadU2(r)
		c.NameAndTypeIndex = ReadU2(r)
	case String:
		c.StringIndex = ReadU2(r)
	case Integer, Float:
		c.Bytes = ReadU4(r)
	case Long, Double:
		c.HighBytes = ReadU4(r)
		c.LowBytes = ReadU4(r)
	case NameAndType:
		c.NameIndex = ReadU2(r)
		c.DescriptorIndex = ReadU2(r)
	case Utf8:
		c.Length = ReadU2(r)
		c.readUtf8Bytes(r)
	case MethodHandle:
		c.ReferenceKind = ReadU1(r)
		c.ReferenceIndex = ReadU2(r)
	case MethodType:
		c.DescriptorIndex = ReadU2(r)
	case InvokeDynamic:
		c.BootstrapMethodAttrIndex = ReadU2(r)
		c.NameAndTypeIndex = ReadU2(r)
	}

	return &c
}

func (c *CpInfo) readUtf8Bytes(r io.Reader) {
	for i := U2(0); i < c.Length; i++ {
		c.Utf8Bytes = append(c.Utf8Bytes, ReadU1(r))
	}
}

// Utf8 returns the string value for Utf8 CpInfo
func (c *CpInfo) Utf8() string {
	var s []byte
	for _, b := range c.Utf8Bytes {
		s = append(s, byte(b))
	}
	return string(s)
}

// Integer returns the signed int value for Integer CpInfo
func (c *CpInfo) Integer() int32 {
	return int32(c.Bytes)
}

// Float returns the float value for Float CpInfo
func (c *CpInfo) Float() float32 {
	return math.Float32frombits(uint32(c.Bytes))
}

// Long returns the signed long value for Long CpInfo
func (c *CpInfo) Long() int64 {
	v := uint64(c.HighBytes)*32 + uint64(c.LowBytes)
	return int64(v)
}

// Double returns the float64 value for Double CpInfo
func (c *CpInfo) Double() float64 {
	v := uint64(c.HighBytes)*32 + uint64(c.LowBytes)
	return math.Float64frombits(uint64(v))
}

// Value finds the utf8 values for the CpInfo and its attributes
func (c *CpInfo) Value() string {
	var s string

	switch c.Tag {
	case Class:
		s = pool.At(c.NameIndex).Value()
	case FieldRef, MethodRef, InterfaceMethodRef:
		s = fmt.Sprintf(
			"ClassIndex=%v NameAndTypeIndex=%v",
			pool.At(c.ClassIndex).Value(),
			pool.At(c.NameAndTypeIndex).Value(),
		)
	case String:
		s = pool.At(c.StringIndex).Value()
	case Integer:
		s = fmt.Sprintf("%v", c.Integer())
	case Float:
		s = fmt.Sprintf("%v", c.Float())
	case Long:
		s = fmt.Sprintf("%v", c.Long())
	case Double:
		s = fmt.Sprintf("%v", c.Double())
	case NameAndType:
		s = fmt.Sprintf(
			"NameIndex=%v DescriptorIndex=%v",
			pool.At(c.NameIndex).Value(),
			pool.At(c.DescriptorIndex).Value(),
		)
	case Utf8:
		s = c.Utf8()
	case MethodHandle:
		s = fmt.Sprintf(
			"ReferenceKind=%v ReferenceIndex=%v",
			c.ReferenceKind,
			pool.At(c.ReferenceIndex).Value(),
		)
	case MethodType:
		s = pool.At(c.DescriptorIndex).Value()
	case InvokeDynamic:
		s = fmt.Sprintf(
			"BootstrapMethodAt=%v NameAndTypeIndex=%v",
			pool.At(c.BootstrapMethodAttrIndex).Value(),
			pool.At(c.NameAndTypeIndex).Value(),
		)
	}

	return s
}

func (c *CpInfo) String() string {
	return fmt.Sprintf("%v{%v}", tagName[c.Tag], c.Value())
}
