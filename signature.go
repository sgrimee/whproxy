package main

import (
	"crypto/hmac"
	"crypto/sha1"
)

// CheckMAC reports whether messageSig is a valid HMAC tag for message.
func validSignature(message, messageSig, secret []byte) bool {
	mac := hmac.New(sha1.New, secret)
	mac.Write(message)
	expectedMAC := mac.Sum(nil)
	// log.Printf("secret: %s", string(secret))
	// log.Printf("messageSig: %x", string(messageSig))
	// log.Printf("computedSig: %x", string(expectedMAC))
	return hmac.Equal(messageSig, expectedMAC)
}
