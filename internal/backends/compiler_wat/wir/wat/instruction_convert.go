package wat

/**************************************
instConvert_i32_wrap_i64:
**************************************/
type instConvert_i32_wrap_i64 struct {
	anInstruction
}

func NewInstConvert_i32_wrap_i64() *instConvert_i32_wrap_i64    { return &instConvert_i32_wrap_i64{} }
func (i *instConvert_i32_wrap_i64) Format(indent string) string { return indent + "i32.wrap/i64" }

/**************************************
instConvert_i32_truncs_f32:
**************************************/
type instConvert_i32_truncs_f32 struct {
	anInstruction
}

func NewInstConvert_i32_truncs_f32() *instConvert_i32_truncs_f32 {
	return &instConvert_i32_truncs_f32{}
}
func (i *instConvert_i32_truncs_f32) Format(indent string) string { return indent + "i32.trunc_s/f32" }

/**************************************
instConvert_i32_truncs_f64:
**************************************/
type instConvert_i32_truncs_f64 struct {
	anInstruction
}

func NewInstConvert_i32_truncs_f64() *instConvert_i32_truncs_f64 {
	return &instConvert_i32_truncs_f64{}
}
func (i *instConvert_i32_truncs_f64) Format(indent string) string { return indent + "i32.trunc_s/f64" }

/**************************************
instConvert_i64_extends_i32:
**************************************/
type instConvert_i64_extends_i32 struct {
	anInstruction
}

func NewInstConvert_i64_extends_i32() *instConvert_i64_extends_i32 {
	return &instConvert_i64_extends_i32{}
}
func (i *instConvert_i64_extends_i32) Format(indent string) string {
	return indent + "i64.extend_s/i32"
}

/**************************************
instConvert_i64_extendu_i32:
**************************************/
type instConvert_i64_extendu_i32 struct {
	anInstruction
}

func NewInstConvert_i64_extendu_i32() *instConvert_i64_extendu_i32 {
	return &instConvert_i64_extendu_i32{}
}
func (i *instConvert_i64_extendu_i32) Format(indent string) string {
	return indent + "i64.extend_u/i32"
}

/**************************************
instConvert_i64_truncs_f32:
**************************************/
type instConvert_i64_truncs_f32 struct {
	anInstruction
}

func NewInstConvert_i64_truncs_f32() *instConvert_i64_truncs_f32 {
	return &instConvert_i64_truncs_f32{}
}
func (i *instConvert_i64_truncs_f32) Format(indent string) string { return indent + "i64.trunc_s/f32" }

/**************************************
instConvert_i64_truncs_f64:
**************************************/
type instConvert_i64_truncs_f64 struct {
	anInstruction
}

func NewInstConvert_i64_truncs_f64() *instConvert_i64_truncs_f64 {
	return &instConvert_i64_truncs_f64{}
}
func (i *instConvert_i64_truncs_f64) Format(indent string) string { return indent + "i64.trunc_s/f64" }

/**************************************
instConvert_f32_converts_i32:
**************************************/
type instConvert_f32_converts_i32 struct {
	anInstruction
}

func NewInstConvert_f32_converts_i32() *instConvert_f32_converts_i32 {
	return &instConvert_f32_converts_i32{}
}
func (i *instConvert_f32_converts_i32) Format(indent string) string {
	return indent + "f32.convert_s/i32"
}

/**************************************
instConvert_f32_convertu_i32:
**************************************/
type instConvert_f32_convertu_i32 struct {
	anInstruction
}

func NewInstConvert_f32_convertu_i32() *instConvert_f32_convertu_i32 {
	return &instConvert_f32_convertu_i32{}
}
func (i *instConvert_f32_convertu_i32) Format(indent string) string {
	return indent + "f32.convert_u/i32"
}

/**************************************
instConvert_f32_converts_i64:
**************************************/
type instConvert_f32_converts_i64 struct {
	anInstruction
}

func NewInstConvert_f32_converts_i64() *instConvert_f32_converts_i64 {
	return &instConvert_f32_converts_i64{}
}
func (i *instConvert_f32_converts_i64) Format(indent string) string {
	return indent + "f32.convert_s/i64"
}

/**************************************
instConvert_f32_convertu_i64:
**************************************/
type instConvert_f32_convertu_i64 struct {
	anInstruction
}

func NewInstConvert_f32_convertu_i64() *instConvert_f32_convertu_i64 {
	return &instConvert_f32_convertu_i64{}
}
func (i *instConvert_f32_convertu_i64) Format(indent string) string {
	return indent + "f32.convert_u/i64"
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
func (i *instConvert_f32_demote_f64) Format(indent string) string {
	return indent + "f32.demote/64"
}

/**************************************
instConvert_f64_converts_i32:
**************************************/
type instConvert_f64_converts_i32 struct {
	anInstruction
}

func NewInstConvert_f64_converts_i32() *instConvert_f64_converts_i32 {
	return &instConvert_f64_converts_i32{}
}
func (i *instConvert_f64_converts_i32) Format(indent string) string {
	return indent + "f64.convert_s/i32"
}

/**************************************
instConvert_f64_convertu_i32:
**************************************/
type instConvert_f64_convertu_i32 struct {
	anInstruction
}

func NewInstConvert_f64_convertu_i32() *instConvert_f64_convertu_i32 {
	return &instConvert_f64_convertu_i32{}
}
func (i *instConvert_f64_convertu_i32) Format(indent string) string {
	return indent + "f64.convert_u/i32"
}

/**************************************
instConvert_f64_converts_i64:
**************************************/
type instConvert_f64_converts_i64 struct {
	anInstruction
}

func NewInstConvert_f64_converts_i64() *instConvert_f64_converts_i64 {
	return &instConvert_f64_converts_i64{}
}
func (i *instConvert_f64_converts_i64) Format(indent string) string {
	return indent + "f64.convert_s/i64"
}

/**************************************
instConvert_f64_convertu_i64:
**************************************/
type instConvert_f64_convertu_i64 struct {
	anInstruction
}

func NewInstConvert_f64_convertu_i64() *instConvert_f64_convertu_i64 {
	return &instConvert_f64_convertu_i64{}
}
func (i *instConvert_f64_convertu_i64) Format(indent string) string {
	return indent + "f64.convert_u/i64"
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
func (i *instConvert_f64_promote_f32) Format(indent string) string {
	return indent + "f64.promote/f32"
}
