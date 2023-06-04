package main

import (
	"context"
	"fmt"
	"log"
	"os"

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
	if err := run(ctx, os.Args[1]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

}

func run(ctx context.Context, fileName string) error {
	srv, err := drive.NewService(ctx)
	if err != nil {
		return fmt.Errorf("Unable to retrieve Drive client: %v", err)
	}
	client := g3drive.New(srv)
	g3f, err := g3drive.Parse(fileName)
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	if err := g3drive.Upload(ctx, client, g3f); err != nil {
		return fmt.Errorf("%w", err)
	}
	return nil
}
