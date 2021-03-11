package cmd

import (
	"bytes"
	"context"
	"encoding/hex"
	"fmt"
	"math"
	"sync"
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

func WriteCombPlc(ctx context.Context, targetAddress, targetPort string,data map[string]interface{}) <-chan error {
	iStartDevNo := math.MaxInt32
	iEndDevNo := math.MinInt32
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
		tx := CreateSendHeader(iStartDevNo, iEndDevNo-iStartDevNo)
		data, err := hex.DecodeString(tx)
		if err != nil {
			errCh <- err
		}
		_,err = client.Write(data)
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
	cmd := fmt.Sprintf("%X", 0x1401)
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
