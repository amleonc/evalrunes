package evalrunes

import (
	"context"
	"sync"
	"unicode"
	"unicode/utf8"
)

// IsPunctuation checks if the given rune is a punctuation symbol as defined
// in the unicode specification, it returns true if it is, and false otherwise.
func IsPunctuation(c rune) bool {
	return unicode.In(c, unicode.Punct)
}

// IsEmoji checks if the given rune is in the range defined for emojis in the
// unicode specification, it returns true if is, and false otherwise.
func IsEmoji(c rune) bool {
	// Check if the character is in the Emoji ranges defined by the Unicode standard.
	// These ranges are defined in the Unicode Standard Annex #44 (https://www.unicode.org/reports/tr44/).
	return (0x1F004 <= c && c <= 0x1F9C0) || // Miscellaneous Symbols and Pictographs
		(0x1F600 <= c && c <= 0x1F64F) || // Emoticons
		(0x1F680 <= c && c <= 0x1F6FF) || // Transport and Map Symbols
		(0x1F900 <= c && c <= 0x1F9FF) || // Supplemental Symbols and Pictographs
		(0x2600 <= c && c <= 0x27BF) // Miscellaneous Symbols
}

// IsAlphanumeric checks if the given rune is an alphanumeric character, it
// returns true if it is, and false otherwise. It works for most human scripts,
// following the unicode specification.
func IsAlphanumeric(character rune) bool {
	return unicode.IsLetter(character) || IsNumber(character)
}

// IsSpace is a wrapper of the [unicode.IsSpace] function to check
// if a given rune is a space character as defined by the unicode specification.
func IsSpace(r rune) bool {
	return unicode.IsSpace(r)
}

// IsNumber is a wrapper of the [unicode.IsNumber] function.
func IsNumber(r rune) bool {
	return unicode.IsNumber(r)
}

func ValidateMaxRuneArrayLength(runes []rune, max int) bool {
	l := utf8.RuneCountInString(string(runes))
	return l <= max
}

func ValidateMinRuneArrayLength(runes []rune, min int) bool {
	l := utf8.RuneCountInString(string(runes))
	return min <= l
}

// ValidateRuneArrayLength checks if the length of the given rune array
// falls within the specified minimum and maximum values.
// The range is inclusive.
func ValidateRuneArrayLength(runes []rune, min, max int) bool {
	length := utf8.RuneCountInString(string(runes))
	return length >= min && length <= max
}

// CompareAgainstValidators checks whether each character in the provided slice of runes
// matches any of the provided validator functions.

func CompareAgainstValidators(runes []rune, validators ...func(rune) bool) bool {
	if len(validators) < 1 {
		return true
	}

	check := func(ctx context.Context, r rune, result chan bool) {
		for _, fn := range validators {
			if fn(r) {
				return
			}
		}
		select {
		case <-ctx.Done():
		default:
			result <- false
		}
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	results := make(chan bool, len(runes))

	var wg sync.WaitGroup
	wg.Add(len(runes))

	for _, r := range runes {
		go func(r rune) {
			defer wg.Done()
			check(ctx, r, results)
		}(r)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	for match := range results {
		if !match {
			cancel()
			return match
		}
	}

	return true
}
