package wow

type realm struct {
	realmType  byte
	locked     byte
	flag       byte
	Name       string `json:"name"`
	Address    string `json:"address"`
	population []byte
	numChars   byte
	category   byte
	realmId    byte
}
