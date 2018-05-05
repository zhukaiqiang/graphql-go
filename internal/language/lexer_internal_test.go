package language

import (
	"errors"
	"strings"
	"testing"
)

type testCase struct {
	name  string
	input string
	token Token
	err   error
}

func TestScanComment(t *testing.T) {
	tests := []testCase{{
		name:  "stops on EOF",
		input: "# a comment",
		token: Token{Kind: COMMENT, Value: "a comment"},
	}, {
		name:  "stops on \\n",
		input: "# a comment\nextra",
		token: Token{Kind: COMMENT, Value: "a comment"},
	}, {
		name:  "stops on \\r",
		input: "# a comment\rextra",
		token: Token{Kind: COMMENT, Value: "a comment"},
	}, {
		name:  "strips leading whitespace",
		input: "#   a comment",
		token: Token{Kind: COMMENT, Value: "a comment"},
	}, {
		name:  "handles empty comment",
		input: "# ",
		token: Token{Kind: COMMENT},
	}}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			sc := NewScanner(strings.NewReader(test.input))

			tok, err := sc.scanComment()

			expectEqualErrors(t, test.err, err)
			expectEqualTokens(t, test.token, tok)
		})
	}
}

func TestScanNumber(t *testing.T) {
	tests := []testCase{{
		name:  "stops after non-digit",
		input: "12345abc123",
		token: Token{Kind: INT, Value: "12345"},
	}, {
		name:  "handles 0 integer",
		input: "0",
		token: Token{Kind: INT, Value: "0"},
	}, {
		name:  "disallows leading 0s",
		input: "01",
		err:   errors.New("invalid number: unexpected digit after 0: '1'"),
	}, {
		name:  "handles negative ints",
		input: "-12345",
		token: Token{Kind: INT, Value: "-12345"},
	}, {
		name:  "handles 0 float",
		input: "0.0",
		token: Token{Kind: FLOAT, Value: "0.0"},
	}, {
		name:  "handles negative floats",
		input: "-1.2345",
		token: Token{Kind: FLOAT, Value: "-1.2345"},
	}, {
		name:  "handles positive exponents",
		input: "-1.2345e+123",
		token: Token{Kind: FLOAT, Value: "-1.2345e+123"},
	}, {
		name:  "handles negative exponents",
		input: "-1.2345E-123",
		token: Token{Kind: FLOAT, Value: "-1.2345E-123"},
	}, {
		name:  "exponent signs are optional",
		input: "1.2e1",
		token: Token{Kind: FLOAT, Value: "1.2e1"},
	}}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			sc := NewScanner(strings.NewReader(test.input))

			tok, err := sc.scanNumber()

			expectEqualErrors(t, test.err, err)
			expectEqualTokens(t, test.token, tok)
		})
	}
}

func TestScanString(t *testing.T) {
	tests := []testCase{{
		name:  "simple string",
		input: `"a simple string"`,
		token: Token{Kind: STRING, Value: "a simple string"},
	}, {
		name:  "respects whitespaces within strings",
		input: `" white   space "`,
		token: Token{Kind: STRING, Value: " white   space "},
	}, {
		name:  "handles empty strings",
		input: `""`,
		token: Token{Kind: STRING},
	}, {
		name:  "handles escaped double quote",
		input: `"quote \""`,
		token: Token{Kind: STRING, Value: `quote "`},
	}, {
		name:  "handles escaped characters",
		input: `"escaped \n\r\b\t\f"`,
		token: Token{Kind: STRING, Value: "escaped \n\r\b\t\f"},
	}, {
		name:  "handles slashes",
		input: `"slashes \\ \/"`,
		token: Token{Kind: STRING, Value: "slashes \\ /"},
	}, {
		name:  "handles escaped unicode",
		input: `"unicode \u1234\u5678\u90AB\uCDEF"`,
		token: Token{Kind: STRING, Value: "unicode \u1234\u5678\u90AB\uCDEF"},
	}}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			sc := NewScanner(strings.NewReader(test.input))

			tok, err := sc.scanString()

			expectEqualErrors(t, test.err, err)
			expectEqualTokens(t, test.token, tok)
		})
	}
}

func TestScanBlockString(t *testing.T) {
	tests := []testCase{{
		name:  "simple block string",
		input: `"""simple block string"""`,
		token: Token{Kind: BLOCK_STRING, Value: "simple block string"},
	}}

	t.Skip("not yet ready to test this")

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			sc := NewScanner(strings.NewReader(test.input))

			tok, err := sc.scanBlockString()

			expectEqualErrors(t, test.err, err)
			expectEqualTokens(t, test.token, tok)
		})
	}
}

func expectEqualTokens(t *testing.T, expected, actual Token) {
	t.Helper()

	if expected.Value != actual.Value {
		t.Errorf("wrong value:\nwant: %q\ngot : %q", expected.Value, actual.Value)
	}

	if expected.Kind != actual.Kind {
		t.Errorf("wrong kind:\nwant: %q\ngot : %q", expected.Kind, actual.Kind)
	}
}

func expectEqualErrors(t *testing.T, expected, actual error) {
	t.Helper()

	switch {
	case expected == actual: // ok
	case expected == nil && actual != nil:
		t.Fatalf("unexpected error:\nwant: %v\ngot : %q", expected, actual)
	case expected != nil && actual == nil:
		t.Fatalf("missing expected error:\nwant: %q\ngot : %v", expected, actual)
	case expected.Error() != actual.Error():
		t.Fatalf("wrong error message:\nwant: %q\ngot : %q", expected, actual)
	}
}
