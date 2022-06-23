package qywx

import (
	"testing"
	"time"
)

func TestFreshTokenTask(t *testing.T) {

	FreshTokenTask("wx861828161a3f015c", "q7EFHBUKk-S1pNBWD0pDXuYjDzahLZ2VaxQ7QfBrYeU", time.Second*15)
}
