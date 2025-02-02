module hitohamain

go 1.22.2

replace karakuripkgs => ../../pkgs

replace futaba => ../../internal/apps/karakuri/futaba

replace hitoha => ../../internal/apps/karakuri/hitoha

require (
	github.com/gorilla/mux v1.8.1
	hitoha v0.0.0-00010101000000-000000000000
)

require (
	futaba v0.0.0-00010101000000-000000000000 // indirect
	github.com/google/uuid v1.6.0 // indirect
	golang.org/x/sys v0.29.0 // indirect
	karakuripkgs v0.0.0-00010101000000-000000000000 // indirect
)
