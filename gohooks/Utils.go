package gohooks

// PassInterface returns an byte array from an interface.
func PassInterface(v interface{}) []byte {
	b, _ := v.([]byte)
	return b
}
