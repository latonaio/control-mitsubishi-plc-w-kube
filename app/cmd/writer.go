package cmd

import (
	"bytes"
	"context"
	"control-mitsubishi-plc-w-kube/lib"
	"encoding/hex"
	"fmt"
	"sync"
)

const (
	REQUEST_SUBHEADER          = "5000"
	HOST_STATION_NO            = "00"
	NETWORK_NO_TO_HOST_STATION = "00"
	PC_NO_TO_HOST_STATION      = "FF"
	TARGET_IO_UNIT_NO          = "FF03"
	BULK_WRITE_CMD             = "1401"
	SUB_CMD                    = "0000"
	WATCH_TIMER                = "1000"
)

type NisPlcMakerSetting struct {
	Content      string
	DataSize     int    // データサイズ(word単位)
	DeviceNumber string // デバイス番号
}

var rcvBufferPool = &sync.Pool{
	New: func() interface{} {
		return bytes.Buffer{}
	},
}

func WriteCombPlc(ctx context.Context, targetAddress, targetPort string, data map[string]interface{}) <-chan error {
	errCh := make(chan error, 1)
	pClient := &PlcClient{}
	client, err := pClient.NewClient(targetAddress, targetPort)
	if err != nil {
		errCh <- err
	}

	// 同時に一回までの書き込みしか許容しない
	var lock sync.Mutex
	go func() {
		lock.Lock()
		defer lock.Unlock()
		initReceiveStream()
		var statusHex string
		//録音開始
		if data["status"] == 0 {
			statusHex = "0100" //録音開始のデバイスnoを1に、録音停止のデバイスnoを0に
			//録音終了
		} else if data["status"] == 1 {
			statusHex = "0001" //録音開始のデバイスnoを0に、録音停止のデバイスnoを1に
		}

		// 録音開始/停止のデバイスnoは8600からの並びなので、始点は8600で固定
		// 録音開始と録音停止のデバイスnoを書き換えるので書き込むデータ長は2で固定
		tx := CreateSendHeader(8600, 2)
		tx = tx + statusHex
		data, err := hex.DecodeString(tx)
		if err != nil {
			errCh <- err
		}
		_, err = client.Write(data)
		if err != nil {
			errCh <- err
		}
	}()
	return errCh
}

func initReceiveStream() {
	rcvBuf := rcvBufferPool.Get().(*bytes.Buffer)
	rcvBuf.Reset()
}

func CreateSendHeader(startDevNum int, writeLen int) string {
	s := lib.GetBytesFrom32bitWithLE(int64(startDevNum))
	s[3] = byte(0xA8)
	startDev := fmt.Sprintf("%X", s)
	wLen := fmt.Sprintf("%X", lib.GetBytesFrom8bitWithLE(int64(writeLen))[0:2])
	dataLen := fmt.Sprintf("%X", len(WATCH_TIMER+BULK_WRITE_CMD+SUB_CMD+startDev+wLen))

	return REQUEST_SUBHEADER +
		NETWORK_NO_TO_HOST_STATION +
		PC_NO_TO_HOST_STATION +
		TARGET_IO_UNIT_NO +
		HOST_STATION_NO +
		dataLen + //固定
		WATCH_TIMER + //固定
		BULK_WRITE_CMD +
		SUB_CMD +
		startDev +
		wLen
}
