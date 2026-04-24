package types

import "strings"

type Type interface {
	Simple() bool
	Constant() bool
	ToString(sb *strings.Builder)
}
