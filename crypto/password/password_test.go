package password

import "testing"

func TestEncryptAndVerify(t *testing.T) {
	password := "123456"

	encrypted, err := Encrypt(password)
	if err != nil {
		t.Fatal(err)
	}

	if err := Verify(encrypted, "12345"); err == nil {
		t.Fatal("verify wrong password should not pass")
	}

	if err := Verify(encrypted, password); err != nil {
		t.Fatal("verify correct password should pass")
	}
}
