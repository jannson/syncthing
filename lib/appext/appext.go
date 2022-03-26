package appext

import (
	"context"
	"crypto/tls"
	"net"
	"net/url"

	"github.com/syncthing/syncthing/lib/protocol"
)

type Manager interface {
	GetDeviceID(state *tls.ConnectionState) (protocol.DeviceID, error)
	ServerConns(linkHost string) <-chan Conn
	Dial(ctx context.Context, deviceId protocol.DeviceID, uri *url.URL) (net.Conn, error)
}

type Conn interface {
	net.Conn
	ConnectionState() tls.ConnectionState
}

var appMgr Manager

func Dial(ctx context.Context, deviceId protocol.DeviceID, uri *url.URL) (net.Conn, error) {
	return appMgr.Dial(ctx, deviceId, uri)
}

func GetDeviceID(state *tls.ConnectionState) (protocol.DeviceID, error) {
	return appMgr.GetDeviceID(state)
}

func ServerConns(linkHost string) <-chan Conn {
	return appMgr.ServerConns(linkHost)
}

func SetMgr(p Manager) {
	appMgr = p
}
