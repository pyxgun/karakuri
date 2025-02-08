module hitoha

go 1.22.2

replace futaba => ../futaba

replace karakuripkgs => ../../../../pkgs

replace karakuri_mod => ./module

require (
	futaba v0.0.0-00010101000000-000000000000
	github.com/gorilla/mux v1.8.1
	karakuri_mod v0.0.0-00010101000000-000000000000
	karakuripkgs v0.0.0-00010101000000-000000000000
)

require (
	github.com/google/uuid v1.6.0 // indirect
	golang.org/x/sys v0.29.0 // indirect
)
