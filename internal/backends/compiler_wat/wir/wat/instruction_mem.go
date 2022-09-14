// 版权 @2022 凹语言 作者。保留所有权利。

package wat

import "strconv"

/**************************************
InstLoad:
**************************************/
type InstLoad struct {
	anInstruction
	typ           ValueType
	offset, align int
}

func NewInstLoad(typ ValueType, offset int, align int) *InstLoad {
	return &InstLoad{typ: typ, offset: offset, align: align}
}
func (i *InstLoad) Format(indent string) string {
	return indent + i.typ.Name() + ".load offset=" + strconv.Itoa(i.offset) + " align=" + strconv.Itoa(i.align)
}

/**************************************
InstStore:
**************************************/
type InstStore struct {
	anInstruction
	typ           ValueType
	offset, align int
}

func NewInstStore(typ ValueType, offset int, align int) *InstStore {
	return &InstStore{typ: typ, offset: offset, align: align}
}
func (i *InstStore) Format(indent string) string {
	return indent + i.typ.Name() + ".store offset=" + strconv.Itoa(i.offset) + " align=" + strconv.Itoa(i.align)
}
