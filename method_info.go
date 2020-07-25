package main

import (
	"io"
)

// MethodInfo represents a method
type MethodInfo = FieldInfo

// NewMethodInfo creates a new MethodInfo
func NewMethodInfo(r io.Reader) *MethodInfo {
	return NewFieldInfo(r)
}
