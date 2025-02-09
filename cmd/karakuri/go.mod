module karakurimain

go 1.22.2

replace karakuripkgs => ../../pkgs

replace karakuri_mod => ../../internal/apps/karakuri/hitoha/module

replace futaba => ../../internal/apps/karakuri/futaba

replace hitoha => ../../internal/apps/karakuri/hitoha

replace karakuri => ../../internal/apps/karakuri/karakuri

require (
	github.com/urfave/cli v1.22.16
	karakuri v0.0.0-00010101000000-000000000000
	karakuri_mod v0.0.0-00010101000000-000000000000
)

require (
	futaba v0.0.0-00010101000000-000000000000 // indirect
	github.com/cpuguy83/go-md2man/v2 v2.0.5 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/gorilla/mux v1.8.1 // indirect
	github.com/russross/blackfriday/v2 v2.1.0 // indirect
	golang.org/x/sys v0.29.0 // indirect
	hitoha v0.0.0-00010101000000-000000000000 // indirect
	karakuripkgs v0.0.0-00010101000000-000000000000 // indirect
)
