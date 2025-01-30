module futaba

go 1.22.2

replace karakuripkgs => ../../../../pkgs

require (
	golang.org/x/sys v0.29.0
	karakuripkgs v0.0.0-00010101000000-000000000000
)

require github.com/google/uuid v1.6.0 // indirect
