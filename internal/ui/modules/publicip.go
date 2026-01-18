package modules

import (
	"github.com/yourusername/bubblefetch/internal/collectors"
	"github.com/yourusername/bubblefetch/internal/ui/theme"
)

type PublicIPModule struct{}

func (m *PublicIPModule) Name() string {
	return "publicip"
}

func (m *PublicIPModule) Render(info *collectors.SystemInfo, styles theme.Styles) string {
	if info.PublicIP == "" {
		return ""
	}
	return renderField("Public IP", info.PublicIP, styles, ": ")
}
