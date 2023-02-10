package sigv4

import (
	"fmt"

	"github.com/grafana/agent/component"
	"github.com/grafana/agent/component/otelcol/auth"
	"github.com/grafana/agent/pkg/river"
	"github.com/open-telemetry/opentelemetry-collector-contrib/extension/sigv4authextension"
	otelcomponent "go.opentelemetry.io/collector/component"
	otelconfig "go.opentelemetry.io/collector/config"
)

func init() {
	component.Register(component.Registration{
		Name:    "otelcol.auth.sigv4",
		Args:    Arguments{},
		Exports: auth.Exports{},

		Build: func(opts component.Options, args component.Arguments) (component.Component, error) {
			fact := sigv4authextension.NewFactory()
			return auth.New(opts, fact, args.(Arguments))
		},
	})
}

// Arguments configures the otelcol.auth.headers component.
type Arguments struct {
	Region     string     `river:"region,attr,optional"`
	Service    string     `river:"service,attr,optional"`
	AssumeRole AssumeRole `river:"assume_role,block,optional"`
}

var _ auth.Arguments = Arguments{}

// Convert implements auth.Arguments.
func (args Arguments) Convert() otelconfig.Extension {
	return &sigv4authextension.Config{
		ExtensionSettings: otelconfig.NewExtensionSettings(otelconfig.NewComponentID("sigv4")),
		Region:            args.Region,
		Service:           args.Service,
		AssumeRole:        args.AssumeRole,
	}
}

// Extensions implements auth.Arguments.
func (args Arguments) Extensions() map[otelconfig.ComponentID]otelcomponent.Extension {
	return nil
}

// Exporters implements auth.Arguments.
func (args Arguments) Exporters() map[otelconfig.DataType]map[otelconfig.ComponentID]otelcomponent.Exporter {
	return nil
}

type AssumeRole struct {
	ARN         string `river:"arn,attr,optional"`
	SessionName string `river:"session_name,attr,optional"`
	STSRegion   string `river:"sts_region,attr,optional"`
}

var _ river.Unmarshaler = (*AssumeRole)(nil)

// UnmarshalRiver implements river.Unmarshaler.
func (h *AssumeRole) UnmarshalRiver(f func(interface{}) error) error {
	type header Header
	if err := f((*header)(h)); err != nil {
		return err
	}

	switch {
	case h.Key == "":
		return fmt.Errorf("key must be set to a non-empty string")
	case h.FromContext == nil && h.Value == nil:
		return fmt.Errorf("either value or from_context must be provided")
	case h.FromContext != nil && h.Value != nil:
		return fmt.Errorf("either value or from_context must be provided, not both")
	}

	return nil
}
