package errors

import "fmt"

func InvalidCharacter(r rune) Syntax {
	return Syntax{
		Message: fmt.Sprintf("syntax error: cannot contain the invalid character %q", r),
	}
}

func InvalidEscapeSequence(r rune) Syntax {
	return Syntax{
		Message: fmt.Sprintf("syntax error: invalid character escape sequence: \\%c", r),
	}
}

func InvalidNumber(r rune, msg string) Syntax {
	return Syntax{
		Message: fmt.Sprintf("invalid number: %s: %q", msg, r),
	}
}

func UnexpectedCharacter(r rune) Syntax {
	var msg string

	if r == '\'' {
		msg = `unexpected single quote character ('), did you mean to use a double quote (")?`
	} else {
		msg = fmt.Sprintf("cannot parse the unexpected character %q", r)
	}

	return Syntax{Message: "syntax error: " + msg}
}

type Syntax struct {
	// TODO: Add pointer-to-Source.
	Message  string
	Position Position
}

func (err Syntax) Error() string {
	return err.Message
}

var _ error = Syntax{}

type Position struct {
	Line   int
	Column int
}

// Implements the `fmt.Stringer` interface.
func (p Position) String() string {
	return fmt.Sprintf("(line %d, column %d)", p.Line, p.Column)
}
