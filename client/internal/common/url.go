package common

func GenerateWsUrl(tls bool, address string) string {
	var scheme string
	if tls {
		scheme = "wss"
	} else {
		scheme = "ws"
	}
	return scheme + "://" + address
}

func GenerateHttpUrl(tls bool, address string) string {
	var scheme string
	if tls {
		scheme = "https"
	} else {
		scheme = "http"
	}
	return scheme + "://" + address
}
