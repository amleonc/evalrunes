# evalrunes

Simple validation for runes.

To evaluate single runes, you can use the provided functions. If you want to test rune arrays (and therefore, strings),
you can use the following pattern:

```go
const exampleString = "abc123"
var match bool

// true
match = evalrunes.CompareAgainstValidators([]rune(exampleString), evalrunes.IsAlphanumeric)

// false
match = evalrunes.CompareAgainstValidators([]rune(exampleString), evalrunes.IsSpace)
```

The function `CompareAgainstValidators` will return `true` **if, and only if** all runes in the array match at least one
of the validator functions provided.