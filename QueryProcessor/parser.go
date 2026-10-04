package queryprocessor

import (
	"strings"
)

var keywords = map[string]bool{
	"CREATE":   true,
	"DATABASE": true,
	"TABLE":    true,
	"SELECT":   true,
	"FROM":     true,
	"WHERE":    true,
	"INSERT":   true,
	"INTO":     true,
	"VALUES":   true,
	"UPDATE":   true,
	"SET":      true,
	"DELETE":   true,
	"INT":      true,
	"STRING":   true,
	"FLOAT":    true,
	"BOOLEAN":  true,
}

func isSymbol(c byte) bool {
	return c == '(' || c == ')' || c == ',' || c == ';' || c == '=' || c == '*'
}

func isData(word string) bool {
	if len(word) >= 2 && word[0] == '"' && word[len(word)-1] == '"' {
		return true
	}
	if len(word) > 0 && ('0' <= word[0] && word[0] <= '9') {
		return true
	}
	return false
}

func classify(word string) Token {
	if keywords[strings.ToUpper(word)] {
		return Token{Type: KEYWORD, Value: word}
	}
	if isData(word) {
		return Token{Type: DATA, Value: word}
	}
	return Token{Type: IDENTIFIER, Value: word}
}

func Tokenize(input string) []Token {
	tokens := []Token{}
	accum := ""
	inQuote := false

	matchAccum := func() {
		if accum != "" {
			tokens = append(tokens, classify(accum))
			accum = ""
		}
	}

	for i := 0; i < len(input); i++ {
		c := input[i]
		if c == '"' {
			inQuote = !inQuote
			accum += string(c)
		} else if inQuote {
			accum += string(c)
		} else if c == ' ' {
			matchAccum()
		} else if isSymbol(c) {
			matchAccum()
			tokens = append(tokens, Token{Type: SYMBOL, Value: string(c)})
		} else {
			accum += string(c)
		}
	}
	matchAccum()
	return tokens
}
