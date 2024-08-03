package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/joho/godotenv"
	"github.com/lll-lll-lll/g3drive"
	"google.golang.org/api/drive/v3"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
}
func main() {
	ctx := context.Background()
	if err := run(ctx, os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

}

func run(ctx context.Context, fileNames []string) error {
	var wg sync.WaitGroup
	wg.Add(len(fileNames))
	errc := make(chan error, len(fileNames))

	srv, err := drive.NewService(ctx)
	if err != nil {
		return fmt.Errorf("Unable to retrieve Drive client: %w", err)
	}
	client := g3drive.New(srv)
	for _, fileName := range fileNames {
		fileName := fileName
		go func(fileName string) error {
			defer wg.Done()
			g3f, err := g3drive.Parse(fileName)
			if err != nil {
				errc <- fmt.Errorf("Failed to parse file %s: %w", fileName, err)
				return fmt.Errorf("%w", err)
			}
			if err := g3drive.Upload(ctx, client, g3f); err != nil {
				errc <- fmt.Errorf("Failed to upload file %s: %w", fileName, err)

				return fmt.Errorf("%w", err)
			}
			return nil
		}(fileName)
	}
	wg.Wait()
	close(errc)
	if len(errc) > 0 {
		return <-errc
	}
	return nil
}
