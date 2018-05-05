package language

import "testing"

func TestIsSourceCharacter(t *testing.T) {
	specials := []rune{rune(0x0009), rune(0x000A), rune(0x000D)}

	for _, r := range specials {
		if !isSourceCharacter(rune(r)) {
			t.Errorf("expected %q to be a source character", rune(r))
		}
	}

	for r := 0x0020; r <= 0xFFFF; r++ {
		if !isSourceCharacter(rune(r)) {
			t.Errorf("expected %q to be a source character", rune(r))
		}
	}
}

func TestIsWhiteSpace(t *testing.T) {
	runes := []rune{' ', '\t', ','}

	for _, r := range runes {
		if !isWhiteSpace(r) {
			t.Errorf("expected %q to be a white space character", rune(r))
		}
	}
}

func TestIsLineTerminator(t *testing.T) {
	runes := []rune{'\n', '\r'}

	for _, r := range runes {
		if !isLineTerminator(r) {
			t.Errorf("expected %q to be a line terminator character", rune(r))
		}
	}
}

func TestIsDigit(t *testing.T) {
	digits := []rune{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9'}

	for _, r := range digits {
		if !isDigit(r) {
			t.Errorf("expected %q to be a digit character (value %d)", r, r)
		}
	}
}
