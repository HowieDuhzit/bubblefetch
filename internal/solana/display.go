package solana

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/howieduhzit/bubblefetch/internal/ui/theme"
	"gopkg.in/yaml.v3"
)

// Display renders token information with theme styling (TUI layout)
func Display(info *TokenInfo, themeName string) (string, error) {
	// Load theme
	thm, err := theme.Load(themeName)
	if err != nil {
		// Fallback to default theme
		thm, _ = theme.Load("default")
	}

	styles := thm.GetStyles()

	var content strings.Builder

	// Get ASCII art (token logo or default)
	var asciiArt string
	if thm.Layout.ShowASCII {
		logo, err := FetchAndConvertLogo(info.LogoURI)
		if err != nil {
			logo = getDefaultTokenASCII()
		}
		asciiArt = styles.ASCII.Render(logo)
	}

	// Build token info lines
	var moduleLines []string

	// Title
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color(thm.Colors.Primary))

	moduleLines = append(moduleLines, titleStyle.Render(fmt.Sprintf("🪙 %s (%s)", info.Name, info.Symbol)))
	moduleLines = append(moduleLines, "") // Empty line

	// Token fields
	fields := []struct {
		label string
		value string
	}{
		{"Contract", truncateAddress(info.Address)},
		{"Symbol", info.Symbol},
		{"Decimals", fmt.Sprintf("%d", info.Decimals)},
		{"Supply", formatSupply(info.Supply, info.Decimals)},
	}

	// Add market data if available
	if info.Price != "" {
		fields = append(fields, struct {
			label string
			value string
		}{"Price", info.Price})
	}

	if info.MarketCap != "" {
		fields = append(fields, struct {
			label string
			value string
		}{"Market Cap", formatMarketCap(info.MarketCap)})
	}

	// Add Bags.fm data if available
	if info.Creator != "" {
		creatorInfo := info.Creator
		if info.CreatorPlatform != "" {
			creatorInfo = fmt.Sprintf("%s (%s)", info.Creator, info.CreatorPlatform)
		}
		fields = append(fields, struct {
			label string
			value string
		}{"Creator", creatorInfo})
	}

	if info.LaunchDate != "" {
		fields = append(fields, struct {
			label string
			value string
		}{"Launched", info.LaunchDate})
	}

	if info.TotalFees != "" {
		fields = append(fields, struct {
			label string
			value string
		}{"Total Fees", info.TotalFees})
	}

	// Render fields
	for _, field := range fields {
		if field.value == "" {
			continue
		}

		line := lipgloss.JoinHorizontal(
			lipgloss.Left,
			styles.Label.Render(field.label),
			styles.Separator.Render(": "),
			styles.Value.Render(field.value),
		)

		moduleLines = append(moduleLines, line)
	}

	// Add price chart if available
	if len(info.PriceHistory) > 0 {
		moduleLines = append(moduleLines, "")

		// Show price change percentage
		var changeLabel string
		if info.PriceChange24h > 0 {
			changeLabel = fmt.Sprintf("Price Chart (24h) 📈 +%.2f%%", info.PriceChange24h)
		} else if info.PriceChange24h < 0 {
			changeLabel = fmt.Sprintf("Price Chart (24h) 📉 %.2f%%", info.PriceChange24h)
		} else {
			changeLabel = "Price Chart (24h)"
		}

		chartLabel := lipgloss.NewStyle().
			Foreground(lipgloss.Color(thm.Colors.Secondary)).
			Render(changeLabel)
		moduleLines = append(moduleLines, chartLabel)

		chart := GeneratePriceChart(info.PriceHistory, 40)
		chartStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color(thm.Colors.Accent))
		moduleLines = append(moduleLines, chartStyle.Render(chart))
	}

	moduleContent := strings.Join(moduleLines, "\n")

	// Join ASCII art and content horizontally (like system info)
	if thm.Layout.ShowASCII && asciiArt != "" {
		content.WriteString(lipgloss.JoinHorizontal(
			lipgloss.Top,
			asciiArt,
			strings.Repeat(" ", thm.Layout.Padding),
			moduleContent,
		))
	} else {
		content.WriteString(moduleContent)
	}

	// Add border
	result := styles.Container.Render(content.String())
	if thm.Layout.BorderStyle != "none" && thm.Layout.BorderStyle != "" {
		result = styles.Border.Render(result)
	}

	return result, nil
}

// truncateAddress shortens a Solana address for display
func truncateAddress(address string) string {
	if len(address) <= 16 {
		return address
	}
	return address[:8] + "..." + address[len(address)-8:]
}

// formatSupply formats token supply with appropriate units
func formatSupply(supply string, decimals int) string {
	// For now, just return the raw supply
	// In production, you'd want to divide by 10^decimals and format with units (K, M, B, T)
	if len(supply) > 15 {
		return supply[:12] + "..."
	}
	return supply
}

// formatMarketCap formats market cap value
func formatMarketCap(marketCap string) string {
	// Remove $ and format
	if len(marketCap) > 1 && marketCap[0] == '$' {
		return marketCap
	}
	return marketCap
}

// DisplayJSON renders token information as JSON
func DisplayJSON(info *TokenInfo) (string, error) {
	data := map[string]interface{}{
		"contract_address": info.Address,
		"name":             info.Name,
		"symbol":           info.Symbol,
		"decimals":         info.Decimals,
		"supply":           info.Supply,
	}

	if info.Price != "" {
		data["price"] = info.Price
	}

	if info.MarketCap != "" {
		data["market_cap"] = info.MarketCap
	}

	if info.Holders != "" {
		data["holders"] = info.Holders
	}

	if info.Description != "" {
		data["description"] = info.Description
	}

	if info.LogoURI != "" {
		data["logo_uri"] = info.LogoURI
	}

	if info.PriceChange24h != 0 {
		data["price_change_24h"] = info.PriceChange24h
	}

	// Add Bags.fm data if available
	if info.Creator != "" {
		data["creator"] = info.Creator
	}

	if info.CreatorPlatform != "" {
		data["creator_platform"] = info.CreatorPlatform
	}

	if info.LaunchDate != "" {
		data["launch_date"] = info.LaunchDate
	}

	if info.TotalFees != "" {
		data["total_fees"] = info.TotalFees
	}

	// Format as pretty JSON
	jsonBytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return "", err
	}

	return string(jsonBytes), nil
}

// DisplayYAML renders token information as YAML
func DisplayYAML(info *TokenInfo) (string, error) {
	data := map[string]interface{}{
		"contract_address": info.Address,
		"name":             info.Name,
		"symbol":           info.Symbol,
		"decimals":         info.Decimals,
		"supply":           info.Supply,
	}

	if info.Price != "" {
		data["price"] = info.Price
	}

	if info.MarketCap != "" {
		data["market_cap"] = info.MarketCap
	}

	if info.Holders != "" {
		data["holders"] = info.Holders
	}

	if info.Description != "" {
		data["description"] = info.Description
	}

	if info.LogoURI != "" {
		data["logo_uri"] = info.LogoURI
	}

	// Add Bags.fm data if available
	if info.Creator != "" {
		data["creator"] = info.Creator
	}

	if info.CreatorPlatform != "" {
		data["creator_platform"] = info.CreatorPlatform
	}

	if info.LaunchDate != "" {
		data["launch_date"] = info.LaunchDate
	}

	if info.TotalFees != "" {
		data["total_fees"] = info.TotalFees
	}

	yamlBytes, err := yaml.Marshal(data)
	if err != nil {
		return "", err
	}

	return string(yamlBytes), nil
}
