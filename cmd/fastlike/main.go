package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/khan/fastlike"
)

func main() {
	var wasm = flag.String("wasm", "", "wasm program to execute")
	var bind = flag.String("b", "localhost:5000", "address to bind to")
	var proxyTo = flag.String("proxy-to", "", "(required) override to send all subrequests to")
	flag.Parse()

	if *wasm == "" {
		fmt.Fprintf(flag.CommandLine.Output(), "-wasm argument is required\n")
		flag.Usage()
		os.Exit(1)
	}

	var opts = []fastlike.InstanceOption{}

	// NOTE: You probably want to specify a proxy-to, otherwise any requests that get proxied
	// without changing the hostname will loop and be blocked.
	if *proxyTo != "" {
		opts = append(opts, fastlike.SubrequestHandlerOption(func(_ string, r *http.Request) (*http.Response, error) {
			r.URL.Host = *proxyTo
			return http.DefaultClient.Do(r)
		}))
	}

	proxy := fastlike.New(*wasm, opts...)

	fmt.Printf("Listening on %s\n", *bind)
	http.ListenAndServe(*bind, proxy)
}
