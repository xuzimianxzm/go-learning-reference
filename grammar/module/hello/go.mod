module hello

go 1.13

replace xuzimian.com/grammar/module/greetings => ../greetings

require (
	github.com/pkg/errors v0.9.1 // indirect
	xuzimian.com/grammar/module/greetings v1.1.0
)
