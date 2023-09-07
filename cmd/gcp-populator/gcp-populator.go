package main

import (
	"flag"
	libclient "github.com/konveyor/forklift-controller/pkg/lib/client/gcp"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"k8s.io/klog/v2"
	"net/http"
)

func main() {
	var (
		googleAuthPath string
		crNamespace    string
		crName         string

		volumePath string
	)
	flag.StringVar(&googleAuthPath, "google-auth-path", "", "Google Auth Path")
	flag.StringVar(&volumePath, "volume-path", "", "Path to populate")
	flag.StringVar(&crName, "cr-name", "", "Custom Resource instance name")
	flag.StringVar(&crNamespace, "cr-namespace", "", "Custom Resource instance namespace")
	flag.Parse()
}

func populate(googleAuthPath string, volumePath string) {
	http.Handle("/metrics", promhttp.Handler())
	go http.ListenAndServe(":2112", nil)
	progressGague := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Subsystem: "volume_populators",
			Name:      "gcp_volume_populator",
			Help:      "Amount of data transferred",
		},
		[]string{"image_id"},
	)
	if err := prometheus.Register(progressGague); err != nil {
		klog.Error("Prometheus progress counter not registered:", err)
	} else {
		klog.Info("Prometheus progress counter registered.")
	}

	client := libclient.Client{
		GoogleAuthPath: googleAuthPath,
	}
	err := client.Connect()
	if err != nil {
		klog.Error("Error connecting to GCP:", err)
		return
	}
}
