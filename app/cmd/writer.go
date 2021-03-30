package cmd

import (
	"bytes"
	"context"
	"control-mitsubishi-plc-w-kube/config"
	"control-mitsubishi-plc-w-kube/lib"
	"encoding/hex"
	"fmt"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
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

type NisPlcMakerSettings struct {
	Settings []*NisPlcMakerSetting `yaml:"settings"`
}

type NisPlcMakerSetting struct {
	Content      string `yaml:"strContent"`
	DataSize     int    `yaml:"iDataSize"`
	DeviceNumber string `yaml:"strDevNo"`
	IO           int    `yaml:"iReadWrite"`
	FlowNumber   int    `yaml:"iFlowNo"`
}

var rcvBufferPool = &sync.Pool{
	New: func() interface{} {
		return bytes.Buffer{}
	},
}

//一文字目がアルファベット、二文字目以降が数値という組み合わせ以外であればエラーを返す
func GetDevNo(strDevNo string) (iDevNo int, err error) {
	iDevNo = 0
	m, _ := regexp.MatchString(`^[a-fA-F\\b]+$`, strDevNo[0:1])
	if !m {
		return 0, errors.New("デバイス番号エラー")
	}
	devNo, err := strconv.Atoi(strDevNo[1:])
	if err != nil {
		return 0, errors.New("デバイス番号エラー")
	}

	return devNo, nil
}


func LoadPlcSettings(cfg *config.Config) (*NisPlcMakerSettings, error) {
	f, err := os.Open(cfg.YamlPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	settings := &NisPlcMakerSettings{}
	b, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	if err := yaml.Unmarshal(b, settings); err != nil {
		return nil, err
	}

	return settings, nil
}


func WriteCombPlc(ctx context.Context, cfg *config.Config, data map[string]interface{}) <-chan error {
	errCh := make(chan error, 1)

	pcs, err := LoadPlcSettings(cfg)
	if err != nil {
		errCh <- err
	}
	pClient := &PlcClient{}
	client, err := pClient.NewClient(cfg.Server.Addr, cfg.Server.Port)
	if err != nil {
		errCh <- err
	}

	for _,v := range pcs.Settings {
		// 同時に一回までの書き込みしか許容しない
		var lock sync.Mutex
		go func(v *NisPlcMakerSetting) {
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

			devNo,err := GetDevNo(v.DeviceNumber)
			tx := CreateSendHeader(devNo, v.DataSize)
			tx = tx + statusHex
			data, err := hex.DecodeString(tx)
			if err != nil {
				errCh <- err
			}
			_, err = client.Write(data)
			if err != nil {
				errCh <- err
			}
		}(v)
	}

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
