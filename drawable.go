package gorough

type Drawable interface {
	Name() string
	Operations() []operation
	Attributes() Attributes
	Styles() *Styles
}

type operation struct {
	code     string
	commands []command
}

type command struct {
	code string
	data []float64
}
