package language

import (
	"io"
)

// Parser represents a parser.
//
// The parser's job is to make sense of the tokens emitted by the lexer,
// and ensure they are in the right order.
type Parser struct {
	s   *Scanner
	buf tokenBuffer
}

type tokenBuffer struct {
	prev Token // last read token
	n    int   // buffer size (max=1)
}

// NewParser returns an initialized instance of Parser.
func NewParser(rd io.Reader) *Parser {
	return &Parser{s: NewScanner(rd)}
}

// scan returns the next token from the underlying scanner.
// If a token has been unscanned then read that instead.
func (p *Parser) scan() (Token, error) {
	// If we have a token on the buffer, read that instead.
	if p.buf.n != 0 {
		p.buf.n = 0
		return p.buf.prev, nil
	}

	// Otherwise read the next token from the scanner.
	tok, err := p.s.Scan()
	if err != nil {
		return tok, err
	}

	// Save it to the buffer in case we unscan it later.
	p.buf.prev = tok
	return tok, nil
}

// unscan pushes the previously read token back onto the buffer.
func (p *Parser) unscan() { p.buf.n = 1 }

// scanIgnoreWhitespace scans the next non-white space token.
func (p *Parser) scanIgnoreWhitespace() (Token, error) {
	tok, err := p.scan()
	if err != nil {
		return tok, err
	}

	if tok.Kind == WHITE_SPACE {
		return p.scan()
	}

	return tok, nil
}
