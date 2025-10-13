module karakuri

replace karakuripkgs => ../../../../pkgs

replace karakuri_mod => ../hitoha/module

replace futaba => ../futaba

replace hitoha => ../hitoha

go 1.22.2

require (
	github.com/google/uuid v1.6.0
	hitoha v0.0.0-00010101000000-000000000000
	karakuri_mod v0.0.0-00010101000000-000000000000
)

require (
	futaba v0.0.0-00010101000000-000000000000 // indirect
	golang.org/x/sys v0.29.0 // indirect
)

require (
	github.com/gorilla/mux v1.8.1 // indirect
	karakuripkgs v0.0.0-00010101000000-000000000000
)
