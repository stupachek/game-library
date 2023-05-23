package jwt

import (
	"fmt"
	"testing"
)

func TestXxx(t *testing.T) {
	token, err := NewJWT("sdf")
	fmt.Printf("%v, %v", token, err)
}
