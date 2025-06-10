package proxy

import (
	"fmt"
	"pb_launcher/configs"

	"github.com/fatih/color"
)

func PrintProxyInfo(c configs.Config) {
	regular := color.New()

	printURLs := func(scheme, port string) {
		pub := fmt.Sprintf("%s://%s", scheme, c.GetDomain())
		if (scheme == "http" && port != "80") ||
			(scheme == "https" && port != "443") {
			pub = fmt.Sprintf("%s:%s", pub, port)
		}
		regular.Printf("├─ Proxy: %s\n", color.CyanString(pub))
	}

	printURLs("http", c.GetBindPort())

	if c.UseHttps() {
		printURLs("https", c.GetBindHttpsPort())
	}
}
