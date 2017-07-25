package roger

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func checkAssignment(t *testing.T, assignmentObj interface{}) {
	client, _ := NewRClient("localhost", 6311)
	sess, _ := client.GetSession()
	defer sess.Close()
	err := sess.Assign("assignedVar", assignmentObj)
	assert.Equal(t, err, nil)
	obj, err := sess.Eval("assignedVar")
	assert.Equal(t, assignmentObj, obj)
	assert.Equal(t, nil, err)
}

func TestIntArrayAssignment(t *testing.T) {
	checkAssignment(t, []int32{1, 2, 3, 4, 5})
}

func TestIntAssignment(t *testing.T) {
	checkAssignment(t, int32(100))
}

func TestDoubleAssignment(t *testing.T) {
	checkAssignment(t, float64(123.4))
}

func TestFloatArrayAssignment(t *testing.T) {
	checkAssignment(t, []float64{1.1, 2.2, 3.3, 4.4, 5.5})
}

func TestStringArrayAssignment(t *testing.T) {
	checkAssignment(t, []string{"test", "str"})
}

func TestByteArrayAssignment(t *testing.T) {
	checkAssignment(t, []byte{'g', 'o', 'l', 'a', 'n', 'g'})
}

func TestStringAssignment(t *testing.T) {
	checkAssignment(t, "testing")
}

func TestLargeStringAssignment(t *testing.T) {
	checkAssignment(t, strings.Repeat("a", 20000000))
}
