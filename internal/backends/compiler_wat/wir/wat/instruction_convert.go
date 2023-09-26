package wat

import "strings"

/**************************************
instConvert_i32_wrap_i64:
**************************************/
type instConvert_i32_wrap_i64 struct {
	anInstruction
}

func NewInstConvert_i32_wrap_i64() *instConvert_i32_wrap_i64 { return &instConvert_i32_wrap_i64{} }
func (i *instConvert_i32_wrap_i64) Format(indent string, sb *strings.Builder) {
	sb.WriteString(indent)
	sb.WriteString("i32.wrap_i64")
}

/**************************************
instConvert_i32_trunc_f32_s:
**************************************/
type instConvert_i32_trunc_f32_s struct {
	anInstruction
}

func NewInstConvert_i32_trunc_f32_s() *instConvert_i32_trunc_f32_s {
	return &instConvert_i32_trunc_f32_s{}
}
func (i *instConvert_i32_trunc_f32_s) Format(indent string, sb *strings.Builder) {
	sb.WriteString(indent)
	sb.WriteString("i32.trunc_f32_s")
}

/**************************************
instConvert_i32_trunc_f64_s:
**************************************/
type instConvert_i32_trunc_f64_s struct {
	anInstruction
}

func NewInstConvert_i32_trunc_f64_s() *instConvert_i32_trunc_f64_s {
	return &instConvert_i32_trunc_f64_s{}
}
func (i *instConvert_i32_trunc_f64_s) Format(indent string, sb *strings.Builder) {
	sb.WriteString(indent)
	sb.WriteString("i32.trunc_f64_s")
}

/**************************************
instConvert_i64_extend_i32_s:
**************************************/
type instConvert_i64_extend_i32_s struct {
	anInstruction
}

func NewInstConvert_i64_extend_i32_s() *instConvert_i64_extend_i32_s {
	return &instConvert_i64_extend_i32_s{}
}
func (i *instConvert_i64_extend_i32_s) Format(indent string, sb *strings.Builder) {
	sb.WriteString(indent)
	sb.WriteString("i64.extend_i32_s")
}

/**************************************
instConvert_i64_extend_i32_u:
**************************************/
type instConvert_i64_extend_i32_u struct {
	anInstruction
}

func NewInstConvert_i64_extend_i32_u() *instConvert_i64_extend_i32_u {
	return &instConvert_i64_extend_i32_u{}
}
func (i *instConvert_i64_extend_i32_u) Format(indent string, sb *strings.Builder) {
	sb.WriteString(indent)
	sb.WriteString("i64.extend_i32_u")
}

/**************************************
instConvert_i64_trunc_f32_s:
**************************************/
type instConvert_i64_trunc_f32_s struct {
	anInstruction
}

func NewInstConvert_i64_trunc_f32_s() *instConvert_i64_trunc_f32_s {
	return &instConvert_i64_trunc_f32_s{}
}
func (i *instConvert_i64_trunc_f32_s) Format(indent string, sb *strings.Builder) {
	sb.WriteString(indent)
	sb.WriteString("i64.trunc_f32_s")
}

/**************************************
instConvert_i64_trunc_f64_s:
**************************************/
type instConvert_i64_trunc_f64_s struct {
	anInstruction
}

func NewInstConvert_i64_trunc_f64_s() *instConvert_i64_trunc_f64_s {
	return &instConvert_i64_trunc_f64_s{}
}
func (i *instConvert_i64_trunc_f64_s) Format(indent string, sb *strings.Builder) {
	sb.WriteString(indent)
	sb.WriteString("i64.trunc_f64_s")
}

/**************************************
instConvert_f32_convert_i32_s:
**************************************/
type instConvert_f32_convert_i32_s struct {
	anInstruction
}

func NewInstConvert_f32_convert_i32_s() *instConvert_f32_convert_i32_s {
	return &instConvert_f32_convert_i32_s{}
}
func (i *instConvert_f32_convert_i32_s) Format(indent string, sb *strings.Builder) {
	sb.WriteString(indent)
	sb.WriteString("f32.convert_i32_s")
}

/**************************************
instConvert_f32_convert_i32_u:
**************************************/
type instConvert_f32_convert_i32_u struct {
	anInstruction
}

func NewInstConvert_f32_convert_i32_u() *instConvert_f32_convert_i32_u {
	return &instConvert_f32_convert_i32_u{}
}
func (i *instConvert_f32_convert_i32_u) Format(indent string, sb *strings.Builder) {
	sb.WriteString(indent)
	sb.WriteString("f32.convert_i32_u")
}

/**************************************
instConvert_f32_convert_i64_s:
**************************************/
type instConvert_f32_convert_i64_s struct {
	anInstruction
}

func NewInstConvert_f32_convert_i64_s() *instConvert_f32_convert_i64_s {
	return &instConvert_f32_convert_i64_s{}
}
func (i *instConvert_f32_convert_i64_s) Format(indent string, sb *strings.Builder) {
	sb.WriteString(indent)
	sb.WriteString("f32.convert_i64_s")
}

/**************************************
instConvert_f32_convert_i64_u:
**************************************/
type instConvert_f32_convert_i64_u struct {
	anInstruction
}

func NewInstConvert_f32_convert_i64_u() *instConvert_f32_convert_i64_u {
	return &instConvert_f32_convert_i64_u{}
}
func (i *instConvert_f32_convert_i64_u) Format(indent string, sb *strings.Builder) {
	sb.WriteString(indent)
	sb.WriteString("f32.convert_i64_u")
}

/**************************************
instConvert_f32_demote_f64:
**************************************/
type instConvert_f32_demote_f64 struct {
	anInstruction
}

func NewInstConvert_f32_demote_f64() *instConvert_f32_demote_f64 {
	return &instConvert_f32_demote_f64{}
}
func (i *instConvert_f32_demote_f64) Format(indent string, sb *strings.Builder) {
	sb.WriteString(indent)
	sb.WriteString("f32.demote_f64")
}

/**************************************
instConvert_f64_convert_i32_s:
**************************************/
type instConvert_f64_convert_i32_s struct {
	anInstruction
}

func NewInstConvert_f64_convert_i32_s() *instConvert_f64_convert_i32_s {
	return &instConvert_f64_convert_i32_s{}
}
func (i *instConvert_f64_convert_i32_s) Format(indent string, sb *strings.Builder) {
	sb.WriteString(indent)
	sb.WriteString("f64.convert_i32_s")
}

/**************************************
instConvert_f64_convert_i32_u:
**************************************/
type instConvert_f64_convert_i32_u struct {
	anInstruction
}

func NewInstConvert_f64_convert_i32_u() *instConvert_f64_convert_i32_u {
	return &instConvert_f64_convert_i32_u{}
}
func (i *instConvert_f64_convert_i32_u) Format(indent string, sb *strings.Builder) {
	sb.WriteString(indent)
	sb.WriteString("f64.convert_i32_u")
}

/**************************************
instConvert_f64_convert_i64_s:
**************************************/
type instConvert_f64_convert_i64_s struct {
	anInstruction
}

func NewInstConvert_f64_convert_i64_s() *instConvert_f64_convert_i64_s {
	return &instConvert_f64_convert_i64_s{}
}
func (i *instConvert_f64_convert_i64_s) Format(indent string, sb *strings.Builder) {
	sb.WriteString(indent)
	sb.WriteString("f64.convert_i64_s")
}

/**************************************
instConvert_f64_convert_i64_u:
**************************************/
type instConvert_f64_convert_i64_u struct {
	anInstruction
}

func NewInstConvert_f64_convert_i64_u() *instConvert_f64_convert_i64_u {
	return &instConvert_f64_convert_i64_u{}
}
func (i *instConvert_f64_convert_i64_u) Format(indent string, sb *strings.Builder) {
	sb.WriteString(indent)
	sb.WriteString("f64.convert_i64_u")
}

/**************************************
instConvert_f64_promote_f32:
**************************************/
type instConvert_f64_promote_f32 struct {
	anInstruction
}

func NewInstConvert_f64_promote_f32() *instConvert_f64_promote_f32 {
	return &instConvert_f64_promote_f32{}
}
func (i *instConvert_f64_promote_f32) Format(indent string, sb *strings.Builder) {
	sb.WriteString(indent)
	sb.WriteString("f64.promote_f32")
}
