package mocks

import (
	"context"
)

// Cursor is an interface that matches mongo.Cursor's required methods
type Cursor interface {
	Close(context.Context) error
	Next(context.Context) bool
	Decode(interface{}) error
	All(context.Context, interface{}) error
	Err() error
}
