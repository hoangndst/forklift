package main

import (
	"flag"
	libclient "github.com/konveyor/forklift-controller/pkg/lib/client/gcp"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"io"
	"k8s.io/klog/v2"
	"net/http"
	"os"
	"strings"
	"time"
)

func main() {
	var (
		googleAuthPath string
		bucketName     string
		objectName     string
		crNamespace    string
		crName         string

		volumePath string
	)
	flag.StringVar(&googleAuthPath, "google-auth-path", "", "Google Auth Path")
	flag.StringVar(&bucketName, "bucket-name", "", "Bucket Name")
	flag.StringVar(&objectName, "object-name", "", "Object Name")
	flag.StringVar(&volumePath, "volume-path", "", "Path to populate")
	flag.StringVar(&crName, "cr-name", "", "Custom Resource instance name")
	flag.StringVar(&crNamespace, "cr-namespace", "", "Custom Resource instance namespace")
	flag.Parse()

	populate(googleAuthPath, volumePath, bucketName, objectName)
}

func populate(googleAuthPath string, volumePath string, bucketName string, objectName string) {
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
	klog.Info("Downloading image from bucket...", bucketName, objectName)
	imageReader, err := client.DownloadImageFromBucket(bucketName, objectName)
	if err != nil {
		klog.Fatal(err)
	}
	defer imageReader.Close()

	flags := os.O_RDWR
	if strings.HasSuffix(volumePath, "disk.img") {
		flags |= os.O_CREATE
	}

	klog.Info("Saving the image to: ", volumePath)
	file, err := os.OpenFile(volumePath, flags, 0650)
	if err != nil {
		klog.Fatal(err)
	}
	defer file.Close()

	err = writeData(imageReader, file, objectName, progressGague)
	if err != nil {
		klog.Fatal(err)
	}
}

type CountingReader struct {
	reader io.ReadCloser
	total  *int64
}

func (cr *CountingReader) Read(p []byte) (int, error) {
	n, err := cr.reader.Read(p)
	*cr.total += int64(n)
	return n, err
}

func writeData(reader io.ReadCloser, file *os.File, imageID string, progress *prometheus.GaugeVec) error {
	total := new(int64)
	countingReader := CountingReader{reader, total}

	done := make(chan bool)
	go func() {
		for {
			select {
			case <-done:
				klog.Info("Total: ", *total)
				klog.Info("Finished!")
				return
			default:
				progress.WithLabelValues(imageID).Set(float64(*total))
				klog.Info("Transferred: ", *total)
				time.Sleep(3 * time.Second)
			}
		}
	}()

	if _, err := io.Copy(file, &countingReader); err != nil {
		klog.Fatal(err)
	}
	done <- true
	progress.WithLabelValues(imageID).Set(float64(*total))

	return nil
}
