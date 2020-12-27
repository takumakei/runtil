module github.com/takumakei/runtil

go 1.15

require (
	github.com/cpuguy83/go-md2man/v2 v2.0.0 // indirect
	github.com/russross/blackfriday/v2 v2.1.0 // indirect
	github.com/urfave/cli/v2 v2.3.0
	go.uber.org/multierr v1.6.0 // indirect
	go.uber.org/zap v1.16.0
)

replace github.com/urfave/cli/v2 => github.com/takumakei/cli/v2 v2.3.0-patch.1
