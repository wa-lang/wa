// 版权 @2024 凹语言 作者。保留所有权利。

package source

import "fmt"

func (k FileKind) String() string {
	switch k {
	case Mod:
		return "wa.mod"
	default:
		return "wa"
	}
}

func (id FileIdentity) String() string {
	return fmt.Sprintf("%s%s%s", id.URI, id.Hash, id.Kind)
}

func (a FileAction) String() string {
	switch a {
	case Open:
		return "Open"
	case Change:
		return "Change"
	case Close:
		return "Close"
	case Save:
		return "Save"
	case Create:
		return "Create"
	case Delete:
		return "Delete"
	case InvalidateMetadata:
		return "InvalidateMetadata"
	default:
		return "Unknown"
	}
}
