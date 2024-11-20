package ui

import (
	"math/rand"
	"strconv"
)

func UniqueKey(id string) string {
	return id + "-" + strconv.Itoa(rand.Intn(100000))
}
