//go:generate mockgen -source=$GOFILE -destination=../../mock/mock_geocoding_client.go -package=mock

package external

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

const (
	defaultGeocodingURL = "https://msearch.gsi.go.jp/address-search/AddressSearch"
	defaultTimeout      = 5 * time.Second
	maxRetries          = 1
)

// Coordinates は緯度経度を表す
type Coordinates struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

// GeocodingResult はジオコーディング結果を表す
type GeocodingResult struct {
	Coordinates Coordinates `json:"coordinates"`
	AddressCode string      `json:"address_code"`
}

// geocodingResponse は国土地理院APIのレスポンス
type geocodingResponse struct {
	Geometry struct {
		Coordinates []float64 `json:"coordinates"` // [lng, lat]
	} `json:"geometry"`
	Properties struct {
		AddressCode string `json:"addressCode"`
	} `json:"properties"`
}

// GeocodingClient はジオコーディングAPIクライアントのインターフェース
type GeocodingClient interface {
	// Geocode は住所から緯度経度を取得する
	Geocode(ctx context.Context, address string) (*GeocodingResult, error)
}

// geocodingClientImpl はGeocodingClientの実装
type geocodingClientImpl struct {
	baseURL    string
	httpClient *http.Client
}

// NewGeocodingClient はGeocodingClientの新しいインスタンスを作成する
func NewGeocodingClient() GeocodingClient {
	return &geocodingClientImpl{
		baseURL: defaultGeocodingURL,
		httpClient: &http.Client{
			Timeout: defaultTimeout,
		},
	}
}

// NewGeocodingClientWithURL はカスタムURLでGeocodingClientを作成する（テスト用）
func NewGeocodingClientWithURL(baseURL string) GeocodingClient {
	return &geocodingClientImpl{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: defaultTimeout,
		},
	}
}

// Geocode は住所から緯度経度を取得する
func (c *geocodingClientImpl) Geocode(ctx context.Context, address string) (*GeocodingResult, error) {
	var lastErr error

	for i := 0; i <= maxRetries; i++ {
		result, err := c.doGeocode(ctx, address)
		if err == nil {
			return result, nil
		}
		lastErr = err

		// コンテキストがキャンセルされた場合はリトライしない
		if ctx.Err() != nil {
			return nil, fmt.Errorf("geocoding cancelled: %w", ctx.Err())
		}

		// リトライ前に少し待機
		if i < maxRetries {
			time.Sleep(100 * time.Millisecond)
		}
	}

	return nil, fmt.Errorf("geocoding failed after retries: %w", lastErr)
}

func (c *geocodingClientImpl) doGeocode(ctx context.Context, address string) (*GeocodingResult, error) {
	// URLを構築
	reqURL := fmt.Sprintf("%s?q=%s", c.baseURL, url.QueryEscape(address))

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var results []geocodingResponse
	if err := json.NewDecoder(resp.Body).Decode(&results); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if len(results) == 0 {
		return nil, fmt.Errorf("no results found for address: %s", address)
	}

	// 最初の結果を使用
	first := results[0]
	if len(first.Geometry.Coordinates) < 2 {
		return nil, fmt.Errorf("invalid coordinates in response")
	}

	return &GeocodingResult{
		Coordinates: Coordinates{
			Lng: first.Geometry.Coordinates[0],
			Lat: first.Geometry.Coordinates[1],
		},
		AddressCode: first.Properties.AddressCode,
	}, nil
}
