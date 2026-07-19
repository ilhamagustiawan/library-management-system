package cmd

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/spf13/cobra"

	"github.com/ilhamagustiawan/library-management-system/backend/auth-service/internal/config"
	authinfra "github.com/ilhamagustiawan/library-management-system/backend/auth-service/internal/infra/auth"
	infraDB "github.com/ilhamagustiawan/library-management-system/backend/auth-service/internal/infra/db"
	oauthClientRepo "github.com/ilhamagustiawan/library-management-system/backend/auth-service/internal/infra/db/repository/oauthclient"
)

var createClientCmd = &cobra.Command{
	Use:   "create-client",
	Short: "Create a confidential OAuth client and print its secret once",
	RunE:  runCreateClient,
}

func runCreateClient(cmd *cobra.Command, _ []string) error {
	name, _ := cmd.Flags().GetString("name")
	clientID, _ := cmd.Flags().GetString("id")
	redirectURI, _ := cmd.Flags().GetString("redirect-uri")
	scopes, _ := cmd.Flags().GetString("scopes")
	kindValue, _ := cmd.Flags().GetString("kind")
	if strings.TrimSpace(name) == "" {
		return fmt.Errorf("--name is required")
	}
	if err := validateClientProvisioning(kindValue, redirectURI, scopes); err != nil {
		return err
	}
	kind := oauthClientRepo.Kind(kindValue)
	if clientID == "" {
		clientID = uuid.NewString()
	}

	cfg, err := config.Load()
	if err != nil {
		return err
	}
	db, err := infraDB.Connect(context.Background(), infraDB.Config{
		DSN: cfg.Database.DSN, MaxOpenConns: cfg.Database.MaxOpenConns, MaxIdleConns: cfg.Database.MaxIdleConns,
		ConnMaxLifetime: cfg.Database.ConnMaxLifetime, ConnMaxIdleTime: cfg.Database.ConnMaxIdleTime,
	})
	if err != nil {
		return err
	}
	defer db.Close()

	secret, err := authinfra.NewTokenGenerator(32).Generate()
	if err != nil {
		return fmt.Errorf("generate client secret: %w", err)
	}
	secretHash, err := authinfra.NewPasswordHasher(cfg.Auth.BcryptCost).Hash(secret)
	if err != nil {
		return fmt.Errorf("hash client secret: %w", err)
	}
	client := &oauthClientRepo.Client{
		ID: clientID, Name: strings.TrimSpace(name), Kind: kind,
		SecretHash: secretHash, RedirectURI: redirectURI,
		Scopes: strings.Fields(scopes), Public: false,
		UserID: sql.NullString{}, CreatedAt: time.Now().UTC(),
	}
	if err := oauthClientRepo.NewStore(db).Create(cmd.Context(), client); err != nil {
		return fmt.Errorf("create OAuth client: %w", err)
	}
	return json.NewEncoder(os.Stdout).Encode(map[string]any{
		"clientId": clientID, "clientSecret": secret, "redirectUri": redirectURI, "scopes": client.ScopeList(),
	})
}

func init() {
	createClientCmd.Flags().String("name", "", "human-readable client name")
	createClientCmd.Flags().String("id", "", "client ID (generated when omitted)")
	createClientCmd.Flags().String("redirect-uri", "", "exact OAuth callback URL")
	createClientCmd.Flags().String("kind", string(oauthClientRepo.KindAuthorizationCode), "authorization_code, client_credentials, or resource_server")
	createClientCmd.Flags().String("scopes", "loans:borrow:self loans:return:self transactions:read:self books:read", "space-separated allowed scopes")
}

func validateClientProvisioning(kindValue, redirectURI, scopes string) error {
	kind := oauthClientRepo.Kind(kindValue)
	switch kind {
	case oauthClientRepo.KindAuthorizationCode:
		if err := validateRedirectURI(redirectURI); err != nil {
			return err
		}
	case oauthClientRepo.KindClientCredentials:
		if redirectURI != "" {
			return fmt.Errorf("--redirect-uri must be empty for client_credentials")
		}
		if len(strings.Fields(scopes)) == 0 {
			return fmt.Errorf("--scopes must contain at least one scope for client_credentials")
		}
	case oauthClientRepo.KindResourceServer:
		if redirectURI != "" || len(strings.Fields(scopes)) != 0 {
			return fmt.Errorf("resource_server requires empty --redirect-uri and --scopes")
		}
	default:
		return fmt.Errorf("--kind must be authorization_code, client_credentials, or resource_server")
	}
	return nil
}

func validateRedirectURI(raw string) error {
	parsed, err := url.Parse(raw)
	if err != nil || parsed.Scheme == "" || parsed.Host == "" || parsed.Fragment != "" {
		return fmt.Errorf("--redirect-uri must be an absolute URL without a fragment")
	}
	loopback := parsed.Hostname() == "localhost" || parsed.Hostname() == "127.0.0.1"
	if parsed.Scheme != "https" && !(parsed.Scheme == "http" && loopback) {
		return fmt.Errorf("--redirect-uri must use HTTPS except on localhost")
	}
	return nil
}
