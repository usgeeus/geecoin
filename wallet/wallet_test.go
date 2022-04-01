package wallet

import (
	"crypto/x509"
	"encoding/hex"
	"io/fs"
	"reflect"
	"testing"
)

const (
	testKey     string = "307702010104208c1d3754297790c3ba42b0e2fcb78f07dffdbd825d8179e0af43d04d8f9b2f77a00a06082a8648ce3d030107a1440342000462cd69f3436de111478a96a02236075727437f96079e5e9163c2a7af06059b75237f1a365c274d864607fc1303a812caaf23b5463012b223d5e45e962b5c7c2c"
	testPayload string = "08f6cef73f056b13148f3e8e1d396cfc5f0d9229af21a59f6ef9f793ba302b22"
	testSig     string = "3319e92e29bf32334e2e1cf668dce1439c2cb5e47f4785909e2543013d65895d17f2e81b7af0914a564c2ab86b7bb9ee4b6e34700ae737a2b14f73a238637782"
)

type fakeLayer struct {
	fakeHasWalletFile func() bool
}

func (fakeLayer) writeFile(name string, data []byte, perm fs.FileMode) error {
	return nil
}
func (f fakeLayer) hasWalletFile() bool {
	return f.fakeHasWalletFile()
}
func (fakeLayer) readFile(name string) ([]byte, error) {
	return x509.MarshalECPrivateKey((makeTestWallet().privateKey))
}

func TestWallet(t *testing.T) {
	t.Run("New Wallet is created", func(t *testing.T) {
		files = fakeLayer{
			fakeHasWalletFile: func() bool { return false },
		}

		tw := Wallet()
		if reflect.TypeOf(tw) != reflect.TypeOf(&wallet{}) {
			t.Error("New Wallet should return a new wallet instance")
		}
	})
	t.Run("Wallet is restored", func(t *testing.T) {
		files = fakeLayer{
			fakeHasWalletFile: func() bool { return true },
		}
		w = nil
		tw := Wallet()
		if reflect.TypeOf(tw) != reflect.TypeOf(&wallet{}) {
			t.Error("New Wallet should return a new wallet instance")
		}
	})
}

func makeTestWallet() *wallet {
	w := &wallet{}
	b, _ := hex.DecodeString(testKey)
	key, _ := x509.ParseECPrivateKey(b)
	w.privateKey = key
	w.Address = aFromK(key)
	return w
}

func TestSign(t *testing.T) {
	s := Sign(testPayload, makeTestWallet())
	_, err := hex.DecodeString(s)
	if err != nil {
		t.Errorf("Sign() should return a hex encoded string, got %s", s)
	}
}

func TestVerify(t *testing.T) {
	type test struct {
		input string
		ok    bool
	}
	tests := []test{
		{testPayload, true},
		{"18f6cef73f056b13148f3e8e1d396cfc5f0d9229af21a59f6ef9f793ba302b22", false},
	}
	for _, tc := range tests {
		w := makeTestWallet()
		ok := Verify(testSig, tc.input, w.Address)
		if ok != tc.ok {
			t.Error("Verify() could not verify testSignature and Payload")
		}
	}
}

func TestRestoreBigInts(t *testing.T) {
	_, _, err := restoreBigInts("xx")
	if err == nil {
		t.Error("restoreBig Ints should return error when payload is not hex.")
	}
}
