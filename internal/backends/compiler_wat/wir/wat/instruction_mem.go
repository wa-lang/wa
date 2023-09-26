// 版权 @2022 凹语言 作者。保留所有权利。

package wat

import (
	"strconv"
	"strings"
)

/**************************************
instLoad:
**************************************/
type instLoad struct {
	anInstruction
	typ           ValueType
	offset, align int
}

func NewInstLoad(typ ValueType, offset int, align int) *instLoad {
	return &instLoad{typ: typ, offset: offset, align: align}
}
func (i *instLoad) Format(indent string, sb *strings.Builder) {
	sb.WriteString(indent)
	sb.WriteString(i.typ.Name())
	sb.WriteString(".load offset=")
	sb.WriteString(strconv.Itoa(i.offset))
	sb.WriteString(" align=")
	sb.WriteString(strconv.Itoa(i.align))
}

/**************************************
instLoad8s:
**************************************/
type instLoad8s struct {
	anInstruction
	offset, align int
}

func NewInstLoad8s(offset int, align int) *instLoad8s {
	return &instLoad8s{offset: offset, align: align}
}
func (i *instLoad8s) Format(indent string, sb *strings.Builder) {
	sb.WriteString(indent)
	sb.WriteString("i32.load8_s offset=")
	sb.WriteString(strconv.Itoa(i.offset))
	sb.WriteString(" align=")
	sb.WriteString(strconv.Itoa(i.align))
}

/**************************************
instLoad8u:
**************************************/
type instLoad8u struct {
	anInstruction
	offset, align int
}

func NewInstLoad8u(offset int, align int) *instLoad8u {
	return &instLoad8u{offset: offset, align: align}
}
func (i *instLoad8u) Format(indent string, sb *strings.Builder) {
	sb.WriteString(indent)
	sb.WriteString("i32.load8_u offset=")
	sb.WriteString(strconv.Itoa(i.offset))
	sb.WriteString(" align=")
	sb.WriteString(strconv.Itoa(i.align))
}

/**************************************
instLoad16s:
**************************************/
type instLoad16s struct {
	anInstruction
	offset, align int
}

func NewInstLoad16s(offset int, align int) *instLoad16s {
	return &instLoad16s{offset: offset, align: align}
}
func (i *instLoad16s) Format(indent string, sb *strings.Builder) {
	sb.WriteString(indent)
	sb.WriteString("i32.load16_s offset=")
	sb.WriteString(strconv.Itoa(i.offset))
	sb.WriteString(" align=")
	sb.WriteString(strconv.Itoa(i.align))
}

/**************************************
instLoad16u:
**************************************/
type instLoad16u struct {
	anInstruction
	offset, align int
}

func NewInstLoad16u(offset int, align int) *instLoad16u {
	return &instLoad16u{offset: offset, align: align}
}
func (i *instLoad16u) Format(indent string, sb *strings.Builder) {
	sb.WriteString(indent)
	sb.WriteString("i32.load16_u offset=")
	sb.WriteString(strconv.Itoa(i.offset))
	sb.WriteString(" align=")
	sb.WriteString(strconv.Itoa(i.align))
}

/**************************************
instStore:
**************************************/
type instStore struct {
	anInstruction
	typ           ValueType
	offset, align int
}

func NewInstStore(typ ValueType, offset int, align int) *instStore {
	return &instStore{typ: typ, offset: offset, align: align}
}
func (i *instStore) Format(indent string, sb *strings.Builder) {
	sb.WriteString(indent)
	sb.WriteString(i.typ.Name())
	sb.WriteString(".store offset=")
	sb.WriteString(strconv.Itoa(i.offset))
	sb.WriteString(" align=")
	sb.WriteString(strconv.Itoa(i.align))
}

/**************************************
instStore8:
**************************************/
type instStore8 struct {
	anInstruction
	offset, align int
}

func NewInstStore8(offset int, align int) *instStore8 {
	return &instStore8{offset: offset, align: align}
}
func (i *instStore8) Format(indent string, sb *strings.Builder) {
	sb.WriteString(indent)
	sb.WriteString("i32.store8 offset=")
	sb.WriteString(strconv.Itoa(i.offset))
	sb.WriteString(" align=")
	sb.WriteString(strconv.Itoa(i.align))
}

/**************************************
instStore16:
**************************************/
type instStore16 struct {
	anInstruction
	offset, align int
}

func NewInstStore16(offset int, align int) *instStore16 {
	return &instStore16{offset: offset, align: align}
}
func (i *instStore16) Format(indent string, sb *strings.Builder) {
	sb.WriteString(indent)
	sb.WriteString("i32.store16 offset=")
	sb.WriteString(strconv.Itoa(i.offset))
	sb.WriteString(" align=")
	sb.WriteString(strconv.Itoa(i.align))
}
