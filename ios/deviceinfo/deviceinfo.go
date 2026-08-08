package deviceinfo

import (
	"fmt"

	"github.com/danielpaulus/go-ios/ios"
	"github.com/danielpaulus/go-ios/ios/coredevice"
	"github.com/danielpaulus/go-ios/ios/xpc"
	"github.com/google/uuid"
)

type Connection struct {
	conn     *xpc.Connection
	deviceId string
}

func NewDeviceInfo(deviceEntry ios.DeviceEntry) (*Connection, error) {
	xpcConn, err := ios.ConnectToXpcServiceTunnelIface(deviceEntry, "com.apple.coredevice.deviceinfo")
	if err != nil {
		return nil, fmt.Errorf("NewDeviceInfo: failed to connect to XPC service: %w", err)
	}
	return &Connection{conn: xpcConn, deviceId: uuid.New().String()}, nil
}

func (c Connection) Close() error {
	return c.conn.Close()
}

func (c Connection) GetDisplayInfo() (map[string]interface{}, error) {
	request := coredevice.BuildRequest(c.deviceId, "com.apple.coredevice.feature.getdisplayinfo", map[string]interface{}{})
	err := c.conn.Send(request, xpc.HeartbeatRequestFlag)
	if err != nil {
		return nil, fmt.Errorf("GetDisplayInfo: failed to send 'getdisplayinfo' request: %w", err)
	}

	response, err := c.conn.ReceiveOnServerClientStream()
	if err != nil {
		return nil, fmt.Errorf("GetDisplayInfo: failed to receive response: %w", err)
	}

	if output, ok := response["CoreDevice.output"].(map[string]interface{}); ok {
		return output, nil
	}

	return nil, fmt.Errorf("GetDisplayInfo: could not parse response %+v", response)
}

func (c Connection) QueryMobileGestalt(keys []string) (map[string]interface{}, error) {
	request := coredevice.BuildRequest(c.deviceId, "com.apple.coredevice.feature.querymobilegestalt", map[string]interface{}{
		"keys": keys,
	})
	if err := c.conn.Send(request, xpc.HeartbeatRequestFlag); err != nil {
		return nil, fmt.Errorf("QueryMobileGestalt: failed to send request: %w", err)
	}

	response, err := c.conn.ReceiveOnServerClientStream()
	if err != nil {
		return nil, fmt.Errorf("QueryMobileGestalt: failed to receive response: %w", err)
	}

	if output, ok := response["CoreDevice.output"].(map[string]interface{}); ok {
		return output, nil
	}

	return nil, fmt.Errorf("QueryMobileGestalt: could not parse response %+v", response)
}
