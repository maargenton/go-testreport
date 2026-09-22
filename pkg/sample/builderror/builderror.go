//go:build sample

package builderror

func Foo(s string) string {
	return s
}

func Bar(s string) string {
	var v = 123
	return s
}
