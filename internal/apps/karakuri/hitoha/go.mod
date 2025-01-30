module hitoha

go 1.22.2

replace karakuripkgs => ../../../../pkgs

require (
	github.com/gorilla/mux v1.8.1
	karakuripkgs v0.0.0-00010101000000-000000000000
)

require github.com/google/uuid v1.6.0 // indirect
