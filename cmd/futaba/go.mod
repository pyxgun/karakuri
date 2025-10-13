module futabamain

go 1.24.0

replace futaba => ../../internal/apps/karakuri/futaba

replace karakuripkgs => ../../pkgs

require (
	futaba v0.0.0-00010101000000-000000000000
	github.com/urfave/cli v1.22.16
	karakuripkgs v0.0.0-00010101000000-000000000000
)

require (
	github.com/cpuguy83/go-md2man/v2 v2.0.5 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/russross/blackfriday/v2 v2.1.0 // indirect
	golang.org/x/sys v0.37.0 // indirect
)
