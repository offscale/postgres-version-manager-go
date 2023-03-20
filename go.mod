module postgres-version-manager-go

go 1.20

require (
	github.com/fergusstrange/embedded-postgres v1.21.0
	github.com/mkideal/cli v0.2.7
)

require (
	github.com/labstack/gommon v0.3.0 // indirect
	github.com/lib/pq v1.10.4 // indirect
	github.com/mattn/go-colorable v0.1.7 // indirect
	github.com/mattn/go-isatty v0.0.12 // indirect
	github.com/mkideal/expr v0.1.0 // indirect
	github.com/xi2/xz v0.0.0-20171230120015-48954b6210f8 // indirect
	golang.org/x/crypto v0.0.0-20201221181555-eec23a3978ad // indirect
	golang.org/x/sys v0.0.0-20200323222414-85ca7c5b95cd // indirect
	golang.org/x/term v0.0.0-20201117132131-f5c789dd3221 // indirect
)

replace github.com/fergusstrange/embedded-postgres v1.21.0 => github.com/offscale/embedded-postgres v0.0.0-20230320182601-7fe884ab92d8
