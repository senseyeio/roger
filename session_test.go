package roger

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSessionCommand(t *testing.T) {
	client, _ := NewRClient("localhost", 6311)
	sess, err := client.GetSession()
	defer sess.Close()
	obj, err := sess.SendCommand("2.2").GetResultObject()
	assert.Equal(t, obj, float64(2.2))
	assert.Equal(t, err, nil)
}

func TestSessionEval(t *testing.T) {
	client, _ := NewRClient("localhost", 6311)
	sess, err := client.GetSession()
	defer sess.Close()
	obj, err := sess.Eval("2.2")
	assert.Equal(t, obj, float64(2.2))
	assert.Equal(t, err, nil)
}

func TestMultipleSessionCommands(t *testing.T) {
	client, _ := NewRClient("localhost", 6311)
	sess, err := client.GetSession()
	defer sess.Close()
	assert.Equal(t, err, nil)
	sess.SendCommand("x <- 2.2")
	obj, err := sess.SendCommand("x").GetResultObject()
	assert.Equal(t, obj, float64(2.2))
	assert.Equal(t, err, nil)
}

func TestSessionClose(t *testing.T) {
	client, _ := NewRClient("localhost", 6311)
	sess, err := client.GetSession()
	sess.Close()
	_, err = sess.SendCommand("2").GetResultObject()
	assert.NotEqual(t, err, nil)
}
