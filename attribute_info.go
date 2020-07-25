package main

import (
	"fmt"
	"io"
	"strings"
)

// AttributeName represent the possible values for NameIndex in the Constant Pool
type AttributeName string

// FrameType represents a type of StackMapFrame
type FrameType U1

// FrameTypeName represents the name of a FrameType
type FrameTypeName string

// VariableInfo represents the tag of the VerificationTypeInfo
type VariableInfo U1

// The values that AttributeName can take
const (
	ConstantValue    AttributeName = "ConstantValue"
	Code             AttributeName = "Code"
	StackMapTable    AttributeName = "StackMapTable"
	Exceptions       AttributeName = "Exceptions"
	BootstrapMethods AttributeName = "BootstrapMethods"
)

// Tag types for the Frame Type
const (
	Same                           FrameTypeName = "Same"
	SameLocalsOneStackItem         FrameTypeName = "SameLocalsOneStackItem"
	SameLocalsOneStackItemExtended FrameTypeName = "SameLocalsOneStackItemExtended"
	ChopFrame                      FrameTypeName = "ChopFrame"
	SameFrameExtended              FrameTypeName = "SameFrameExtended"
	AppendFrame                    FrameTypeName = "AppendFrame"
	FullFrame                      FrameTypeName = "FullFrame"
)

// The values that VariableInfo can take
const (
	ItemTop               VariableInfo = 0
	ItemInteger           VariableInfo = 1
	ItemFloat             VariableInfo = 2
	ItemNull              VariableInfo = 5
	ItemUninitializedThis VariableInfo = 6
	ItemObject            VariableInfo = 7
	ItemUninitialized     VariableInfo = 8
	ItemLong              VariableInfo = 4
	ItemDouble            VariableInfo = 3
)

// AttributeInfo represents an attribute
type AttributeInfo struct {
	NameIndex U2
	Length    U4
	Info      []U1

	// ConstantValue
	ConstantValueIndex U2

	// Code
	MaxStack             U2
	MaxLocals            U2
	CodeLength           U4
	Code                 []U1
	ExceptionTableLength U2
	ExceptionTable       []*ExceptionTableInfo
	AttributesCount      U2
	Attributes           []*AttributeInfo

	// StackMapTable
	NumberOfEntries U2
	Entries         []*StackMapFrame

	// Exceptions
	NumberOfExceptions  U2
	ExceptionIndexTable []U2

	// BootstrapMethods
	NumberOfBootstrapMethods U2
	BootstrapMethods         []*BootstrapMethodInfo
}

// ExceptionTableInfo represents an entry in the ExceptionTable attribute
type ExceptionTableInfo struct {
	StartPC   U2
	EndPC     U2
	HandlerPC U2
	CatchType U2
}

// StackMapFrame represents an entry in the StackTableMap attribute
type StackMapFrame struct {
	FrameType          FrameType
	OffsetDelta        U2
	Stack              []*VerificationTypeInfo
	Locals             []*VerificationTypeInfo
	TrimVariableCount  U1
	NumberOfLocals     U2
	NumberOfStackItems U2
}

// VerificationTypeInfo represents an entry in the Stack or Locals array of StackMapFrame
type VerificationTypeInfo struct {
	Tag        VariableInfo
	CPoolIndex U2
	Offset     U2
}

// BootstrapMethodInfo represents an entry in the BootstrapMethods attribute
type BootstrapMethodInfo struct {
	BootstrapMethodRef         U2
	NumberOfBootstrapArguments U2
	BootstrapArguments         []U2
}

// NewAttributeInfo creates a new AttributeInfo
func NewAttributeInfo(r io.Reader) *AttributeInfo {
	var attr AttributeInfo

	attr.NameIndex = ReadU2(r)
	attr.Length = ReadU4(r)

	switch v := pool.At(attr.NameIndex).Utf8(); AttributeName(v) {
	case ConstantValue:
		attr.ConstantValueIndex = ReadU2(r)
	case Code:
		attr.readCode(r)
	case StackMapTable:
		attr.readStackMapTable(r)
	case Exceptions:
		attr.readExceptions(r)
	case BootstrapMethods:
		attr.readBootstrapMethods(r)
	default:
		attr.readDefault(r)
	}

	return &attr
}

func (attr *AttributeInfo) readCode(r io.Reader) {
	attr.MaxStack = ReadU2(r)
	attr.MaxLocals = ReadU2(r)

	attr.CodeLength = ReadU4(r)
	for i := U4(0); i < attr.CodeLength; i++ {
		attr.Code = append(attr.Code, ReadU1(r))
	}

	attr.ExceptionTableLength = ReadU2(r)
	for i := U2(0); i < attr.ExceptionTableLength; i++ {
		attr.ExceptionTable = append(attr.ExceptionTable, NewExceptionTableInfo(r))
	}

	attr.AttributesCount = ReadU2(r)
	for i := U2(0); i < attr.AttributesCount; i++ {
		attr.Attributes = append(attr.Attributes, NewAttributeInfo(r))
	}
}

func (attr *AttributeInfo) readStackMapTable(r io.Reader) {
	attr.NumberOfEntries = ReadU2(r)
	for i := U2(0); i < attr.NumberOfEntries; i++ {
		attr.Entries = append(attr.Entries, NewStackMapFrame(r))
	}
}

func (attr *AttributeInfo) readExceptions(r io.Reader) {
	attr.NumberOfExceptions = ReadU2(r)
	for i := U2(0); i < attr.NumberOfExceptions; i++ {
		attr.ExceptionIndexTable = append(attr.ExceptionIndexTable, ReadU2(r))
	}
}

func (attr *AttributeInfo) readBootstrapMethods(r io.Reader) {
	attr.NumberOfBootstrapMethods = ReadU2(r)
	for i := U2(0); i < attr.NumberOfBootstrapMethods; i++ {
		attr.BootstrapMethods = append(attr.BootstrapMethods, NewBootstrapMethodInfo(r))
	}
}

func (attr *AttributeInfo) readDefault(r io.Reader) {
	for i := U4(0); i < attr.Length; i++ {
		attr.Info = append(attr.Info, ReadU1(r))
	}
}

// NewExceptionTableInfo creates a new ExceptionTableInfo struct
func NewExceptionTableInfo(r io.Reader) *ExceptionTableInfo {
	var e ExceptionTableInfo

	e.StartPC = ReadU2(r)
	e.EndPC = ReadU2(r)
	e.HandlerPC = ReadU2(r)
	e.CatchType = ReadU2(r)

	return &e
}

// GetName returns the name of the frame type
func (ft FrameType) GetName() FrameTypeName {
	v := int(ft)
	if v >= 0 && v <= 63 {
		return Same
	} else if v >= 64 && ft <= 127 {
		return SameLocalsOneStackItem
	} else if v == 247 {
		return SameLocalsOneStackItemExtended
	} else if v >= 248 && v <= 250 {
		return ChopFrame
	} else if v == 251 {
		return SameFrameExtended
	} else if v >= 252 && v <= 254 {
		return AppendFrame
	} else {
		return FullFrame
	}
}

// NewStackMapFrame creates a new StackMapFrame struct
func NewStackMapFrame(r io.Reader) *StackMapFrame {
	var s StackMapFrame

	s.FrameType = FrameType(ReadU1(r))
	switch s.FrameType.GetName() {
	case Same:
		s.OffsetDelta = U2(s.FrameType)
	case SameLocalsOneStackItem:
		s.OffsetDelta = U2(s.FrameType - 64)
		s.NumberOfStackItems = U2(1)
		s.readStack(r)
	case SameLocalsOneStackItemExtended:
		s.OffsetDelta = ReadU2(r)
		s.NumberOfStackItems = U2(1)
		s.readStack(r)
	case ChopFrame:
		s.OffsetDelta = ReadU2(r)
		s.TrimVariableCount = U1(251 - s.FrameType)
	case SameFrameExtended:
		s.OffsetDelta = ReadU2(r)
	case AppendFrame:
		s.OffsetDelta = ReadU2(r)
		s.NumberOfLocals = U2(s.FrameType - 251)
		s.readLocals(r)
	case FullFrame:
		s.OffsetDelta = ReadU2(r)
		s.NumberOfLocals = ReadU2(r)
		s.readLocals(r)
		s.NumberOfStackItems = ReadU2(r)
		s.readStack(r)
	}

	return &s
}

func (s *StackMapFrame) readStack(r io.Reader) {
	for i := U2(0); i < s.NumberOfStackItems; i++ {
		s.Stack = append(s.Stack, NewVerificationTypeInfo(r))
	}
}

func (s *StackMapFrame) readLocals(r io.Reader) {
	for i := U2(0); i < s.NumberOfLocals; i++ {
		s.Locals = append(s.Locals, NewVerificationTypeInfo(r))
	}
}

// NewVerificationTypeInfo returns a new VerificationTypeInfo struct
func NewVerificationTypeInfo(r io.Reader) *VerificationTypeInfo {
	var v VerificationTypeInfo

	v.Tag = VariableInfo(ReadU1(r))

	if v.Tag == ItemObject {
		v.CPoolIndex = ReadU2(r)
	}

	if v.Tag == ItemUninitialized {
		v.Offset = ReadU2(r)
	}

	return &v
}

// NewBootstrapMethodInfo created a new BootstrapMethodInfo struct
func NewBootstrapMethodInfo(r io.Reader) *BootstrapMethodInfo {
	var b BootstrapMethodInfo

	b.BootstrapMethodRef = ReadU2(r)

	b.NumberOfBootstrapArguments = ReadU2(r)
	for i := U2(0); i < b.NumberOfBootstrapArguments; i++ {
		b.BootstrapArguments = append(b.BootstrapArguments, ReadU2(r))
	}

	return &b
}

func (ai *AttributeInfo) String() string {
	var attrs []string
	for _, v := range ai.Attributes {
		attrs = append(attrs, v.String())
	}

	return fmt.Sprintf(
		"AttributeInfo{NameIndex=%v Length=%v, Attrs=%v}",
		pool.At(ai.NameIndex).String(),
		ai.Length,
		"["+strings.Join(attrs, ",")+"]",
	)
}
