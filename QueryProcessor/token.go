package queryprocessor

type TokenType int

const (
	IDENTIFIER TokenType = iota
	SYMBOL
	KEYWORD
	DATA
)

func (t TokenType) String() string {
	switch t {
	case IDENTIFIER:
		return "IDENTIFIER"
	case SYMBOL:
		return "SYMBOL"
	case KEYWORD:
		return "KEYWORD"
	case DATA:
		return "DATA"
	}
	return "UNKNOWN"
}

type Token struct {
	Type  TokenType
	Value string
}
