package roger

import "errors"

type authType int

const (
	atPlain authType = 1
	atCrypt authType = 2
)

func login(sess *session) error {
	if sess.authReq == false {
		return nil
	}
	if sess.authReq == true && (sess.user == "" || sess.password == "") {
		return errors.New("Authentication is required and no credentials have been specified")
	}
	if sess.key == "" {
		sess.key = "rs"
	}
	cmd := sess.user + "\n" + sess.password
	if sess.authType == atCrypt {
		cmd = sess.user + "\n" + crypt(sess.password, sess.key)
	}

	packet := sess.sendCommand(cmdLogin, cmd)
	if packet.IsError() {
		_, err := packet.GetResultObject()
		return errors.New("Authentication failed: " + err.Error())
	}
	return nil
}
