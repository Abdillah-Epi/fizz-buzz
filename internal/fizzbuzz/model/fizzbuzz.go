package model

type Request struct {
	Int1  int
	Int2  int
	Limit int
	Str1  string
	Str2  string
}

type Result struct {
	Values []string
}
