package ios

import "fmt"

type DeletePair struct {
	BundleID            string
	ClientVersionString string
	MessageType         string
	ProgName            string
	LibUSBMuxVersion    uint32 `plist:"kLibUSBMuxVersion"`
	PairRecordID        string
}

func newDeletePair(udid string) DeletePair {
	data := DeletePair{
		BundleID:            "go.ios.control",
		ClientVersionString: "go-ios-1.0.0",
		MessageType:         "DeletePairRecord",
		ProgName:            "go-ios",
		LibUSBMuxVersion:    3,
		PairRecordID:        udid,
	}
	return data
}

func (muxConn *UsbMuxConnection) DeletePair(udid string) error {
	err := muxConn.Send(newDeletePair(udid))
	if err != nil {
		return err
	}
	resp, err := muxConn.ReadMessage()
	if err != nil {
		return err
	}
	muxresponse := MuxResponsefromBytes(resp.Payload)
	if !muxresponse.IsSuccessFull() {
		return fmt.Errorf("DeletePairRecord failed with error code %d", muxresponse.Number)
	}
	return nil
}

func DeletePairRecord(udid string) error {
	muxConnection, err := NewUsbMuxConnectionSimple()
	if err != nil {
		return err
	}
	defer muxConnection.Close()
	return muxConnection.DeletePair(udid)
}
