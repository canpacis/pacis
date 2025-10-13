module github.com/canpacis/pacis/registry

go 1.24.5

replace github.com/canpacis/pacis => ../

replace components => .

require (
	components v0.0.0-00010101000000-000000000000
	github.com/canpacis/pacis v0.3.0
)

require github.com/Oudwins/tailwind-merge-go v0.2.1 // indirect
