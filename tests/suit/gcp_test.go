package suit

import (
	"context"
	"fmt"
	"log"

	"cloud.google.com/go/storage"
)

func main() {
	// Create a context and client
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	// List the buckets in your project
	buckets := client.Buckets(ctx, "your-project-id")

	if err != nil {
		log.Fatalf("Failed to list buckets: %v", err)
	}

	fmt.Println("Buckets:")
	for _, bucket := range buckets {
		fmt.Printf("%s\n", bucket.Name)
	}
}
