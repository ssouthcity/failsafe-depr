package main

func NewIntPtr(s int) *int {
	return &s
}

func ListContainsStr(l []string, t string) bool {
	for _, s := range l {
		if s == t {
			return true
		}
	}
	return false
}
