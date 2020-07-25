package main

import (
	"fmt"
	"io"
)

// FieldInfo represents a field
type FieldInfo struct {
	AccessFlags     U2
	NameIndex       U2
	DescriptorIndex U2
	AttributesCount U2
	Attributes      []*AttributeInfo
}

// NewFieldInfo creates a new FieldInfo
func NewFieldInfo(r io.Reader) *FieldInfo {
	var fi FieldInfo

	fi.AccessFlags = ReadU2(r)
	fi.NameIndex = ReadU2(r)
	fi.DescriptorIndex = ReadU2(r)

	fi.AttributesCount = ReadU2(r)
	fi.Attributes = readAttributes(r, fi.AttributesCount)

	return &fi
}

func (fi *FieldInfo) String() string {
	return fmt.Sprintf(
		"Info{AccessFlags=0x%04x NameIndex=%v DescriptorIndex=%v AttributesCount=%v}",
		fi.AccessFlags,
		pool.At(fi.NameIndex).String(),
		pool.At(fi.DescriptorIndex).String(),
		fi.AttributesCount,
	)
}
