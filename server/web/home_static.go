package web

import _ "embed"

//go:embed home.tpl.html
var homeTpl []byte

func HomeTpl() []byte { return homeTpl }

