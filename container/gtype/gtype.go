// Package gtype provides kinds of high performance and concurrent-safe basic variable types.
package gtype

// Type is alias of Interface.
type Type = Interface

// New is alias of NewInterface.
// See NewInterface.
func New(value ...interface{}) *Type {
	return NewInterface(value...)
}
