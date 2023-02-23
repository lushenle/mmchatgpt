package main

import (
	"context"
	"crypto/tls"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/lushenle/mmchatgpt/config"
	"github.com/lushenle/mmchatgpt/webhook"
	log "github.com/sirupsen/logrus"
)

func main() {
	err := config.LoadConfig()
	if err != nil {
		log.Warn("Configuration file not found, try to read env...")
	}

	var param webhook.WhSvrParam
	var whsrv webhook.WebHookServer

	flag.IntVar(&param.Port, "port", 3000, "Webhook Server Port")
	flag.StringVar(&param.CertFile, "tlsCertFile", "", "x509 certification file")
	flag.StringVar(&param.KeyFile, "tlsKeyFile", "", "x509 private key file")
	flag.Parse()

	if param.CertFile != "" && param.KeyFile != "" {
		certificate, err := tls.LoadX509KeyPair(param.CertFile, param.KeyFile)
		if err != nil {
			log.Errorf("Failed to load key pair: %s", err)
			return
		}
		whsrv = webhook.WebHookServer{
			Server: &http.Server{
				Addr: fmt.Sprintf(":%d", param.Port),
				TLSConfig: &tls.Config{
					Certificates: []tls.Certificate{certificate},
				},
			},
		}
	} else {
		whsrv = webhook.WebHookServer{
			Server: &http.Server{
				Addr: fmt.Sprintf(":%d", param.Port),
			},
		}
	}

	// http server handler
	mux := http.NewServeMux()
	mux.HandleFunc("/webhook", whsrv.ServeHTTP)
	whsrv.Server.Handler = mux

	// start webhook server
	if whsrv.Server.TLSConfig != nil {
		go func() {
			if err := whsrv.Server.ListenAndServeTLS("", ""); err != nil {
				log.Errorf("Failed to listen and serve webhook: %s", err)
			}
		}()
	} else {
		go func() {
			if err := whsrv.Server.ListenAndServe(); err != nil {
				log.Errorf("Failed to listen and serve webhook: %s", err)
			}
		}()
	}
	log.Infof("Starting server on %d...\n", param.Port)

	// os signal
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	<-signalChan
	log.Info("Got OS shutdown signal, gracefully shutting down...")
	if err := whsrv.Server.Shutdown(context.Background()); err != nil {
		log.Errorf("HTTP Server Shutdown err: %s", err)
	}
}
