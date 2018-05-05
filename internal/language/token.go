package language

type Token struct {
	Kind   TokenKind
	Start  int
	End    int
	Line   int
	Column int
	Value  string
}

// TokenKind is an enumeration of the possible tokens.
type TokenKind int

// An enumeration of different kinds of tokens that the lexer emits.
const (
	ILLEGAL      = iota
	SOF          // <SOF>
	EOF          // <EOF>
	WHITE_SPACE  //
	BANG         // !
	DOLLAR       // $
	AMP          // &
	PAREN_L      // (
	PAREN_R      // )
	SPREAD       // ...
	COLON        // :
	EQUALS       // =
	NEGATIVE     // -
	POSITIVE     // +
	AT           // @
	BRACKET_L    // [
	BRACKET_R    // ]
	BRACE_L      // {
	BRACE_R      // }
	PIPE         // |
	NAME         // Name
	INT          // Int
	FLOAT        // Float
	STRING       // String
	BLOCK_STRING // BlockString
	COMMENT      // Comment
)

func (k TokenKind) String() string {
	name := kindNames[k]
	if name == "" {
		return "ILLEGAL"
	}

	return name
}

var kindNames = map[TokenKind]string{
	SOF:          "<SOF>",
	EOF:          "<EOF>",
	BANG:         "!",
	DOLLAR:       "$",
	AMP:          "&",
	PAREN_L:      "(",
	PAREN_R:      ")",
	SPREAD:       "...",
	COLON:        ":",
	EQUALS:       "=",
	AT:           "@",
	BRACKET_L:    "[",
	BRACKET_R:    "]",
	BRACE_L:      "{",
	PIPE:         "|",
	BRACE_R:      "}",
	NAME:         "Name",
	INT:          "Int",
	FLOAT:        "Float",
	STRING:       "String",
	BLOCK_STRING: "BlockString",
	COMMENT:      "Comment",
}

// The following functions define character classes.

// SourceCharacter ::
//	/[\u0009\u000A\u000D\u0020-\uFFFF]/
func isSourceCharacter(r rune) bool {
	return r >= rune(0x0020) && r <= rune(0xFFFF) ||
		r == rune(0x0009) ||
		r == rune(0x000A) ||
		r == rune(0x000D)
}

var whiteSpaceSet = map[rune]struct{}{
	'\t': {}, // Horizontal Tab (U+0009)
	' ':  {}, // Space (U+0020)
	',':  {}, // Commas are insignificant, treated as white space
}

func isWhiteSpace(r rune) bool {
	_, ok := whiteSpaceSet[r]
	return ok
}

func isLineTerminator(r rune) bool {
	return r == '\n' || r == '\r'
}

var punctuatorSet = map[string]struct{}{
	"!":   {}, // BANG
	"$":   {}, // DOLLAR
	"(":   {}, // PAREN_L
	")":   {}, // PAREN_R
	"...": {}, // SPREAD
	":":   {}, // COLON
	"=":   {}, // EQUALS
	"@":   {}, // AT
	"[":   {}, // BRACKET_L
	"]":   {}, // BRACKET_R
	"{":   {}, // BRACE_L
	"|":   {}, // PIPE
	"}":   {}, // BRACE_R
}

func isPunctuator(s string) bool {
	_, ok := punctuatorSet[s]
	return ok
}

func isName(s string) bool {
	if s == "" {
		return false
	}

	for i, r := range s {
		if i == 0 && !isFirstNameCharacter(r) {
			return false
		}

		if !isNameCharacter(r) {
			return false
		}
	}

	return true
}

func isFirstNameCharacter(r rune) bool {
	return isLetter(r) || r == '_'
}

func isNameCharacter(r rune) bool {
	return isLetter(r) || isDigit(r) || r == '_'
}

var operationTypeSet = map[string]struct{}{
	"query":        {}, // a read-only fetch
	"mutation":     {}, // a write followed by a fetch
	"subscription": {}, // a long-lived request that fetches data in response to source events
}

func isOperationType(s string) bool {
	_, ok := operationTypeSet[s]
	return ok
}

func isFirstNumberCharacter(r rune) bool {
	return r == '-' || isDigit(r)
}

func isDigit(r rune) bool {
	return '0' <= r && r <= '9'
}

func isNonZeroDigit(r rune) bool {
	return r != '0' && isDigit(r)
}

func isSign(r rune) bool {
	return r == '-' || r == '+'
}

func isExponentIndicator(r rune) bool {
	return r == 'e' || r == 'E'
}

func isBooleanValue(s string) bool {
	return s == "true" || s == "false"
}

func isNullValue(s string) bool {
	return s == "null"
}

// EnumValue ::
// 	Name but not `true` or `false` or `null`
func isEnumValue(s string) bool {
	return !isBooleanValue(s) && !isNullValue(s) && isName(s)
}

func isLetter(r rune) bool {
	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z')
}
