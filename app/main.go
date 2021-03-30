package main

import (
	"bitbucket.org/latonaio/aion-core/pkg/go-client/msclient"
	"context"
	"control-mitsubishi-plc-w-kube/cmd"
	"control-mitsubishi-plc-w-kube/config"
	"control-mitsubishi-plc-w-kube/pkg"
	"log"
	"os"
	"os/signal"
	"syscall"
)

const msName = "control-mitsubishi-plc-w-kube"

func main() {
	errCh := make(chan error, 1)
	ctx, cancel := context.WithCancel(context.Background())
	cfg, err := config.New()
	if err != nil {
		errCh <- err
	}
	defer cancel()
	s := pkg.New(ctx, cfg)
	go s.Start(errCh)
	kc, err := msclient.NewKanbanClient(ctx)
	if err != nil {
		errCh <- err
		return
	}
	kw, err := kc.GetKanbanCh(msName, kc.GetProcessNumber())
	if err != nil {
		errCh <- err
		return
	}
	quitC := make(chan os.Signal, 1)
	signal.Notify(quitC, syscall.SIGTERM, os.Interrupt)
loop:
	for {
		select {
		case data := <-kw:
			if data != nil {
				res, _ := data.GetMetadataByMap()
				cmd.WriteCombPlc(ctx, cfg, res)
			}
		case err := <-errCh:
			log.Printf("err = %v", err)
			break loop
		case <-quitC:
			if err := s.Shutdown(ctx); err != nil {
				errCh <- err
			}
		}
	}
}
