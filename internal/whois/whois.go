package whois

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"time"
)

const (
	whoisPort    = "43"
	whoisTimeout = 8 * time.Second
)

// Lookup performs a WHOIS lookup plus DNS record scan for a target domain.
func Lookup(target string, includeRaw bool) (Result, error) {
	target = strings.TrimSpace(target)
	if target == "" {
		return Result{}, fmt.Errorf("whois target is empty")
	}

	whoisText, err := lookupWhois(target)
	if err != nil {
		return Result{}, err
	}
	parsed := parseWhois(whoisText)
	if !includeRaw {
		parsed.Raw = ""
	}

	return Result{
		Target: target,
		Whois:  parsed,
		DNS:    lookupDNS(target),
	}, nil
}

func lookupWhois(target string) (string, error) {
	if ip := net.ParseIP(target); ip != nil {
		return lookupWhoisChain(target)
	}

	domain := strings.ToLower(target)
	tld := domain
	if idx := strings.LastIndex(domain, "."); idx != -1 {
		tld = domain[idx+1:]
	}

	ref, _ := lookupRefer("whois.iana.org", tld)
	if ref == "" {
		ref = "whois.iana.org"
	}

	resp, err := queryWhois(ref, domain)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(resp), nil
}

func lookupWhoisChain(target string) (string, error) {
	ref, _ := lookupRefer("whois.iana.org", target)
	if ref == "" {
		ref = "whois.iana.org"
	}
	resp, err := queryWhois(ref, target)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(resp), nil
}

func lookupRefer(server, query string) (string, error) {
	resp, err := queryWhois(server, query)
	if err != nil {
		return "", err
	}

	scanner := bufio.NewScanner(strings.NewReader(resp))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		lower := strings.ToLower(line)
		if strings.HasPrefix(lower, "refer:") || strings.HasPrefix(lower, "whois:") {
			fields := strings.Fields(line)
			if len(fields) >= 2 {
				return fields[1], nil
			}
		}
	}

	return "", nil
}

func queryWhois(server, query string) (string, error) {
	conn, err := net.DialTimeout("tcp", net.JoinHostPort(server, whoisPort), whoisTimeout)
	if err != nil {
		return "", fmt.Errorf("whois connect failed: %w", err)
	}
	defer conn.Close()

	_ = conn.SetDeadline(time.Now().Add(whoisTimeout))
	if _, err := fmt.Fprintf(conn, "%s\r\n", query); err != nil {
		return "", fmt.Errorf("whois query failed: %w", err)
	}

	var b strings.Builder
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		b.WriteString(scanner.Text())
		b.WriteString("\n")
	}
	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("whois read failed: %w", err)
	}

	return b.String(), nil
}

func lookupDNS(target string) dnsResult {
	if ip := net.ParseIP(target); ip != nil {
		hosts, err := net.LookupAddr(target)
		if err != nil || len(hosts) == 0 {
			return dnsResult{PTR: []string{}}
		}
		return dnsResult{PTR: hosts}
	}

	result := dnsResult{}

	if ips, err := net.LookupIP(target); err == nil && len(ips) > 0 {
		for _, ip := range ips {
			if ip.To4() != nil {
				result.A = append(result.A, ip.String())
			} else {
				result.AAAA = append(result.AAAA, ip.String())
			}
		}
	}

	if cname, err := net.LookupCNAME(target); err == nil && cname != "" {
		result.CNAME = cname
	}

	if mx, err := net.LookupMX(target); err == nil && len(mx) > 0 {
		for _, record := range mx {
			result.MX = append(result.MX, fmt.Sprintf("%d %s", record.Pref, record.Host))
		}
	}

	if ns, err := net.LookupNS(target); err == nil && len(ns) > 0 {
		for _, record := range ns {
			result.NS = append(result.NS, record.Host)
		}
	}

	if txt, err := net.LookupTXT(target); err == nil && len(txt) > 0 {
		result.TXT = append(result.TXT, txt...)
	}

	return result
}

type whoisField struct {
	Label string
	Value string
}

type whoisParsed struct {
	Fields []whoisField
	Raw    string
}

type dnsResult struct {
	A    []string
	AAAA []string
	CNAME string
	MX   []string
	NS   []string
	TXT  []string
	PTR  []string
}

type Result struct {
	Target string
	Whois  whoisParsed
	DNS    dnsResult
}

func parseWhois(raw string) whoisParsed {
	lines := strings.Split(raw, "\n")
	rawTrim := strings.TrimSpace(raw)

	var registrar string
	var created string
	var updated string
	var expires string
	var domainID string
	var statuses []string
	var nameServers []string

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "%") || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.ToLower(strings.TrimSpace(parts[0]))
		value := strings.TrimSpace(parts[1])
		if value == "" {
			continue
		}

		switch {
		case key == "registrar":
			if registrar == "" {
				registrar = value
			}
		case key == "registry domain id" || key == "domain id":
			if domainID == "" {
				domainID = value
			}
		case key == "creation date" || key == "created" || key == "registered on":
			if created == "" {
				created = value
			}
		case key == "updated date" || key == "last updated" || key == "updated":
			if updated == "" {
				updated = value
			}
		case key == "registry expiry date" || key == "expiration date" || key == "expiry date" || key == "paid-till":
			if expires == "" {
				expires = value
			}
		case key == "domain status" || key == "status":
			statuses = append(statuses, value)
		case key == "name server" || key == "nserver":
			nameServers = append(nameServers, value)
		}
	}

	fields := []whoisField{}
	if registrar != "" {
		fields = append(fields, whoisField{Label: "Registrar", Value: registrar})
	}
	if domainID != "" {
		fields = append(fields, whoisField{Label: "Registry ID", Value: domainID})
	}
	if created != "" {
		fields = append(fields, whoisField{Label: "Created", Value: created})
	}
	if updated != "" {
		fields = append(fields, whoisField{Label: "Updated", Value: updated})
	}
	if expires != "" {
		fields = append(fields, whoisField{Label: "Expires", Value: expires})
	}
	if len(statuses) > 0 {
		fields = append(fields, whoisField{Label: "Status", Value: strings.Join(unique(statuses), ", ")})
	}
	if len(nameServers) > 0 {
		fields = append(fields, whoisField{Label: "Name Servers", Value: strings.Join(unique(nameServers), ", ")})
	}

	return whoisParsed{
		Fields: fields,
		Raw:    rawTrim,
	}
}

func unique(values []string) []string {
	seen := make(map[string]struct{}, len(values))
	result := make([]string, 0, len(values))
	for _, value := range values {
		v := strings.TrimSpace(value)
		if v == "" {
			continue
		}
		if _, ok := seen[v]; ok {
			continue
		}
		seen[v] = struct{}{}
		result = append(result, v)
	}
	return result
}
