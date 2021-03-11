package plc_mitsubishi

import "fmt"

type RequestBuilder struct {
	*CommonMessageFormat
	ReqDataLen int
	WatchTimer string
	IO         int
}

func NewRequestBuilder(io int, watchTimer int) *RequestBuilder {
	msgFmt := newFmt()
	return &RequestBuilder{
		CommonMessageFormat: msgFmt,
		ReqDataLen:          12,
		WatchTimer:          "1000",
		IO:                  io,
	}
}

func (rb *RequestBuilder) buildRequestMsg(startDevNum int, writeLen int) string {
	var cmd string
	if rb.IO == INPUT {
		cmd = BULK_READ_CMD
	} else {
		cmd = BULK_WRITE_CMD
	}

	return rb.Subheader +
		rb.AccessRoute.NetworkNo +
		rb.AccessRoute.PcNo +
		rb.AccessRoute.IONo +
		rb.AccessRoute.StationNo +
		fmt.Sprintf("%X", rb.ReqDataLen) +
		rb.WatchTimer +
		cmd +
		SUB_CMD +
		fmt.Sprintf("%X", startDevNum) +
		fmt.Sprintf("%X", writeLen)
}

func CreateSendHeader(startDevNum int, writeLen int) string {
	// subheader
	subHeader := "5000"
	//ネットワーク番号
	netNum := fmt.Sprintf("%X", 0)
	//PC番号
	pcNum := fmt.Sprintf("%X", 0xFF)
	//要求先ユニットI/O番号
	io := fmt.Sprintf("%X", 0x3FF)
	//要求先ユニット局番号
	unit := fmt.Sprintf("%X", 0)
	//要求データ長
	dataLen := fmt.Sprintf("%X", 12)
	//CPU監視タイマ
	cpuTimer := fmt.Sprintf("%X", 0x1)
	//コマンド
	cmd := fmt.Sprintf("%X", 0x0401)
	//サブコマンド
	subCmd := fmt.Sprintf("%X", 0x00)
	//要求データ部
	startDev := fmt.Sprintf("%X", startDevNum)
	wLen := fmt.Sprintf("%X", writeLen)

	return subHeader +
		netNum +
		pcNum +
		io +
		unit +
		dataLen +
		cpuTimer +
		cmd +
		subCmd +
		startDev +
		wLen
}
