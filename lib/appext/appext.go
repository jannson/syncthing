package appext

import (
	"context"
	"crypto/tls"
	"errors"
	"net"
	"net/http"
	"net/url"

	"github.com/syncthing/syncthing/lib/events"
	"github.com/syncthing/syncthing/lib/protocol"
)

var ErrNext = errors.New("next")

type GuiHandler interface {
	GetHttpHandle() http.Handler
	GetDiskEvents(folderId string) []events.Event
}

type Manager interface {
	SetGuiHandler(localAddr string, guiHandle GuiHandler)
	GetGuiListener(network, localAddr string) (net.Listener, error)
	GetDeviceID(state *tls.ConnectionState) (protocol.DeviceID, error)
	ServerConns(linkHost string) <-chan Conn
	Dial(ctx context.Context, localAddr string, deviceId protocol.DeviceID, uri *url.URL) (net.Conn, error)
}

type Conn interface {
	net.Conn
	ConnectionState() tls.ConnectionState
}

var appMgr Manager

func Dial(ctx context.Context, localAddr string, deviceId protocol.DeviceID, uri *url.URL) (net.Conn, error) {
	return appMgr.Dial(ctx, localAddr, deviceId, uri)
}

func GetDeviceID(state *tls.ConnectionState) (protocol.DeviceID, error) {
	if appMgr != nil {
		return appMgr.GetDeviceID(state)
	}
	const clCount = 1
	certs := state.PeerCertificates
	if cl := len(certs); cl != clCount {
		return protocol.DeviceID{}, errors.New("got peer  certificate list of length !=1")
	}
	remoteCert := certs[0]
	remoteID := protocol.NewDeviceID(remoteCert.Raw)
	return remoteID, nil
}

func ServerConns(linkHost string) <-chan Conn {
	return appMgr.ServerConns(linkHost)
}

func SetGuiHandler(localAddr string, guiHandle GuiHandler) {
	appMgr.SetGuiHandler(localAddr, guiHandle)
}

func GetGuiListener(network, localAddr string) (net.Listener, error) {
	return appMgr.GetGuiListener(network, localAddr)
}

func SetMgr(p Manager) {
	appMgr = p
}
