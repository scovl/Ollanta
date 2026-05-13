module github.com/scovl/ollanta/ollantacore

go 1.21

require (
	github.com/BurntSushi/toml v1.5.0
	github.com/scovl/ollanta/domain v0.0.0
	go.opentelemetry.io/otel v1.38.0
)

replace github.com/scovl/ollanta/domain => ../domain
