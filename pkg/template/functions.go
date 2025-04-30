package template

import (
	"fmt"
	"regexp"
	"strings"
	"text/template"
)

var StringFuncs = template.FuncMap{
	"regexMatch":           regexMatch,
	"regexFind":            regexFind,
	"regexFindAll":         regexFindAll,
	"regexFindSubmatch":    regexFindSubmatch,
	"regexFindAllSubmatch": regexFindAllSubmatch,
	"regexReplaceAll":      regexReplaceAll,
	"regexSplit":           regexSplit,
	"regexSplitN":          regexSplitN,
	"join":                 join,
}

// regexMatch compiles the given regex and checks if it matches the input. It
// returns an error if the regex fails to compile.
func regexMatch(regex string, s string) (bool, error) {
	return regexp.MatchString(regex, s)
}

// regexFind returns the first fragment of the input than matches the regex or
// an empty string if nothing matches.  It returns an error if the regex fails
// to compile.
func regexFind(regex string, s string) (string, error) {
	r, err := regexp.Compile(regex)
	if err != nil {
		return "", err
	}
	return r.FindString(s), nil
}

// regexFindAll returns a list of all fragments of the input that match the
// regex or an empty list if nothing matches.  It returns an error if the regex
// fails to compile.
func regexFindAll(regex string, s string) ([]string, error) {
	r, err := regexp.Compile(regex)
	if err != nil {
		return []string{}, err
	}
	return r.FindAllString(s, -1), nil
}

// regexFindSubmatch returns the i-th sub-match from the first fragment of the
// input that matches the regex or an empty string if nothing matches. It
// returns an error if the regex fails to compile or if the sub-match index is
// out of bounds.
func regexFindSubmatch(regex string, i int, s string) (string, error) {
	r, err := regexp.Compile(regex)
	if err != nil {
		return "", err
	}
	if i > r.NumSubexp() {
		return "", fmt.Errorf("invalid submatch index %d, regex has %d submatches", i, r.NumSubexp())
	}
	var m = r.FindStringSubmatch(s)
	if m != nil {
		return m[i], nil
	}
	return "", nil
}

// regexFindAllSubmatch returns a list of the requested 1-based sub-matches from
// all the fragments of the input that match the regex or an empty list if
// nothing matches. It returns an error if the regex fails to compile or if the
// sub-match index is out of bounds.
func regexFindAllSubmatch(regex string, i int, s string) ([]string, error) {
	r, err := regexp.Compile(regex)
	if err != nil {
		return []string{}, err
	}
	if i > r.NumSubexp() {
		return []string{}, fmt.Errorf(
			"invalid submatch index %d, regex has %d submatches", i, r.NumSubexp())
	}
	var mm = r.FindAllStringSubmatch(s, -1)
	var result = make([]string, 0, len(mm))
	for _, m := range mm {
		result = append(result, m[i])
	}
	return result, nil
}

// regexReplaceAll replaces all matches of the regex in the input with the
// replacement string which can contain sub-match reference ($1, $2, ...). It
// returns an error if the regex fails to compile.
func regexReplaceAll(regex string, repl string, s string) (string, error) {
	r, err := regexp.Compile(regex)
	if err != nil {
		return "", err
	}
	return r.ReplaceAllString(s, repl), nil
}

// regexSplit splits the input into a list of strings using the regex as a
// delimiter. It returns an error if the regex fails to compile.
func regexSplit(regex string, s string) ([]string, error) {
	r, err := regexp.Compile(regex)
	if err != nil {
		return []string{}, err
	}
	return r.Split(s, -1), nil
}

// regexSplitN splits the input into a list of strings using the regex as a
// delimiter. The n parameter specifies the maximum number of substrings to
// return, with the last substring containing the remainder unsplit. It returns
// an error if the regex fails to compile.
func regexSplitN(regex string, n int, s string) ([]string, error) {
	r, err := regexp.Compile(regex)
	if err != nil {
		return []string{}, err
	}
	return r.Split(s, n), nil
}

func join(sep string, s []string) string {
	if len(s) == 0 {
		return ""
	}
	if len(s) == 1 {
		return s[0]
	}

	var result strings.Builder
	result.WriteString(s[0])
	for _, v := range s[1:] {
		result.WriteString(sep)
		result.WriteString(v)
	}
	return result.String()
}
