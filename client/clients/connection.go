package clients

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/Doridian/wsvpn/shared"
	"github.com/Doridian/wsvpn/shared/sockets/adapters"
	"github.com/gorilla/websocket"
	"github.com/marten-seemann/webtransport-go"
)

func (c *Client) connectAdapter() error {
	clientUrl := *c.ServerUrl

	clientUrl.Scheme = strings.ToLower(clientUrl.Scheme)

	switch clientUrl.Scheme {
	case "webtransport":
		clientUrl.Scheme = "https"
		dialer := webtransport.Dialer{}
		dialer.TLSClientConf = c.TLSConfig

		if c.ProxyUrl != nil {
			return errors.New("proxy is not support for WebTransport at the moment")
		}

		_, conn, err := dialer.Dial(context.Background(), clientUrl.String(), c.Headers)
		if err != nil {
			return err
		}

		c.adapter = adapters.NewWebTransportAdapter(conn, false)
	case "ws":
	case "wss":
		dialer := websocket.Dialer{}
		if c.ProxyUrl != nil {
			c.log.Printf("Using HTTP proxy %s", c.ProxyUrl.Redacted())
			dialer.Proxy = func(_ *http.Request) (*url.URL, error) {
				return c.ProxyUrl, nil
			}
		}
		dialer.TLSClientConfig = c.TLSConfig

		conn, _, err := dialer.Dial(clientUrl.String(), c.Headers)
		if err != nil {
			return err
		}

		c.adapter = adapters.NewWebSocketAdapter(conn)
	default:
		panic(fmt.Errorf("invalid protocol: %s", clientUrl.Scheme))
	}

	tlsConnState, ok := c.adapter.GetTLSConnectionState()
	if ok {
		c.log.Printf("TLS %s %s connection established with cipher=%s", shared.TlsVersionString(tlsConnState.Version), c.adapter.Name(), tls.CipherSuiteName(tlsConnState.CipherSuite))
	} else {
		c.log.Printf("Unencrypted %s connection established", c.adapter.Name())
	}

	return nil
}
