package gorough

import (
	"github.com/NovikovRoman/gorough/internal/tools"
)

type Drawable interface {
	Name() string
	Operations() []operation
	Attributes() Attributes
	Style() Style
	Filler() Filler
}
