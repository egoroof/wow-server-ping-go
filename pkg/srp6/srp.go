package srp6

import (
	"crypto/rand"
	"crypto/sha1"
	"math/big"
	"strings"
)

// based on https://github.com/Kangaroux/go-wow-srp6
// and https://gtker.com/implementation-guide-for-the-world-of-warcraft-flavor-of-srp6/

// generates Client Private Key
func ClientPrivateKey() []byte {
	clientPrivateKey := make([]byte, 32)
	rand.Read(clientPrivateKey)
	return clientPrivateKey
}

// calculates Client Public Key
func ClientPublicKey(clientPrivateKey []byte) []byte {
	publicKey := big.NewInt(0).Exp(g, bytesToInt(clientPrivateKey), n)

	return intToBytes(32, publicKey)
}

// returns a 40 byte key that will be used for header encryption/decryption
func ClientSessionKey(
	username, password string,
	salt, clientPublicKey, clientPrivateKey, serverPublicKey []byte,
) []byte {
	x := calculateX(username, password, salt)
	u := calculateU(clientPublicKey, serverPublicKey)
	s := calculateClientSKey(clientPrivateKey, serverPublicKey, x, u)
	return calculateInterleave(s)
}

// returns a proof that the client should send after receiving the auth challenge.
// The server should compare this with the proof received by the client and verify they match.
// If they match, the client has proven they know the session key.
func ClientProof(
	username string,
	salt,
	clientPublicKey,
	serverPublicKey,
	sessionKey []byte,
) []byte {
	hUsername := sha1.Sum([]byte(strings.ToUpper(username)))
	h := sha1.New()
	h.Write(xorHash)
	h.Write(hUsername[:])
	h.Write(salt)
	h.Write(clientPublicKey)
	h.Write(serverPublicKey)
	h.Write(sessionKey)
	return h.Sum(nil)
}

// returns an intermediate value X used for generating the password verifier
func calculateX(username, password string, salt []byte) []byte {
	h := sha1.New()
	inner := sha1.Sum([]byte(strings.ToUpper(username) + ":" + strings.ToUpper(password)))
	h.Write(salt)
	h.Write(inner[:])
	return h.Sum(nil)
}

// returns an intermediate value U used for generating the session key
func calculateU(clientPublicKey, serverPublicKey []byte) []byte {
	h := sha1.New()
	h.Write(clientPublicKey)
	h.Write(serverPublicKey)
	return h.Sum(nil)
}

// returns an intermediate 32 byte S key used to generate the session key
func calculateClientSKey(clientPrivateKey, serverPublicKey, x, u []byte) []byte {
	s := big.NewInt(0).Exp(g, bytesToInt(x), n)
	s.Mul(s, k)
	s.Sub(bytesToInt(serverPublicKey), s)
	inner := big.NewInt(0).Mul(bytesToInt(u), bytesToInt(x))
	inner.Add(inner, bytesToInt(clientPrivateKey))
	s.Exp(s, inner, n)
	return intToBytes(32, s)
}

// returns a 40 byte array containing an interleaved S-key
func calculateInterleave(S []byte) []byte {
	// If the leading byte is zero, remove the leading TWO bytes
	for len(S) > 0 && S[0] == 0 {
		S = S[2:]
	}

	lenS := len(S)
	even, odd := make([]byte, lenS/2), make([]byte, lenS/2)

	// Split the even/odd bytes into separate arrays
	for i := 0; i < lenS/2; i++ {
		even[i] = S[i*2]
		odd[i] = S[i*2+1]
	}

	hEven := sha1.Sum(even)
	hOdd := sha1.Sum(odd)
	interleaved := make([]byte, 40)

	// Interleave the even bytes and odd bytes together, alternating each byte
	for i := range 20 {
		interleaved[i*2] = hEven[i]
		interleaved[i*2+1] = hOdd[i]
	}

	return interleaved
}
