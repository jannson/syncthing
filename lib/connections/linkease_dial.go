package connections

import (
	"context"
	"crypto/tls"
	"net"
	"net/url"
	"time"

	"github.com/syncthing/syncthing/lib/appext"
	"github.com/syncthing/syncthing/lib/config"
	"github.com/syncthing/syncthing/lib/protocol"
)

func init() {
	factory := &linkeaseDialerFactory{}
	dialers["link2"] = factory
}

type linkeaseDialer struct {
	commonDialer
}

type linkeaseConn struct {
	net.Conn
}

func (c linkeaseConn) ConnectionState() tls.ConnectionState {
	return tls.ConnectionState{}
}

func (d *linkeaseDialer) Dial(ctx context.Context, deviceId protocol.DeviceID, uri *url.URL) (internalConn, error) {
	tc, err := appext.Dial(ctx, deviceId, uri)
	if err != nil {
		return internalConn{}, err
	}
	return internalConn{linkeaseConn{tc}, connTypeLinkEaseClient, tcpPriority}, nil
}

type linkeaseDialerFactory struct{}

func (linkeaseDialerFactory) New(opts config.OptionsConfiguration, tlsCfg *tls.Config) genericDialer {
	return &linkeaseDialer{commonDialer{
		trafficClass:      opts.TrafficClass,
		reconnectInterval: time.Duration(opts.ReconnectIntervalS) * time.Second,
		tlsCfg:            tlsCfg,
	}}
}

func (linkeaseDialerFactory) Priority() int {
	return tcpPriority
}

func (linkeaseDialerFactory) AlwaysWAN() bool {
	return false
}

func (linkeaseDialerFactory) Valid(_ config.Configuration) error {
	// Always valid
	return nil
}

func (linkeaseDialerFactory) String() string {
	return "LinkEase Dialer"
}
