module github.com/JacobPotter/go-zendesk

go 1.22

require (
	github.com/google/go-querystring v1.1.0
	github.com/stretchr/testify v1.9.0
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	golang.org/x/exp v0.0.0-20240719175910-8a7402abbf56
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

retract (
	// Versions that were for testing knowledge share
	v0.33.2
	v0.33.1
)