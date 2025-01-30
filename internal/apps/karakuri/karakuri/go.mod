module karakuri

replace karakuripkgs => ../../../../pkgs

replace hitoha => ../hitoha

go 1.22.2

require (
	github.com/google/uuid v1.6.0
	hitoha v0.0.0-00010101000000-000000000000
)

require (
	github.com/gorilla/mux v1.8.1 // indirect
	karakuripkgs v0.0.0-00010101000000-000000000000
)
