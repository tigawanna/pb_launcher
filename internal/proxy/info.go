package proxy

import (
	"pb_launcher/configs"
	"pb_launcher/utils/networktools"

	"github.com/fatih/color"
)

func PrintProxyInfo(c configs.Config) {
	regular := color.New()
	regular.Printf("├─ Proxy: %s\n",
		color.CyanString(
			networktools.BuildHostURL("http", c.GetDomain(), c.GetHttpPort()),
		),
	)
	if c.IsHttpsEnabled() {
		regular.Printf("├─ Proxy: %s\n",
			color.CyanString(
				networktools.BuildHostURL("https", c.GetDomain(), c.GetHttpsPort()),
			),
		)
	}
}
