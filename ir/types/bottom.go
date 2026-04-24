package types

import "strings"

// Global
var BottomType = &bottomType{}

type bottomType struct{}

func (b *bottomType) Simple() bool {
	return true
}

func (b *bottomType) Constant() bool {
	return false
}

func (b *bottomType) ToString(sb *strings.Builder) {
	sb.WriteString("Bottom")
}
