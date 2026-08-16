package shared

import (
	"fmt"
	"strconv"
	"strings"
)

type JSONPointer interface {
	fmt.Stringer
}

type jsonPointer struct {
	segments []any
}

var RootPointer JSONPointer = &jsonPointer{segments: []any{}}

func (jp jsonPointer) String() string {
	var b strings.Builder
	b.WriteString("#")
	for _, segment := range jp.segments {
		b.WriteString("/")
		switch v := segment.(type) {
		case string:
			b.WriteString(v)
		case int:
			b.WriteString(strconv.Itoa(v))
		}
	}
	return b.String()
}

type JSONPointerBuilder interface {
	Key(key string) JSONPointerBuilder
	Index(index int) JSONPointerBuilder
	Build() JSONPointer
}

type jsonPointerBuilder struct {
	segments []any
}

func NewJSONPointerBuilder() JSONPointerBuilder {
	return &jsonPointerBuilder{}
}

func (b *jsonPointerBuilder) Key(key string) JSONPointerBuilder {
	b.segments = append(b.segments, key)
	return b
}

func (b *jsonPointerBuilder) Index(index int) JSONPointerBuilder {
	b.segments = append(b.segments, index)
	return b
}

func (b *jsonPointerBuilder) Build() JSONPointer {
	// because we're not making a copy of the segment slice, the builder should not be reused after calling Build()
	return &jsonPointer{segments: b.segments}
}
