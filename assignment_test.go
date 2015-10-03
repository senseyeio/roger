package roger

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type Info struct {
	A string `r:"v_a"`
	B string `r:"v_b"`
}

func TestStructToR(t *testing.T) {
	client, _ := NewRClient("localhost", 6311)
	sess, err := client.GetSession()
	defer sess.Close()
	// obj, err := sess.SendCommand("2.2").GetResultObject()
	// assert.Equal(t, obj, float64(2.2))
	// assert.Equal(t, err, nil)

	info := &Info{A: "1, 2, 3, 4", B: `"a" "b" "c" "d"`}
	sess.StructToR(info)
	obj_a, err := sess.SendCommand("v_a[1]").GetResultObject()
	obj_b, err := sess.SendCommand("v_b[4]").GetResultObject()

	assert.Equal(t, obj_a, float64(1))
	assert.Equal(t, obj_b, "d")
	assert.Equal(t, err, nil)
}
