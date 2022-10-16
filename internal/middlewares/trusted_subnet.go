package middlewares

import (
	"net"
	"net/http"
)

func TrustedSubnet(cidr string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			if cidr != "" {
				if r.RemoteAddr == "" {
					http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)

					return
				}

				ipNet := net.ParseIP(r.RemoteAddr)

				_, cidrNet, err := net.ParseCIDR(cidr)
				if err != nil {
					http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)

					return
				}

				if !cidrNet.Contains(ipNet) {
					http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)

					return
				}
			}

			next.ServeHTTP(w, r)
		}

		return http.HandlerFunc(fn)
	}
}
