package sql

import (
	"fmt"
	"testing"
)

func TestFields(t *testing.T) {
	s := []string{"model_1.field_a", "model_1.field_b", "local_fielda"}
	r := BuildFields(BaseFields, s)
	fmt.Printf("%+v\n", r)
	fmt.Printf("%T\n", r)
}
