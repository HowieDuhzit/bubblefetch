package solana

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// TokenInfo represents Solana token metadata and statistics
type TokenInfo struct {
	Address        string
	Name           string
	Symbol         string
	Decimals       int
	Supply         string
	Price          string
	MarketCap      string
	Holders        string
	Description    string
	LogoURI        string
	PriceHistory   []PricePoint
	PriceChange24h float64
	// Bags.fm data (if available)
	Creator        string  // Token creator username/handle
	CreatorPlatform string  // Platform (e.g., Twitter, Discord)
	TotalFees      string  // Total fees collected
	LaunchDate     string  // When token was launched
}

// PricePoint represents a price at a specific time
type PricePoint struct {
	Timestamp int64
	Price     float64
}

// FetchTokenInfo retrieves token data from Solana blockchain
func FetchTokenInfo(contractAddress string) (*TokenInfo, error) {
	return FetchTokenInfoWithAPIKey(contractAddress, "")
}

// FetchTokenInfoWithAPIKey retrieves token data with optional Bags.fm API key
func FetchTokenInfoWithAPIKey(contractAddress, bagsAPIKey string) (*TokenInfo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Try to get token metadata from Solana RPC
	metadata, err := fetchTokenMetadata(ctx, contractAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch token metadata: %w", err)
	}

	// Try to get market data (price, market cap, logo)
	marketData, _ := fetchMarketData(ctx, contractAddress)

	info := &TokenInfo{
		Address:     contractAddress,
		Name:        metadata.Name,
		Symbol:      metadata.Symbol,
		Decimals:    metadata.Decimals,
		Supply:      metadata.Supply,
		Description: metadata.Description,
	}

	if marketData != nil {
		info.Price = marketData.Price
		info.MarketCap = marketData.MarketCap
		info.Holders = marketData.Holders
		info.LogoURI = marketData.LogoURI
		info.PriceChange24h = marketData.PriceChange24h
	}

	// Try to get price history
	priceHistory, _ := fetchPriceHistory(ctx, contractAddress)
	if priceHistory != nil {
		info.PriceHistory = priceHistory
	}

	// Try to get Bags.fm data (enhanced info if available)
	if bagsAPIKey != "" {
		bagsData, _ := fetchBagsData(ctx, contractAddress, bagsAPIKey)
		if bagsData != nil {
			info.Creator = bagsData.Creator
			info.CreatorPlatform = bagsData.CreatorPlatform
			info.TotalFees = bagsData.TotalFees
			info.LaunchDate = bagsData.LaunchDate
		}
	}

	return info, nil
}

type rpcRequest struct {
	Jsonrpc string        `json:"jsonrpc"`
	ID      int           `json:"id"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
}

type rpcResponse struct {
	Result json.RawMessage `json:"result"`
	Error  *rpcError       `json:"error"`
}

type rpcError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type tokenMetadata struct {
	Name        string
	Symbol      string
	Decimals    int
	Supply      string
	Description string
}

type marketData struct {
	Price          string
	MarketCap      string
	Holders        string
	LogoURI        string
	PriceChange24h float64
}

func fetchTokenMetadata(ctx context.Context, address string) (*tokenMetadata, error) {
	// Use public Solana RPC endpoint
	rpcURL := "https://api.mainnet-beta.solana.com"

	// Get token account info
	reqBody := rpcRequest{
		Jsonrpc: "2.0",
		ID:      1,
		Method:  "getAccountInfo",
		Params: []interface{}{
			address,
			map[string]interface{}{
				"encoding": "jsonParsed",
			},
		},
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", rpcURL, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var rpcResp rpcResponse
	if err := json.Unmarshal(respBody, &rpcResp); err != nil {
		return nil, err
	}

	if rpcResp.Error != nil {
		return nil, fmt.Errorf("RPC error: %s", rpcResp.Error.Message)
	}

	// Parse account info
	var accountInfo struct {
		Value struct {
			Data struct {
				Parsed struct {
					Info struct {
						Decimals   int    `json:"decimals"`
						Supply     string `json:"supply"`
						MintAuthority string `json:"mintAuthority"`
					} `json:"info"`
					Type string `json:"type"`
				} `json:"parsed"`
			} `json:"data"`
		} `json:"value"`
	}

	if err := json.Unmarshal(rpcResp.Result, &accountInfo); err != nil {
		return nil, err
	}

	metadata := &tokenMetadata{
		Decimals: accountInfo.Value.Data.Parsed.Info.Decimals,
		Supply:   accountInfo.Value.Data.Parsed.Info.Supply,
	}

	// Try to fetch metadata from common metadata programs
	metadataAccount, err := fetchMetadataAccount(ctx, address)
	if err == nil && metadataAccount != nil {
		metadata.Name = metadataAccount.Name
		metadata.Symbol = metadataAccount.Symbol
		metadata.Description = metadataAccount.Description
	}

	// If name/symbol not found, use address as fallback
	if metadata.Name == "" {
		metadata.Name = "Unknown Token"
	}
	if metadata.Symbol == "" {
		metadata.Symbol = address[:8] + "..."
	}

	return metadata, nil
}

type metadataAccount struct {
	Name        string
	Symbol      string
	Description string
	URI         string
	LogoURI     string
}

func fetchMetadataAccount(ctx context.Context, mintAddress string) (*metadataAccount, error) {
	// Try to fetch from DexScreener API for token info
	url := fmt.Sprintf("https://api.dexscreener.com/latest/dex/tokens/%s", mintAddress)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result struct {
		Pairs []struct {
			BaseToken struct {
				Name    string `json:"name"`
				Symbol  string `json:"symbol"`
				Address string `json:"address"`
			} `json:"baseToken"`
			Info struct {
				ImageURL string `json:"imageUrl"`
			} `json:"info"`
		} `json:"pairs"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	if len(result.Pairs) == 0 {
		return nil, fmt.Errorf("no pairs found")
	}

	logoURI := result.Pairs[0].Info.ImageURL
	if logoURI == "" {
		// Try to get from token metadata URI
		logoURI = fmt.Sprintf("https://raw.githubusercontent.com/solana-labs/token-list/main/assets/mainnet/%s/logo.png", mintAddress)
	}

	return &metadataAccount{
		Name:    result.Pairs[0].BaseToken.Name,
		Symbol:  result.Pairs[0].BaseToken.Symbol,
		LogoURI: logoURI,
	}, nil
}

func fetchPriceHistory(ctx context.Context, address string) ([]PricePoint, error) {
	// Fetch 24h price history from DexScreener
	url := fmt.Sprintf("https://api.dexscreener.com/latest/dex/tokens/%s", address)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, nil
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil
	}

	var result struct {
		Pairs []struct {
			PriceUsd    string `json:"priceUsd"`
			PriceChange struct {
				H24 float64 `json:"h24"`
				H6  float64 `json:"h6"`
				H1  float64 `json:"h1"`
				M5  float64 `json:"m5"`
			} `json:"priceChange"`
		} `json:"pairs"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, nil
	}

	if len(result.Pairs) == 0 {
		return nil, nil
	}

	// Parse current price
	currentPriceStr := result.Pairs[0].PriceUsd
	if currentPriceStr == "" {
		return nil, nil
	}

	currentPrice := parseFloat(currentPriceStr)
	if currentPrice == 0 {
		return nil, nil
	}

	priceChanges := result.Pairs[0].PriceChange

	// Generate price history from current price and price changes
	// We'll create 40 data points going back 24 hours
	priceHistory := make([]PricePoint, 40)
	now := time.Now().Unix()

	// Calculate prices at different time intervals
	// Working backwards from current price using percentage changes

	for i := 0; i < 40; i++ {
		hoursAgo := float64(24) * (float64(39-i) / 39.0)
		timestamp := now - int64(hoursAgo*3600)

		var price float64
		if hoursAgo <= 0.083 { // Last 5 minutes
			changePercent := priceChanges.M5 * (hoursAgo / 0.083)
			price = currentPrice / (1 + changePercent/100)
		} else if hoursAgo <= 1 { // Last hour
			changePercent := priceChanges.H1 * (hoursAgo / 1.0)
			price = currentPrice / (1 + changePercent/100)
		} else if hoursAgo <= 6 { // Last 6 hours
			changePercent := priceChanges.H6 * (hoursAgo / 6.0)
			price = currentPrice / (1 + changePercent/100)
		} else { // Last 24 hours
			changePercent := priceChanges.H24 * (hoursAgo / 24.0)
			price = currentPrice / (1 + changePercent/100)
		}

		priceHistory[i] = PricePoint{
			Timestamp: timestamp,
			Price:     price,
		}
	}

	return priceHistory, nil
}

func parseFloat(s string) float64 {
	var f float64
	fmt.Sscanf(s, "%f", &f)
	return f
}

func fetchMarketData(ctx context.Context, address string) (*marketData, error) {
	// Try to fetch from DexScreener API
	url := fmt.Sprintf("https://api.dexscreener.com/latest/dex/tokens/%s", address)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, nil // Return nil without error for missing market data
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil
	}

	var result struct {
		Pairs []struct {
			PriceUsd    string  `json:"priceUsd"`
			MarketCap   float64 `json:"marketCap"`
			Fdv         float64 `json:"fdv"`
			Liquidity   struct {
				Usd float64 `json:"usd"`
			} `json:"liquidity"`
			Info struct {
				ImageURL string `json:"imageUrl"`
			} `json:"info"`
			PriceChange struct {
				H24 float64 `json:"h24"`
			} `json:"priceChange"`
		} `json:"pairs"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, nil
	}

	if len(result.Pairs) == 0 {
		return nil, nil
	}

	pair := result.Pairs[0]

	// Format market cap
	marketCapStr := ""
	marketCapValue := pair.MarketCap
	if marketCapValue == 0 {
		marketCapValue = pair.Fdv
	}
	if marketCapValue > 0 {
		marketCapStr = formatLargeNumber(marketCapValue)
	}

	return &marketData{
		Price:          fmt.Sprintf("$%s", pair.PriceUsd),
		MarketCap:      marketCapStr,
		Holders:        "N/A", // DexScreener doesn't provide holder count
		LogoURI:        pair.Info.ImageURL,
		PriceChange24h: pair.PriceChange.H24,
	}, nil
}

// formatLargeNumber formats large numbers with K, M, B, T suffixes
func formatLargeNumber(n float64) string {
	if n >= 1e12 {
		return fmt.Sprintf("$%.2fT", n/1e12)
	} else if n >= 1e9 {
		return fmt.Sprintf("$%.2fB", n/1e9)
	} else if n >= 1e6 {
		return fmt.Sprintf("$%.2fM", n/1e6)
	} else if n >= 1e3 {
		return fmt.Sprintf("$%.2fK", n/1e3)
	}
	return fmt.Sprintf("$%.2f", n)
}

// bagsData represents additional token info from Bags.fm API
type bagsData struct {
	Creator         string
	CreatorPlatform string
	TotalFees       string
	LaunchDate      string
}

// fetchBagsData retrieves additional token information from Bags.fm API
func fetchBagsData(ctx context.Context, mintAddress, apiKey string) (*bagsData, error) {
	const baseURL = "https://public-api-v2.bags.fm/api/v1"

	// Try to fetch token launch creators (public analytics endpoint)
	creatorData, _ := fetchBagsCreators(ctx, baseURL, mintAddress, apiKey)

	// Try to fetch token lifetime fees (public analytics endpoint)
	feesData, _ := fetchBagsFees(ctx, baseURL, mintAddress, apiKey)

	// If we got no data from either endpoint, return nil
	if creatorData == nil && feesData == nil {
		return nil, nil
	}

	data := &bagsData{}

	if creatorData != nil {
		data.Creator = creatorData.Creator
		data.CreatorPlatform = creatorData.Platform
		data.LaunchDate = creatorData.LaunchDate
	}

	if feesData != nil {
		data.TotalFees = feesData.TotalFees
	}

	return data, nil
}

type bagsCreatorResponse struct {
	Creator    string
	Platform   string
	LaunchDate string
}

func fetchBagsCreators(ctx context.Context, baseURL, mintAddress, apiKey string) (*bagsCreatorResponse, error) {
	// Endpoint: GET /token-launch/creator/v3?tokenMint={mint}
	url := fmt.Sprintf("%s/token-launch/creator/v3?tokenMint=%s", baseURL, mintAddress)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	// Add API key (required)
	if apiKey != "" {
		req.Header.Set("x-api-key", apiKey)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// If 404 or other error, token wasn't launched via Bags.fm
	if resp.StatusCode != 200 {
		return nil, nil
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Parse Bags.fm response format: {"success": true, "response": [...]}
	var result struct {
		Success  bool `json:"success"`
		Response []struct {
			Username         string `json:"username"`
			Provider         string `json:"provider"`
			BagsUsername     string `json:"bagsUsername"`
			TwitterUsername  string `json:"twitterUsername"`
			ProviderUsername string `json:"providerUsername"`
		} `json:"response"`
		Error string `json:"error"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	if !result.Success {
		return nil, nil
	}

	// Extract creator info
	if len(result.Response) == 0 {
		return nil, nil
	}

	creator := result.Response[0]

	// Use provider username (e.g., Twitter handle)
	username := creator.ProviderUsername
	if username == "" {
		username = creator.Username
	}

	// Capitalize platform name
	platform := creator.Provider
	if platform == "twitter" {
		platform = "Twitter"
	} else if platform == "discord" {
		platform = "Discord"
	} else if platform != "" {
		// Capitalize first letter
		platform = strings.ToUpper(platform[:1]) + platform[1:]
	} else {
		platform = "Bags.fm"
	}

	return &bagsCreatorResponse{
		Creator:    username,
		Platform:   platform,
		LaunchDate: "", // Launch date not in creator endpoint
	}, nil
}

type bagsFeesResponse struct {
	TotalFees string
}

func fetchBagsFees(ctx context.Context, baseURL, mintAddress, apiKey string) (*bagsFeesResponse, error) {
	// Endpoint: GET /token-launch/lifetime-fees?tokenMint={mint}
	url := fmt.Sprintf("%s/token-launch/lifetime-fees?tokenMint=%s", baseURL, mintAddress)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	// Add API key (required)
	if apiKey != "" {
		req.Header.Set("x-api-key", apiKey)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// If not 200, token wasn't launched via Bags.fm or other error
	if resp.StatusCode != 200 {
		return nil, nil
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Parse Bags.fm response format: {"success": true, "response": "863007244"}
	// Response is lamports as a string
	var result struct {
		Success  bool   `json:"success"`
		Response string `json:"response"`
		Error    string `json:"error"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	if !result.Success {
		return nil, nil
	}

	// Convert lamports to SOL (1 SOL = 1,000,000,000 lamports)
	lamports := parseFloat(result.Response)
	sol := lamports / 1e9

	// Format as currency
	totalFees := formatLargeNumber(sol)
	// Replace $ with SOL symbol since these are in SOL, not USD
	totalFees = fmt.Sprintf("◎%.2f", sol)
	if sol >= 1000 {
		totalFees = fmt.Sprintf("◎%.2fK", sol/1000)
	} else if sol >= 1000000 {
		totalFees = fmt.Sprintf("◎%.2fM", sol/1000000)
	}

	return &bagsFeesResponse{
		TotalFees: totalFees,
	}, nil
}
