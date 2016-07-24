package main

import (
	"github.com/pborman/uuid"
)

func NewId() string {
	// return a string containing a random unique id
	return uuid.New()
}
