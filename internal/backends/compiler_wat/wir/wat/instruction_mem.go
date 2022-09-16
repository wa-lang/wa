// 版权 @2022 凹语言 作者。保留所有权利。

package wat

import "strconv"

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
func (i *instLoad) Format(indent string) string {
	return indent + i.typ.Name() + ".load offset=" + strconv.Itoa(i.offset) + " align=" + strconv.Itoa(i.align)
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
func (i *instStore) Format(indent string) string {
	return indent + i.typ.Name() + ".store offset=" + strconv.Itoa(i.offset) + " align=" + strconv.Itoa(i.align)
}
