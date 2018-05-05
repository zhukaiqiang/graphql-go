package language

import (
	"bufio"
	"bytes"
	"io"

	"github.com/graph-gophers/graphql-go/errors"
)

// The lexer breaks up a series of characters into _tokens_.
// These tokens then get fed to a parser which performs _semantic analysis_.

var eof = rune(0)

// Scanner represents a lexical scanner.
type Scanner struct {
	rd *bufio.Reader
}

// NewScanner returns an initialized instance of scanner.
func NewScanner(rd io.Reader) *Scanner {
	return &Scanner{rd: bufio.NewReader(rd)}
}

// Scan returns the next token and literal value.
func (s *Scanner) Scan() (Token, error) {
	// Read the next rune.
	r := s.read()

	switch {
	case !isSourceCharacter(r):
		return illegalToken(string(r)), errors.InvalidCharacter(r)

	case isWhiteSpace(r):
		// Consume all contiguous whitespace.
		s.unread()
		return s.scanWhiteSpace()

	case isFirstNameCharacter(r):
		s.unread()
		return s.scanName()

	case isFirstNumberCharacter(r):
		s.unread()
		return s.scanNumber()
	}

	switch r {
	case '!':
		return Token{Kind: BANG}, nil
	case '"':
		if s.isBlockString() {
			s.unread()
			return s.scanBlockString()
		}

		s.unread()
		return s.scanString()
	case '#':
		s.unread()
		return s.scanComment()
	case '$':
		return Token{Kind: DOLLAR}, nil
	case '&':
		return Token{Kind: AMP}, nil
	case '(':
		return Token{Kind: PAREN_L}, nil
	case ')':
		return Token{Kind: PAREN_R}, nil
	case '.':
		for i := 0; i < 2; i++ {
			if next := s.read(); next != '.' {
				break
			}
		}

		return Token{Kind: SPREAD}, nil
	case ':':
		return Token{Kind: COLON}, nil
	case '=':
		return Token{Kind: EQUALS}, nil
	case '@':
		return Token{Kind: AT}, nil
	case '[':
		return Token{Kind: BRACKET_L}, nil
	case ']':
		return Token{Kind: BRACKET_R}, nil
	case '{':
		return Token{Kind: BRACE_L}, nil
	case '|':
		return Token{Kind: PIPE}, nil
	case '}':
		return Token{Kind: BRACE_R}, nil
	case eof:
		return Token{Kind: EOF}, nil
	}

	return Token{Kind: ILLEGAL, Value: string(r)}, errors.UnexpectedCharacter(r)
}

// read the next rune from the bufferred reader.
// Returns the rune(0) if an error occurs (or io.EOF is returned).
func (s *Scanner) read() rune {
	r, _, err := s.rd.ReadRune()
	if err != nil {
		return eof
	}
	return r
}

// unread places the previously read rune back on the reader.
func (s *Scanner) unread() { _ = s.rd.UnreadRune() }

func (s *Scanner) peek() rune {
	r := s.read()
	s.unread()
	return r
}

func (s *Scanner) isBlockString() bool {
	// TODO: Fixup the reading here.
	if s.read() != '"' {
		s.unread()
		return false
	}

	if s.read() != '"' {
		s.unread()
		s.unread()
		return false
	}

	return true
}

func (s *Scanner) scanComment() (Token, error) {
	var buf bytes.Buffer

	if s.peek() == '#' {
		// Remove leading '#', if not already removed.
		_ = s.read()
	}

	// Remove leading whitespace.
	tok, err := s.scanWhiteSpace()
	if err != nil {
		return tok, err
	}

	for r := s.read(); r != eof; r = s.read() {
		if r != 0x0009 && r < 0x0020 {
			break
		}

		buf.WriteRune(r)
	}

	return Token{Kind: COMMENT, Value: buf.String()}, nil
}

// scanWhiteSpace consumes the current rune and all contiguous whitespace.
func (s *Scanner) scanWhiteSpace() (Token, error) {
	// Create a buffer and read the current character into it.
	var buf bytes.Buffer
	buf.WriteRune(s.read())

	// Read every subsequent whitespace character into the buffer.
	// Non-whitespace characters and EOF will cause the loop to exit.
loop:
	for r := s.read(); ; r = s.read() {
		switch {
		case r == eof:
			break loop
		case !isWhiteSpace(r):
			s.unread()
			break loop
		default:
			buf.WriteRune(r)
		}
	}

	return Token{Kind: WHITE_SPACE, Value: buf.String()}, nil
}

// scanName consumes the current rune and all contigusous ident runes.
func (s *Scanner) scanName() (Token, error) {
	// TODO: Implement
	// Create a buffer and read the current character into it.
	var buf bytes.Buffer
	r := s.read()

	// TODO: Names can be prefixed by '_', so we must check for that as well.
	if !isLetter(r) {
		return illegalToken(string(r)), errors.InvalidCharacter(r)
	}

	buf.WriteRune(r)

	// Read every subsequent ident character into the buffer.
	// Non-ident characters and EOF will cause the loop to exit.
	for r := s.read(); r != eof; r = s.read() {
		if !isNameCharacter(r) {
			s.unread()
			break
		}

		_, _ = buf.WriteRune(r)
	}

	value := buf.String()

	// If the string matches a keyword then return that keyword.
	switch value {
	case "schema":
	case "query":
	case "mutation":
	case "subscription":
	}

	return Token{Kind: NAME, Value: value}, nil
}

func (s *Scanner) scanNumber() (Token, error) {
	var buf bytes.Buffer
	var kind TokenKind = INT

	if s.peek() == '-' {
		buf.WriteRune(s.read())
	}

	if s.peek() == '0' {
		buf.WriteRune(s.read())
		// Numbers with leading zeros are invalid.
		// Ensure that there are no digits following this 0.
		if isDigit(s.peek()) {
			return Token{}, errors.InvalidNumber(s.peek(), "unexpected digit after 0")
		}
	}

	for r := s.read(); isDigit(r); r = s.read() {
		buf.WriteRune(r)
	}
	s.unread()

	if s.peek() == '.' {
		kind = FLOAT
		buf.WriteRune(s.read())

		for r := s.read(); isDigit(r); r = s.read() {
			buf.WriteRune(r)
		}
		s.unread()
	}

	if r := s.peek(); r == 'e' || r == 'E' {
		buf.WriteRune(s.read())

		if next := s.peek(); next == '-' || next == '+' {
			buf.WriteRune(s.read())
		}

		for r := s.read(); isDigit(r); r = s.read() {
			buf.WriteRune(r)
		}
		s.unread()
	}

	return Token{Kind: kind, Value: buf.String()}, nil
}

func (s *Scanner) scanBlockString() (Token, error) {
	// TODO: Implement.
	return Token{Kind: BLOCK_STRING}, errors.New("scanBlockString: not yet implemented")
}

func (s *Scanner) scanString() (Token, error) {
	var buf bytes.Buffer

	if r := s.read(); r != '"' {
		return Token{}, errors.Errorf(`expected a double quote ("), got %q`, r)
	}

	for r := s.read(); r != '"'; r = s.read() {
		if r == eof {
			return Token{}, errors.Errorf("unexpected EOF")
		}

		if isLineTerminator(r) {
			return Token{}, errors.Errorf("unexpected line terminator")
		}

		if r == '\\' {
			var escaped rune

			switch next := s.read(); next {
			case '"', '\\', '/':
				escaped = next
			case 'b':
				escaped = '\b'
			case 'f':
				escaped = '\f'
			case 'n':
				escaped = '\n'
			case 'r':
				escaped = '\r'
			case 't':
				escaped = '\t'
			case 'u':
				// TODO: handle escaped unicode.
			default:
				return Token{}, errors.InvalidEscapeSequence(r)
			}

			if escaped != eof {
				buf.WriteRune(escaped)
				continue
			}
		}

		buf.WriteRune(r)
	}

	return Token{Kind: STRING, Value: buf.String()}, nil
}

//
func illegalToken(val string) Token {
	return Token{Kind: ILLEGAL, Value: val}
}
