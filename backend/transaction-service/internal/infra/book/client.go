package book

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/ilhamagustiawan/library-management-system/backend/transaction-service/internal/domain/errs"
)

const maxResponseBytes = 1 << 20

type TokenSource interface {
	Token(context.Context) (string, error)
}

type Config struct {
	BaseURL string
	Timeout time.Duration
}

type Client struct {
	baseURL string
	tokens  TokenSource
	http    *http.Client
}

type reservation struct {
	TransactionID string `json:"transactionId"`
	BookID        string `json:"bookId"`
	Status        string `json:"status"`
}
type successEnvelope struct {
	Code string      `json:"code"`
	Data reservation `json:"data"`
}
type errorEnvelope struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func NewClient(config Config, tokens TokenSource) (*Client, error) {
	parsed, err := url.Parse(strings.TrimRight(strings.TrimSpace(config.BaseURL), "/"))
	if err != nil || (parsed.Scheme != "http" && parsed.Scheme != "https") || parsed.Host == "" || parsed.RawQuery != "" || parsed.Fragment != "" {
		return nil, fmt.Errorf("BOOK_SERVICE_URL must be an absolute HTTP(S) origin")
	}
	if tokens == nil {
		return nil, fmt.Errorf("Book Service token source is required")
	}
	if config.Timeout <= 0 {
		config.Timeout = 2 * time.Second
	}
	return &Client{baseURL: parsed.String(), tokens: tokens, http: &http.Client{Timeout: config.Timeout}}, nil
}

func (c *Client) Reserve(ctx context.Context, bookID, transactionID string) error {
	request, err := c.request(ctx, http.MethodPut, "/internal/v1/books/"+url.PathEscape(bookID)+"/reservations/"+url.PathEscape(transactionID), nil)
	if err != nil {
		return err
	}
	response, err := c.http.Do(request)
	if err != nil {
		return fmt.Errorf("%w: reserve book stock: %v", errs.ErrDependency, err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK && response.StatusCode != http.StatusCreated {
		return mapResponseError(response)
	}
	var envelope successEnvelope
	if err := decodeStrict(response.Body, &envelope); err != nil {
		return fmt.Errorf("%w: decode Book Service reservation: %v", errs.ErrDependency, err)
	}
	if envelope.Code == "" || envelope.Data.TransactionID != transactionID || envelope.Data.BookID != bookID || envelope.Data.Status != "active" {
		return fmt.Errorf("%w: Book Service returned a mismatched reservation", errs.ErrDependency)
	}
	return nil
}

func (c *Client) Release(ctx context.Context, bookID, transactionID string) error {
	request, err := c.request(ctx, http.MethodDelete, "/internal/v1/books/"+url.PathEscape(bookID)+"/reservations/"+url.PathEscape(transactionID), nil)
	if err != nil {
		return err
	}
	response, err := c.http.Do(request)
	if err != nil {
		return fmt.Errorf("%w: release book reservation: %v", errs.ErrDependency, err)
	}
	defer response.Body.Close()
	if response.StatusCode == http.StatusNoContent || response.StatusCode == http.StatusOK {
		return nil
	}
	return mapResponseError(response)
}

func (c *Client) request(ctx context.Context, method, path string, body io.Reader) (*http.Request, error) {
	token, err := c.tokens.Token(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w: obtain Book Service token: %v", errs.ErrDependency, err)
	}
	request, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, body)
	if err != nil {
		return nil, fmt.Errorf("create Book Service request: %w", err)
	}
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", "Bearer "+token)
	if body != nil {
		request.Header.Set("Content-Type", "application/json")
	}
	return request, nil
}

func mapResponseError(response *http.Response) error {
	var envelope errorEnvelope
	_ = decodeStrict(response.Body, &envelope)
	switch response.StatusCode {
	case http.StatusNotFound:
		return errs.ErrBookNotFound
	case http.StatusConflict:
		if envelope.Code == errs.CodeStockUnavailable {
			return errs.ErrStockUnavailable
		}
		return fmt.Errorf("%w: Book Service rejected the reservation with code %s", errs.ErrDependency, envelope.Code)
	default:
		return fmt.Errorf("%w: Book Service returned HTTP %d", errs.ErrDependency, response.StatusCode)
	}
}

func decodeStrict(reader io.Reader, target any) error {
	decoder := json.NewDecoder(io.LimitReader(reader, maxResponseBytes))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(target); err != nil {
		return err
	}
	if decoder.Decode(&struct{}{}) != io.EOF {
		return fmt.Errorf("response contains multiple JSON values")
	}
	return nil
}
