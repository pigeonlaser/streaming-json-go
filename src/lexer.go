package main

import (
	"fmt"
	"strings"
)

// token const
const (
	TOKEN_EOF                  = iota // end-of-file
	TOKEN_LEFT_BRACKET                // [
	TOKEN_RIGHT_BRACKET               // ]
	TOKEN_LEFT_BRACE                  // {
	TOKEN_RIGHT_BRACE                 // }
	TOKEN_COLON                       // :
	TOKEN_DOT                         // .
	TOKEN_COMMA                       // ,
	TOKEN_QUOTE                       // "
	TOKEN_ESCAPE_CHARACTER            // \
	TOKEN_NULL                        // null
	TOKEN_TRUE                        // true
	TOKEN_FLASE                       // false
	TOKEN_ALPHABET_LOWERCASE_A        // a
	TOKEN_ALPHABET_LOWERCASE_E        // e
	TOKEN_ALPHABET_LOWERCASE_F        // f
	TOKEN_ALPHABET_LOWERCASE_L        // l
	TOKEN_ALPHABET_LOWERCASE_N        // n
	TOKEN_ALPHABET_LOWERCASE_R        // r
	TOKEN_ALPHABET_LOWERCASE_S        // s
	TOKEN_ALPHABET_LOWERCASE_T        // t
	TOKEN_ALPHABET_LOWERCASE_U        // u
	TOKEN_OTHERS                      // anything else in json
)

// token symbol const
const (
	TOKEN_LEFT_BRACKET_SYMBOL         = '['
	TOKEN_RIGHT_BRACKET_SYMBOL        = ']'
	TOKEN_LEFT_BRACE_SYMBOL           = '{'
	TOKEN_RIGHT_BRACE_SYMBOL          = '}'
	TOKEN_COLON_SYMBOL                = ':'
	TOKEN_DOT_SYMBOL                  = '.'
	TOKEN_COMMA_SYMBOL                = ','
	TOKEN_QUOTE_SYMBOL                = '"'
	TOKEN_ESCAPE_CHARACTER_SYMBOL     = '\\'
	TOKEN_ALPHABET_LOWERCASE_A_SYMBOL = 'a'
	TOKEN_ALPHABET_LOWERCASE_E_SYMBOL = 'e'
	TOKEN_ALPHABET_LOWERCASE_F_SYMBOL = 'f'
	TOKEN_ALPHABET_LOWERCASE_L_SYMBOL = 'l'
	TOKEN_ALPHABET_LOWERCASE_N_SYMBOL = 'n'
	TOKEN_ALPHABET_LOWERCASE_R_SYMBOL = 'r'
	TOKEN_ALPHABET_LOWERCASE_S_SYMBOL = 's'
	TOKEN_ALPHABET_LOWERCASE_T_SYMBOL = 't'
	TOKEN_ALPHABET_LOWERCASE_U_SYMBOL = 'u'
)

var tokenNameMap = map[int]string{
	TOKEN_EOF:                  "EOF",
	TOKEN_LEFT_BRACKET:         "[",
	TOKEN_RIGHT_BRACKET:        "]",
	TOKEN_LEFT_BRACE:           "{",
	TOKEN_RIGHT_BRACE:          "}",
	TOKEN_COLON:                ":",
	TOKEN_DOT:                  ".",
	TOKEN_COMMA:                ",",
	TOKEN_QUOTE:                "\"",
	TOKEN_ESCAPE_CHARACTER:     "\\",
	TOKEN_NULL:                 "null",
	TOKEN_TRUE:                 "true",
	TOKEN_FLASE:                "false",
	TOKEN_ALPHABET_LOWERCASE_A: "a",
	TOKEN_ALPHABET_LOWERCASE_E: "e",
	TOKEN_ALPHABET_LOWERCASE_F: "f",
	TOKEN_ALPHABET_LOWERCASE_L: "l",
	TOKEN_ALPHABET_LOWERCASE_N: "n",
	TOKEN_ALPHABET_LOWERCASE_R: "r",
	TOKEN_ALPHABET_LOWERCASE_S: "s",
	TOKEN_ALPHABET_LOWERCASE_T: "t",
	TOKEN_ALPHABET_LOWERCASE_U: "u",
}

var leftPairTokens = map[int]bool{
	TOKEN_LEFT_BRACKET: true,
	TOKEN_LEFT_BRACE:   true,
}

var rightPairTokens = map[int]bool{
	TOKEN_RIGHT_BRACKET: true,
	TOKEN_RIGHT_BRACE:   true,
}

var mirrorTokenMap = map[int]int{
	TOKEN_LEFT_BRACKET: TOKEN_RIGHT_BRACKET,
	TOKEN_LEFT_BRACE:   TOKEN_RIGHT_BRACE,
	TOKEN_QUOTE:        TOKEN_QUOTE,
}

type Lexer struct {
	JSONContent      strings.Builder
	JSONSegment      string
	TokenStack       []int
	MirrorTokenStack []int
}

func NewLexer() *Lexer {
	return &Lexer{}
}

func (lexer *Lexer) getTopTokenOnStack() int {
	tokenStackLen := len(lexer.TokenStack)
	if tokenStackLen == 0 {
		return TOKEN_EOF
	}
	return lexer.TokenStack[tokenStackLen-1]
}

func (lexer *Lexer) getTopTokenOnMirrorStack() int {
	mirrotTokenStackLen := len(lexer.MirrorTokenStack)
	if mirrotTokenStackLen == 0 {
		return TOKEN_EOF
	}
	return lexer.MirrorTokenStack[mirrotTokenStackLen-1]
}

func (lexer *Lexer) popTokenStack() int {
	tokenStackLen := len(lexer.TokenStack)
	if tokenStackLen == 0 {
		return TOKEN_EOF
	}
	token := lexer.TokenStack[tokenStackLen-1]
	lexer.TokenStack = lexer.TokenStack[:tokenStackLen-1]
	return token
}

func (lexer *Lexer) popMirrorTokenStack() int {
	mirrorTokenStackLen := len(lexer.MirrorTokenStack)
	if mirrorTokenStackLen == 0 {
		return TOKEN_EOF
	}
	token := lexer.MirrorTokenStack[mirrorTokenStackLen-1]
	lexer.MirrorTokenStack = lexer.MirrorTokenStack[:mirrorTokenStackLen-1]
	return token
}

func (lexer *Lexer) pushTokenStack(token int) {
	lexer.TokenStack = append(lexer.TokenStack, token)
}

func (lexer *Lexer) pushMirrorTokenStack(token int) {
	lexer.MirrorTokenStack = append(lexer.MirrorTokenStack, token)
}

func (lexer *Lexer) dumpMirrorTokenStackToString() string {
	var stackInString strings.Builder
	for i := len(lexer.MirrorTokenStack) - 1; i >= 0; i-- {
		stackInString.WriteString(tokenNameMap[lexer.MirrorTokenStack[i]])
	}
	return stackInString.String()
}

func (lexer *Lexer) isLeftPairToken(token int) bool {
	itIs, hit := leftPairTokens[token]
	return itIs && hit
}

func (lexer *Lexer) isRightPairToken(token int) bool {
	itIs, hit := rightPairTokens[token]
	return itIs && hit
}

func (lexer *Lexer) skipJSONSegment(n int) {
	lexer.JSONSegment = lexer.JSONSegment[n:]
}

func (lexer *Lexer) streamStoppedInAnObject() bool {
	fmt.Printf("[DUMP] streamStoppedInAnObject.MirrorTokenStack: '%+v'\n", lexer.MirrorTokenStack)
	// only "}" left
	if lexer.getTopTokenOnMirrorStack() == TOKEN_RIGHT_BRACE {
		return true
	}

	// ":", "null", "}" left
	case1 := []int{
		TOKEN_RIGHT_BRACE,
		TOKEN_ALPHABET_LOWERCASE_L,
		TOKEN_ALPHABET_LOWERCASE_L,
		TOKEN_ALPHABET_LOWERCASE_U,
		TOKEN_ALPHABET_LOWERCASE_N,
		TOKEN_COLON,
	}
	if matchStack(lexer.MirrorTokenStack, case1) {
		return true
	}
	// "null", "}" left
	case2 := []int{
		TOKEN_RIGHT_BRACE,
		TOKEN_ALPHABET_LOWERCASE_L,
		TOKEN_ALPHABET_LOWERCASE_L,
		TOKEN_ALPHABET_LOWERCASE_U,
		TOKEN_ALPHABET_LOWERCASE_N,
	}
	if matchStack(lexer.MirrorTokenStack, case2) {
		return true
	}
	return false
}

// check if JSON stream stopped in an object properity's key, like `{"field`
func (lexer *Lexer) streamStoppedInAnObjectKey() bool {
	tokens := []int{
		TOKEN_RIGHT_BRACE,
		TOKEN_ALPHABET_LOWERCASE_L,
		TOKEN_ALPHABET_LOWERCASE_L,
		TOKEN_ALPHABET_LOWERCASE_U,
		TOKEN_ALPHABET_LOWERCASE_N,
		TOKEN_COLON,
		TOKEN_QUOTE,
	}
	return matchStack(lexer.MirrorTokenStack, tokens)
}

// check if JSON stream stopped in an object properity's value, like `{"field": "value`
func (lexer *Lexer) streamStoppedInAnObjectValue() bool {
	tokens := []int{
		TOKEN_RIGHT_BRACE,
		TOKEN_QUOTE,
	}
	return matchStack(lexer.MirrorTokenStack, tokens)
}

func (lexer *Lexer) streamStoppedInAnArray() bool {
	return lexer.getTopTokenOnMirrorStack() == TOKEN_RIGHT_BRACKET
}

func (lexer *Lexer) streamStoppedInAString() bool {
	return lexer.getTopTokenOnStack() == TOKEN_QUOTE && lexer.getTopTokenOnMirrorStack() == TOKEN_QUOTE
}

func (lexer *Lexer) matchToken() (int, byte) {
	// finish
	fmt.Printf("[DUMP] len(lexer.JSONSegment): %d\n", len(lexer.JSONSegment))
	fmt.Printf("[DUMP] lexer.JSONSegment: '%s'\n", lexer.JSONSegment)
	if len(lexer.JSONSegment) == 0 {
		return TOKEN_EOF, byte(0)
	}
	tokenSymbol := lexer.JSONSegment[0]

	// check token
	switch tokenSymbol {
	case TOKEN_LEFT_BRACKET_SYMBOL:
		lexer.skipJSONSegment(1)
		return TOKEN_LEFT_BRACKET, tokenSymbol
	case TOKEN_RIGHT_BRACKET_SYMBOL:
		lexer.skipJSONSegment(1)
		return TOKEN_RIGHT_BRACKET, tokenSymbol
	case TOKEN_LEFT_BRACE_SYMBOL:
		lexer.skipJSONSegment(1)
		return TOKEN_LEFT_BRACE, tokenSymbol
	case TOKEN_RIGHT_BRACE_SYMBOL:
		lexer.skipJSONSegment(1)
		return TOKEN_RIGHT_BRACE, tokenSymbol
	case TOKEN_COLON_SYMBOL:
		lexer.skipJSONSegment(1)
		return TOKEN_COLON, tokenSymbol
	case TOKEN_DOT_SYMBOL:
		lexer.skipJSONSegment(1)
		return TOKEN_DOT, tokenSymbol
	case TOKEN_COMMA_SYMBOL:
		lexer.skipJSONSegment(1)
		return TOKEN_COMMA, tokenSymbol
	case TOKEN_QUOTE_SYMBOL:
		lexer.skipJSONSegment(1)
		return TOKEN_QUOTE, tokenSymbol
	case TOKEN_ALPHABET_LOWERCASE_A_SYMBOL:
		lexer.skipJSONSegment(1)
		return TOKEN_ALPHABET_LOWERCASE_A, tokenSymbol
	case TOKEN_ALPHABET_LOWERCASE_E_SYMBOL:
		lexer.skipJSONSegment(1)
		return TOKEN_ALPHABET_LOWERCASE_E, tokenSymbol
	case TOKEN_ALPHABET_LOWERCASE_F_SYMBOL:
		lexer.skipJSONSegment(1)
		return TOKEN_ALPHABET_LOWERCASE_F, tokenSymbol
	case TOKEN_ALPHABET_LOWERCASE_L_SYMBOL:
		lexer.skipJSONSegment(1)
		return TOKEN_ALPHABET_LOWERCASE_L, tokenSymbol
	case TOKEN_ALPHABET_LOWERCASE_N_SYMBOL:
		lexer.skipJSONSegment(1)
		return TOKEN_ALPHABET_LOWERCASE_N, tokenSymbol
	case TOKEN_ALPHABET_LOWERCASE_R_SYMBOL:
		lexer.skipJSONSegment(1)
		return TOKEN_ALPHABET_LOWERCASE_R, tokenSymbol
	case TOKEN_ALPHABET_LOWERCASE_S_SYMBOL:
		lexer.skipJSONSegment(1)
		return TOKEN_ALPHABET_LOWERCASE_S, tokenSymbol
	case TOKEN_ALPHABET_LOWERCASE_T_SYMBOL:
		lexer.skipJSONSegment(1)
		return TOKEN_ALPHABET_LOWERCASE_T, tokenSymbol
	case TOKEN_ALPHABET_LOWERCASE_U_SYMBOL:
		lexer.skipJSONSegment(1)
		return TOKEN_ALPHABET_LOWERCASE_U, tokenSymbol
	default:
		lexer.skipJSONSegment(1)
		return TOKEN_OTHERS, tokenSymbol
	}
}

func (lexer *Lexer) AppendString(str string) error {
	lexer.JSONSegment = str
	for {
		token, tokenSymbol := lexer.matchToken()
		fmt.Printf("[DUMP] AppendString.token: %s\n", tokenNameMap[token])

		switch token {
		case TOKEN_EOF:
			// nothing to do with TOKEN_EOF
		case TOKEN_OTHERS:
			lexer.JSONContent.WriteByte(tokenSymbol)
		case TOKEN_QUOTE:
			fmt.Printf("    case TOKEN_QUOTE:\n")
			fmt.Printf("    lexer.streamStoppedInAnObject():%+v\n", lexer.streamStoppedInAnObject())
			fmt.Printf("    lexer.getTopTokenOnMirrorStack():%+v\n", lexer.getTopTokenOnMirrorStack())

			lexer.JSONContent.WriteByte(tokenSymbol)
			lexer.pushTokenStack(token)
			if lexer.streamStoppedInAnObject() {
				fmt.Printf("    lexer.streamStoppedInAnObject()\n")
				// push `null`, `:`, `"` into mirror stack
				lexer.pushMirrorTokenStack(TOKEN_ALPHABET_LOWERCASE_L)
				lexer.pushMirrorTokenStack(TOKEN_ALPHABET_LOWERCASE_L)
				lexer.pushMirrorTokenStack(TOKEN_ALPHABET_LOWERCASE_U)
				lexer.pushMirrorTokenStack(TOKEN_ALPHABET_LOWERCASE_N)
				lexer.pushMirrorTokenStack(TOKEN_COLON)
				lexer.pushMirrorTokenStack(TOKEN_QUOTE)
			} else if lexer.streamStoppedInAnArray() {
				fmt.Printf("    lexer.streamStoppedInAnArray()\n")

				// push `"` into mirror stack
				lexer.pushMirrorTokenStack(TOKEN_QUOTE)
			} else if lexer.streamStoppedInAString() {
				fmt.Printf("    lexer.streamStoppedInAString()\n")

				// check if stopped in key of object's properity or value of object's properity
				if lexer.streamStoppedInAnObjectKey() {
					fmt.Printf("    lexer.streamStoppedInAnObjectKey()\n")

					// pop `"` from mirror stack
					lexer.popMirrorTokenStack()
				} else if lexer.streamStoppedInAnObjectValue() {
					fmt.Printf("    lexer.streamStoppedInAnObjectValue()\n")

					// pop `"` from mirror stack
					lexer.popMirrorTokenStack()
				} else {
					return fmt.Errorf("invalied quote token in json stream, incompleted object properity")
				}
			} else {
				return fmt.Errorf("invalied quote token in json stream")
			}
		case TOKEN_COLON:
			lexer.JSONContent.WriteByte(tokenSymbol)
			lexer.pushTokenStack(token)
			if lexer.streamStoppedInAnObject() {
				// pop `:` from mirror stack
				lexer.popMirrorTokenStack()
			}
		case TOKEN_ALPHABET_LOWERCASE_A:
			lexer.JSONContent.WriteByte(tokenSymbol)
			// in a string, just skip token
			if lexer.streamStoppedInAString() {
				continue
			}
			// check if `f` in token stack and `a`, `l`, `s`, `e in mirror stack
			itIsPartOfTokenFalse := func() bool {

				left := []int{
					TOKEN_ALPHABET_LOWERCASE_F,
				}
				if !matchStack(lexer.TokenStack, left) {
					return false
				}
				right := []int{
					TOKEN_ALPHABET_LOWERCASE_E,
					TOKEN_ALPHABET_LOWERCASE_S,
					TOKEN_ALPHABET_LOWERCASE_L,
					TOKEN_ALPHABET_LOWERCASE_A,
				}
				if !matchStack(lexer.MirrorTokenStack, right) {
					return false
				}
				return true
			}
			if !itIsPartOfTokenFalse() {
				continue
			}
			lexer.pushTokenStack(token)
			lexer.popMirrorTokenStack()
		case TOKEN_ALPHABET_LOWERCASE_E:
			lexer.JSONContent.WriteByte(tokenSymbol)
			// in a string, just skip token
			if lexer.streamStoppedInAString() {
				continue
			}
			// check if `f`, `a`, `l`, `s` in token stack and `e` in mirror stack
			itIsPartOfTokenFalse := func() bool {
				left := []int{
					TOKEN_ALPHABET_LOWERCASE_F,
					TOKEN_ALPHABET_LOWERCASE_A,
					TOKEN_ALPHABET_LOWERCASE_L,
					TOKEN_ALPHABET_LOWERCASE_S,
				}
				if !matchStack(lexer.TokenStack, left) {
					return false
				}
				right := []int{
					TOKEN_ALPHABET_LOWERCASE_E,
				}
				if !matchStack(lexer.MirrorTokenStack, right) {
					return false
				}
				return true
			}
			// check if `t`, `r`, `u` in token stack and `e` in mirror stack
			itIsPartOfTokenTrue := func() bool {
				left := []int{
					TOKEN_ALPHABET_LOWERCASE_T,
					TOKEN_ALPHABET_LOWERCASE_R,
					TOKEN_ALPHABET_LOWERCASE_U,
				}
				if !matchStack(lexer.TokenStack, left) {
					return false
				}
				right := []int{
					TOKEN_ALPHABET_LOWERCASE_E,
				}
				if !matchStack(lexer.MirrorTokenStack, right) {
					return false
				}
				return true
			}
			if !itIsPartOfTokenFalse() && !itIsPartOfTokenTrue() {
				continue
			}
			lexer.pushTokenStack(token)
			lexer.popMirrorTokenStack()
		case TOKEN_ALPHABET_LOWERCASE_F:
			lexer.JSONContent.WriteByte(tokenSymbol)
			// in a string, just skip token
			if lexer.streamStoppedInAString() {
				continue
			}

			lexer.pushTokenStack(token)
			// pop `n`, `u`, `l`, `l`
			lexer.popMirrorTokenStack()
			lexer.popMirrorTokenStack()
			lexer.popMirrorTokenStack()
			lexer.popMirrorTokenStack()
			// push `a`, `l`, `s`, `e`
			lexer.pushMirrorTokenStack(TOKEN_ALPHABET_LOWERCASE_E)
			lexer.pushMirrorTokenStack(TOKEN_ALPHABET_LOWERCASE_S)
			lexer.pushMirrorTokenStack(TOKEN_ALPHABET_LOWERCASE_L)
			lexer.pushMirrorTokenStack(TOKEN_ALPHABET_LOWERCASE_A)
		case TOKEN_ALPHABET_LOWERCASE_L:
			lexer.JSONContent.WriteByte(tokenSymbol)
			// in a string, just skip token
			if lexer.streamStoppedInAString() {
				continue
			}
			// check if `f`, `a` in token stack and, `l`, `s`, `e` in mirror stack
			itIsPartOfTokenFalse := func() bool {
				left := []int{
					TOKEN_ALPHABET_LOWERCASE_F,
					TOKEN_ALPHABET_LOWERCASE_A,
				}
				if !matchStack(lexer.TokenStack, left) {
					return false
				}
				right := []int{
					TOKEN_ALPHABET_LOWERCASE_E,
					TOKEN_ALPHABET_LOWERCASE_S,
					TOKEN_ALPHABET_LOWERCASE_L,
				}
				if !matchStack(lexer.MirrorTokenStack, right) {
					return false
				}
				return true
			}
			// check if `n`, `u` in token stack and `l`, `l` in mirror stack
			itIsPartOfTokenNull1 := func() bool {
				fmt.Printf("[]RUN itIsPartOfTokenNull1() !!!!!!!\n")

				left := []int{
					TOKEN_ALPHABET_LOWERCASE_N,
					TOKEN_ALPHABET_LOWERCASE_U,
				}
				fmt.Printf("    lexer.TokenStack: %+v\n", lexer.TokenStack)
				if !matchStack(lexer.TokenStack, left) {
					fmt.Printf("left does not match !!!!!!!\n")

					return false
				}
				right := []int{
					TOKEN_ALPHABET_LOWERCASE_L,
					TOKEN_ALPHABET_LOWERCASE_L,
				}
				if !matchStack(lexer.MirrorTokenStack, right) {
					fmt.Printf("does not match !!!!!!!\n")
					return false
				}
				fmt.Printf("match !!!!!!!\n")

				return true
			}
			// check if `n`, `u`, `l` in token stack and `l` in mirror stack
			itIsPartOfTokenNull2 := func() bool {
				left := []int{
					TOKEN_ALPHABET_LOWERCASE_N,
					TOKEN_ALPHABET_LOWERCASE_U,
					TOKEN_ALPHABET_LOWERCASE_L,
				}
				if !matchStack(lexer.TokenStack, left) {
					return false
				}
				right := []int{
					TOKEN_ALPHABET_LOWERCASE_L,
				}
				if !matchStack(lexer.MirrorTokenStack, right) {
					return false
				}
				return true
			}
			if !itIsPartOfTokenFalse() && !itIsPartOfTokenNull1() && !itIsPartOfTokenNull2() {
				continue
			}
			lexer.pushTokenStack(token)
			lexer.popMirrorTokenStack()
		case TOKEN_ALPHABET_LOWERCASE_N:
			lexer.JSONContent.WriteByte(tokenSymbol)
			// in a string, just skip token
			if lexer.streamStoppedInAString() {
				continue
			}

			lexer.pushTokenStack(token)
			lexer.popMirrorTokenStack()
		case TOKEN_ALPHABET_LOWERCASE_R:
			lexer.JSONContent.WriteByte(tokenSymbol)
			// in a string, just skip token
			if lexer.streamStoppedInAString() {
				continue
			}
			// check if `t` in token stack and `r`, `u`, `e in mirror stack
			itIsPartOfTokenTrue := func() bool {

				left := []int{
					TOKEN_ALPHABET_LOWERCASE_T,
				}
				if !matchStack(lexer.TokenStack, left) {
					return false
				}
				right := []int{
					TOKEN_ALPHABET_LOWERCASE_E,
					TOKEN_ALPHABET_LOWERCASE_U,
					TOKEN_ALPHABET_LOWERCASE_R,
				}
				if !matchStack(lexer.MirrorTokenStack, right) {
					return false
				}
				return true
			}
			if !itIsPartOfTokenTrue() {
				continue
			}
			lexer.pushTokenStack(token)
			lexer.popMirrorTokenStack()
		case TOKEN_ALPHABET_LOWERCASE_S:
			lexer.JSONContent.WriteByte(tokenSymbol)
			// in a string, just skip token
			if lexer.streamStoppedInAString() {
				continue
			}
			// check if `f`, `a`, `l` in token stack and `s`, `e in mirror stack
			itIsPartOfTokenFalse := func() bool {

				left := []int{
					TOKEN_ALPHABET_LOWERCASE_F,
					TOKEN_ALPHABET_LOWERCASE_A,
					TOKEN_ALPHABET_LOWERCASE_L,
				}
				if !matchStack(lexer.TokenStack, left) {
					return false
				}
				right := []int{
					TOKEN_ALPHABET_LOWERCASE_E,
					TOKEN_ALPHABET_LOWERCASE_S,
				}
				if !matchStack(lexer.MirrorTokenStack, right) {
					return false
				}
				return true
			}
			if !itIsPartOfTokenFalse() {
				continue
			}
			lexer.pushTokenStack(token)
			lexer.popMirrorTokenStack()
		case TOKEN_ALPHABET_LOWERCASE_T:
			lexer.JSONContent.WriteByte(tokenSymbol)
			// in a string, just skip token
			if lexer.streamStoppedInAString() {
				continue
			}

			lexer.pushTokenStack(token)
			// pop `n`, `u`, `l`, `l`
			lexer.popMirrorTokenStack()
			lexer.popMirrorTokenStack()
			lexer.popMirrorTokenStack()
			lexer.popMirrorTokenStack()
			// push `r`, `u`, `e`
			lexer.pushMirrorTokenStack(TOKEN_ALPHABET_LOWERCASE_E)
			lexer.pushMirrorTokenStack(TOKEN_ALPHABET_LOWERCASE_U)
			lexer.pushMirrorTokenStack(TOKEN_ALPHABET_LOWERCASE_R)
		case TOKEN_ALPHABET_LOWERCASE_U:
			lexer.JSONContent.WriteByte(tokenSymbol)
			// in a string, just skip token
			if lexer.streamStoppedInAString() {
				continue
			}
			// check if `t`, `r` in token stack and, `u`, `e` in mirror stack
			itIsPartOfTokenTrue := func() bool {
				left := []int{
					TOKEN_ALPHABET_LOWERCASE_T,
					TOKEN_ALPHABET_LOWERCASE_R,
				}
				if !matchStack(lexer.TokenStack, left) {
					return false
				}
				right := []int{
					TOKEN_ALPHABET_LOWERCASE_E,
					TOKEN_ALPHABET_LOWERCASE_U,
				}
				if !matchStack(lexer.MirrorTokenStack, right) {
					return false
				}
				return true
			}
			// check if `n` in token stack and `u`, `l`, `l` in mirror stack
			itIsPartOfTokenNull := func() bool {
				left := []int{
					TOKEN_ALPHABET_LOWERCASE_N,
				}
				if !matchStack(lexer.TokenStack, left) {
					return false
				}
				right := []int{
					TOKEN_ALPHABET_LOWERCASE_L,
					TOKEN_ALPHABET_LOWERCASE_L,
					TOKEN_ALPHABET_LOWERCASE_U,
				}
				if !matchStack(lexer.MirrorTokenStack, right) {
					return false
				}
				return true
			}
			if !itIsPartOfTokenTrue() && !itIsPartOfTokenNull() {
				continue
			}
			lexer.pushTokenStack(token)
			lexer.popMirrorTokenStack()
		default:
			lexer.JSONContent.WriteByte(tokenSymbol)
			if lexer.isLeftPairToken(token) {
				lexer.pushTokenStack(token)
				lexer.pushMirrorTokenStack(mirrorTokenMap[token])
			} else if lexer.isRightPairToken(token) {
				lexer.pushTokenStack(token)
				lexer.popMirrorTokenStack()
			}
		}

		// check if end
		if token == TOKEN_EOF {
			break
		}
	}
	return nil
}

// complete missing parts for incomplete number, properity of object, null and boolean.
func (lexer *Lexer) completeMissingParts() string {
	// check if "," or ":" symbol on top of lexer.TokenStack
	if lexer.streamStoppedInAnObject() {
		switch lexer.getTopTokenOnStack() {
		case TOKEN_DOT:
			return `0`
		case TOKEN_COLON:
			return `: null`
		}
	}
	return ""
}

func (lexer *Lexer) CompleteJSON() string {
	mirrorTokens := lexer.dumpMirrorTokenStackToString()
	fmt.Printf("[DUMP] mirrorTokens: %s\n", mirrorTokens)
	return lexer.JSONContent.String() + lexer.dumpMirrorTokenStackToString()
}

// {         }
// {"      ":null}
