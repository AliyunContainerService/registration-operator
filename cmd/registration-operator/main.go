package main

import (
	goflag "flag"
	"fmt"
	"k8s.io/klog/v2"
	"math/rand"
	"net/http"
	"os"
	log "github.com/sirupsen/logrus"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	utilflag "k8s.io/component-base/cli/flag"
	"k8s.io/component-base/logs"

	"open-cluster-management.io/registration-operator/pkg/cmd/operator"
	"open-cluster-management.io/registration-operator/pkg/version"
)

const defaultHealthzAddr = ":10254"

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	pflag.CommandLine.SetNormalizeFunc(utilflag.WordSepNormalizeFunc)
	pflag.CommandLine.AddGoFlagSet(goflag.CommandLine)

	logs.InitLogs()
	go func() {
		startHealthzServer(defaultHealthzAddr)
	}()
	defer logs.FlushLogs()

	command := newNucleusCommand()
	if err := command.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

func newNucleusCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "registration-operator",
		Short: "Nucleus Operator",
		Run: func(cmd *cobra.Command, args []string) {
			_ = cmd.Help()
			os.Exit(1)
		},
	}

	if v := version.Get().String(); len(v) == 0 {
		cmd.Version = "<unknown>"
	} else {
		cmd.Version = v
	}

	cmd.AddCommand(operator.NewHubOperatorCmd())
	cmd.AddCommand(operator.NewKlusterletOperatorCmd())

	return cmd
}

func startHealthzServer(addr string) {
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(http.StatusOK)
	})
	s := &http.Server{
		Addr:    addr,
		Handler: mux,
	}
	log.Info("start healthz server and listen on", "addr", addr)
	if err := s.ListenAndServe(); err != nil {
		log.Error(err, "healthz server cause a error")
		klog.Fatal(err)
	}
}
