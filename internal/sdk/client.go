package sdk

import (
	"context"
	"fmt"

	"github.com/patrickcping/pingone-go-sdk-v2/pingone"
)

type Client struct {
	API *pingone.Client
}

type Config struct {
	ClientID             string
	ClientSecret         string
	EnvironmentID        string
	AccessToken          string
	Region               string
	APIHostnameOverride  *string
	AuthHostnameOverride *string
	ProxyURL             *string
}

func (c *Config) APIClient(ctx context.Context, version string) (*Client, error) {

	userAgent := fmt.Sprintf("pingtools PingOne-CleanConfig/%s go", version)

	config := &pingone.Config{
		ClientID:             &c.ClientID,
		ClientSecret:         &c.ClientSecret,
		EnvironmentID:        &c.EnvironmentID,
		AccessToken:          &c.AccessToken,
		Region:               c.Region,
		APIHostnameOverride:  c.APIHostnameOverride,
		AuthHostnameOverride: c.AuthHostnameOverride,
		UserAgentOverride:    &userAgent,
		ProxyURL:             c.ProxyURL,
	}

	client, err := config.APIClient(ctx)
	if err != nil {
		return nil, err
	}

	cliClient := &Client{
		API: client,
	}

	return cliClient, nil
}
