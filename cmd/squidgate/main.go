package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/jacobalberty/authgate/internal/client/peers"
)

func main() {
	ctx := context.Background()

	if err := start(ctx, os.Stdin, os.Args[1:]); err != nil {
		panic(err)
	}
}

func start(ctx context.Context, in io.Reader, args []string) error {
	var (
		filename = "squidgate.json"
	)

	if len(args) > 0 {
		filename = args[0]
	}

	config, err := loadConfig(filename)
	if err != nil {
		return err
	}

	nb, err := peers.New(config.RootEndpoint, config.Token)

	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(in)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		parts := strings.Fields(line)
		if len(parts) != 2 {
			fmt.Println("ERR")
			continue
		}

		group, ip := parts[0], parts[1]
		ok, err := nb.IsPeerInGroup(ctx, group, ip)
		if err != nil {
			fmt.Fprintln(os.Stderr, "error checking group:", err)
			fmt.Println("ERR")
			continue
		}

		if ok {
			fmt.Println("OK")
		} else {
			fmt.Println("ERR")
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "helper read error:", err)
	}
	return nil
}

// loadConfig loads the configuration for the NetBird API client.
func loadConfig(filen string) (*Config, error) {
	// Load our config from the json in the file
	if filen == "" {
		return nil, fmt.Errorf("config file name is empty")
	}
	if _, err := os.Stat(filen); os.IsNotExist(err) {
		return nil, fmt.Errorf("config file %s does not exist", filen)
	}

	file, err := os.Open(filen)
	if err != nil {
		return nil, fmt.Errorf("error opening config file %s: %w", filen, err)
	}
	defer file.Close()
	var config Config
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return nil, fmt.Errorf("error decoding config file %s: %w", filen, err)
	}

	return &config, nil

}

type Config struct {
	RootEndpoint string `json:"root_endpoint"`
	Token        string `json:"token"`
}
