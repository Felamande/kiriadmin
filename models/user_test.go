package models

import "testing"

func TestUserEncrypt(t *testing.T) {
	enb, _ := RsaEncrypt([]byte("Hello"))
	t.Log(string(enb))
}
