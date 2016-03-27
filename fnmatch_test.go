package fnmatch

import (
	"gospec"
	"testing"
)

func TestAllSpecs(t *testing.T) {
	r := gospec.NewRunner()
	r.AddSpec("MatchSpec", MatchSpec)
	gospec.MainGoTest(r, t)
}

// This is a set of tests ported from a set of tests for C fnmatch
// found at http://www.mail-archive.com/bug-gnulib@gnu.org/msg14048.html
func TestMatch(t *testing.T) {
	assert := func(p, s string) {
		if !Match(p, s, 0) {
			t.Errorf("Assertion failed: Match(%#v, %#v, 0)", p, s)
		}
	}
	assert("", "")
	assert("*", "")
	assert("*", "foo")
	assert("*", "bar")
	assert("*", "*")
	assert("**", "f")
	assert("**", "foo.txt")
	assert("*.*", "foo.txt")
	assert("foo*.txt", "foobar.txt")
	assert("foo.txt", "foo.txt")
	assert("foo\\.txt", "foo.txt")
	if Match("foo\\.txt", "foo.txt", FNM_NOESCAPE) {
		t.Errorf("Assertion failed: Match(%#v, %#v, FNM_NOESCAPE) == false", "foo\\.txt", "foo.txt")
	}
}

func MatchSpec(c *gospec.Context) {
	c.Specify("A wildcard pattern \"*\" should match anything", func() {
		p := "*"
		c.Then(Match(p, "", 0)).Should.Equal(true)
		c.Then(Match(p, "foo", 0)).Should.Equal(true)
		c.Then(Match(p, "*", 0)).Should.Equal(true)
		c.Then(Match(p, "   ", 0)).Should.Equal(true)
		c.Then(Match(p, ".foo", 0)).Should.Equal(true)
		c.Then(Match(p, "わたし", 0)).Should.Equal(true)
		testSlash := func(flags int, result bool) {
			c.Then(Match(p, "foo/bar", flags)).Should.Equal(result)
			c.Then(Match(p, "/", flags)).Should.Equal(result)
			c.Then(Match(p, "/foo", flags)).Should.Equal(result)
			c.Then(Match(p, "foo/", flags)).Should.Equal(result)
		}
		c.Specify("Including / when flags is 0", func() { testSlash(0, true) })
		c.Specify("Except / when used with FNM_PATHNAME", func() { testSlash(FNM_PATHNAME, false) })
		c.Specify("Except . in some circumstances when used with FNM_PERIOD", func() {
			c.Then(Match("*", ".foo", FNM_PERIOD)).Should.Equal(false)
			c.Then(Match("/*", "/.foo", FNM_PERIOD)).Should.Equal(true)
			c.Then(Match("/*", "/.foo", FNM_PERIOD|FNM_PATHNAME)).Should.Equal(false)
		})
	})

	c.Specify("A question mark pattern \"?\" should match a single character", func() {
		p := "?"
		c.Then(Match(p, "", 0)).Should.Equal(false)
		c.Then(Match(p, "f", 0)).Should.Equal(true)
		c.Then(Match(p, ".", 0)).Should.Equal(true)
		c.Then(Match(p, "?", 0)).Should.Equal(true)
		c.Then(Match(p, "foo", 0)).Should.Equal(false)
		c.Then(Match(p, "わ", 0)).Should.Equal(true)
		c.Then(Match(p, "わた", 0)).Should.Equal(false)
		c.Specify("Even / when flags is 0", func() { c.Then(Match(p, "/", 0)).Should.Equal(true) })
		c.Specify("Except / when flags is FNM_PATHNAME", func() { c.Then(Match(p, "/", FNM_PATHNAME)).Should.Equal(false) })
		c.Specify("Except . sometimes when FNM_PERIOD is given", func() {
			c.Then(Match("?", ".", FNM_PERIOD)).Should.Equal(false)
			c.Then(Match("foo?", "foo.", FNM_PERIOD)).Should.Equal(true)
			c.Then(Match("/?", "/.", FNM_PERIOD)).Should.Equal(true)
			c.Then(Match("/?", "/.", FNM_PERIOD|FNM_PATHNAME)).Should.Equal(false)
		})
	})

	c.Specify("A range pattern", func() {
		p := "[a-z]"
		c.Specify("Should match a single character inside its range", func() {
			c.Then(Match(p, "a", 0)).Should.Equal(true)
			c.Then(Match(p, "q", 0)).Should.Equal(true)
			c.Then(Match(p, "z", 0)).Should.Equal(true)
			c.Then(Match("[わ]", "わ", 0)).Should.Equal(true)
		})
		c.Specify("Should not match characters outside its range", func() {
			c.Then(Match(p, "-", 0)).Should.Equal(false)
			c.Then(Match(p, " ", 0)).Should.Equal(false)
			c.Then(Match(p, "D", 0)).Should.Equal(false)
			c.Then(Match(p, "é", 0)).Should.Equal(false)
		})
		c.Specify("Should only match one character", func() {
			c.Then(Match(p, "ab", 0)).Should.Equal(false)
			c.Then(Match(p, "", 0)).Should.Equal(false)
		})
		c.Specify("Should not consume more of the pattern than necessary", func() { c.Then(Match(p+"foo", "afoo", 0)).Should.Equal(true) })

		c.Specify("Should match a - if it's the first or last character or backslash-escaped", func() {
			c.Then(Match("[-az]", "-", 0)).Should.Equal(true)
			c.Then(Match("[-az]", "a", 0)).Should.Equal(true)
			c.Then(Match("[-az]", "b", 0)).Should.Equal(false)
			c.Then(Match("[az-]", "-", 0)).Should.Equal(true)
			c.Then(Match("[a\\-z]", "-", 0)).Should.Equal(true)
			c.Then(Match("[a\\-z]", "b", 0)).Should.Equal(false)
			c.Specify("And ignore \\ when FNM_NOESCAPE is given", func() {
				c.Then(Match("[a\\-z]", "\\", FNM_NOESCAPE)).Should.Equal(true)
				c.Then(Match("[a\\-z]", "-", FNM_NOESCAPE)).Should.Equal(false)
			})
		})
		c.Specify("Should be negated if starting with ^ or !", func() {
			c.Then(Match("[^a-z]", "a", 0)).Should.Equal(false)
			c.Then(Match("[!a-z]", "b", 0)).Should.Equal(false)
			c.Then(Match("[!a-z]", "é", 0)).Should.Equal(true)
			c.Then(Match("[!a-z]", "わ", 0)).Should.Equal(true)
			c.Specify("And still match - if following the negation character", func() {
				c.Then(Match("[^-az]", "-", 0)).Should.Equal(false)
				c.Then(Match("[^-az]", "b", 0)).Should.Equal(true)
			})
		})
		c.Specify("Should support multiple characters/ranges", func() {
			c.Then(Match("[abc]", "a", 0)).Should.Equal(true)
			c.Then(Match("[abc]", "c", 0)).Should.Equal(true)
			c.Then(Match("[abc]", "d", 0)).Should.Equal(false)
			c.Then(Match("[a-cg-z]", "c", 0)).Should.Equal(true)
			c.Then(Match("[a-cg-z]", "h", 0)).Should.Equal(true)
			c.Then(Match("[a-cg-z]", "d", 0)).Should.Equal(false)
		})
		c.Specify("Should not match / when flags is FNM_PATHNAME", func() {
			c.Then(Match("[abc/def]", "/", 0)).Should.Equal(true)
			c.Then(Match("[abc/def]", "/", FNM_PATHNAME)).Should.Equal(false)
			c.Then(Match("[.-0]", "/", 0)).Should.Equal(true) // .-0 includes /
			c.Then(Match("[.-0]", "/", FNM_PATHNAME)).Should.Equal(false)
		})
		c.Specify("Should normally be case-sensitive", func() {
			c.Then(Match("[a-z]", "A", 0)).Should.Equal(false)
			c.Then(Match("[A-Z]", "a", 0)).Should.Equal(false)
			c.Specify("Except when FNM_CASEFOLD is given", func() {
				c.Then(Match("[a-z]", "A", FNM_CASEFOLD)).Should.Equal(true)
				c.Then(Match("[A-Z]", "a", FNM_CASEFOLD)).Should.Equal(true)
			})
		})
		// What about [a-c-f]? How should that behave? It's undocumented.
	})

	c.Specify("A backslash should escape the following character", func() {
		c.Then(Match("\\\\", "\\", 0)).Should.Equal(true)
		c.Then(Match("\\*", "*", 0)).Should.Equal(true)
		c.Then(Match("\\*", "foo", 0)).Should.Equal(false)
		c.Then(Match("\\?", "?", 0)).Should.Equal(true)
		c.Then(Match("\\?", "f", 0)).Should.Equal(false)
		c.Then(Match("\\[a-z]", "[a-z]", 0)).Should.Equal(true)
		c.Then(Match("\\[a-z]", "a", 0)).Should.Equal(false)
		c.Then(Match("\\foo", "foo", 0)).Should.Equal(true)
		c.Then(Match("\\わ", "わ", 0)).Should.Equal(true)
		c.Specify("Unless FNM_NOESCAPE is given", func() {
			c.Then(Match("\\\\", "\\", FNM_NOESCAPE)).Should.Equal(false)
			c.Then(Match("\\\\", "\\\\", FNM_NOESCAPE)).Should.Equal(true)
			c.Then(Match("\\*", "foo", FNM_NOESCAPE)).Should.Equal(false)
			c.Then(Match("\\*", "\\*", FNM_NOESCAPE)).Should.Equal(true)
		})
	})

	c.Specify("Literal characters should match themselves", func() {
		c.Then(Match("foo", "foo", 0)).Should.Equal(true)
		c.Then(Match("foo", "foobar", 0)).Should.Equal(false)
		c.Then(Match("foobar", "foo", 0)).Should.Equal(false)
		c.Then(Match("foo", "Foo", 0)).Should.Equal(false)
		c.Then(Match("わたし", "わたし", 0)).Should.Equal(true)
		c.Specify("And perform case-folding when FNM_CASEFOLD is given", func() {
			c.Then(Match("foo", "FOO", FNM_CASEFOLD)).Should.Equal(true)
			c.Then(Match("FoO", "fOo", FNM_CASEFOLD)).Should.Equal(true)
		})
	})

	c.Specify("FNM_LEADING_DIR should ignore trailing /*", func() {
		c.Then(Match("foo", "foo/bar", 0)).Should.Equal(false)
		c.Then(Match("foo", "foo/bar", FNM_LEADING_DIR)).Should.Equal(true)
		c.Then(Match("*", "foo/bar", FNM_PATHNAME)).Should.Equal(false)
		c.Then(Match("*", "foo/bar", FNM_PATHNAME|FNM_LEADING_DIR)).Should.Equal(true)
	})
}
