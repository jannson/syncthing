package connections

import (
	"context"
	"crypto/tls"
	"net/url"

	"github.com/syncthing/syncthing/lib/appext"
	"github.com/syncthing/syncthing/lib/config"
	"github.com/syncthing/syncthing/lib/connections/registry"
	"github.com/syncthing/syncthing/lib/nat"
	"github.com/syncthing/syncthing/lib/util"
)

func init() {
	factory := &linkeaseListenerFactory{}
	listeners["link2"] = factory
}

type linkeaseListener struct {
	util.ServiceWithError
	onAddressesChangedNotifier

	uri     *url.URL
	cfg     config.Wrapper
	conns   chan internalConn
	factory listenerFactory
}

func (link *linkeaseListener) serve(ctx context.Context) error {
	l.Debugln("serve linkease listener, uri=", link.uri)
	link.notifyAddressesChanged(link)
	registry.Register(link.uri.Scheme, link.uri.Hostname())
	appConns := appext.ServerConns(link.uri.String())

LINK_LOOP:
	for {
		select {
		case <-ctx.Done():
			break LINK_LOOP
		case appConn := <-appConns:
			l.Debugln("got linkease server connection")
			link.conns <- internalConn{appConn, connTypeLinkEaseServer, tcpPriority}
		}
	}

	return nil
}

func (link *linkeaseListener) URI() *url.URL {
	return link.uri
}

func (link *linkeaseListener) WANAddresses() []*url.URL {
	return []*url.URL{link.uri}
}

func (link *linkeaseListener) LANAddresses() []*url.URL {
	return []*url.URL{}
}

func (link *linkeaseListener) String() string {
	return link.uri.String()
}

func (link *linkeaseListener) Factory() listenerFactory {
	return link.factory
}

func (link *linkeaseListener) NATType() string {
	return "unknown"
}

type linkeaseListenerFactory struct{}

func (f *linkeaseListenerFactory) New(uri *url.URL, cfg config.Wrapper, tlsCfg *tls.Config, conns chan internalConn, natService *nat.Service) genericListener {
	l := &linkeaseListener{
		uri:     uri,
		cfg:     cfg,
		factory: f,
		conns:   conns,
	}
	l.ServiceWithError = util.AsServiceWithError(l.serve, l.String())
	return l
}

func (linkeaseListenerFactory) Valid(_ config.Configuration) error {
	// Always valid
	return nil
}
