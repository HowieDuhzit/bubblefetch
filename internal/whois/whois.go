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
func Lookup(target string) (string, error) {
	target = strings.TrimSpace(target)
	if target == "" {
		return "", fmt.Errorf("whois target is empty")
	}

	var b strings.Builder
	b.WriteString("Target: ")
	b.WriteString(target)
	b.WriteString("\n\nWHOIS\n")

	whoisText, err := lookupWhois(target)
	if err != nil {
		return "", err
	}
	b.WriteString(whoisText)
	b.WriteString("\n\nDNS\n")
	b.WriteString(lookupDNS(target))

	return b.String(), nil
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

func lookupDNS(target string) string {
	if ip := net.ParseIP(target); ip != nil {
		hosts, err := net.LookupAddr(target)
		if err != nil || len(hosts) == 0 {
			return "PTR: (none)\n"
		}
		return "PTR: " + strings.Join(hosts, ", ") + "\n"
	}

	var b strings.Builder

	if ips, err := net.LookupIP(target); err == nil && len(ips) > 0 {
		var a []string
		var aaaa []string
		for _, ip := range ips {
			if ip.To4() != nil {
				a = append(a, ip.String())
			} else {
				aaaa = append(aaaa, ip.String())
			}
		}
		if len(a) > 0 {
			b.WriteString("A: ")
			b.WriteString(strings.Join(a, ", "))
			b.WriteString("\n")
		}
		if len(aaaa) > 0 {
			b.WriteString("AAAA: ")
			b.WriteString(strings.Join(aaaa, ", "))
			b.WriteString("\n")
		}
	} else {
		b.WriteString("A/AAAA: (none)\n")
	}

	if cname, err := net.LookupCNAME(target); err == nil && cname != "" {
		b.WriteString("CNAME: ")
		b.WriteString(cname)
		b.WriteString("\n")
	}

	if mx, err := net.LookupMX(target); err == nil && len(mx) > 0 {
		var records []string
		for _, record := range mx {
			records = append(records, fmt.Sprintf("%d %s", record.Pref, record.Host))
		}
		b.WriteString("MX: ")
		b.WriteString(strings.Join(records, ", "))
		b.WriteString("\n")
	}

	if ns, err := net.LookupNS(target); err == nil && len(ns) > 0 {
		var records []string
		for _, record := range ns {
			records = append(records, record.Host)
		}
		b.WriteString("NS: ")
		b.WriteString(strings.Join(records, ", "))
		b.WriteString("\n")
	}

	if txt, err := net.LookupTXT(target); err == nil && len(txt) > 0 {
		b.WriteString("TXT:\n")
		for _, record := range txt {
			b.WriteString("  - ")
			b.WriteString(record)
			b.WriteString("\n")
		}
	}

	return b.String()
}
