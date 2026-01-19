# Export Formats & JSON Schema

Bubblefetch exports structured data via `--export json|yaml|text`. JSON/YAML share the same schema (YAML is a direct serialization of JSON).

## JSON Schema (v1)

The JSON output is a single object with these fields:

```json
{
  "OS": "string",
  "Kernel": "string",
  "Hostname": "string",
  "Uptime": "string",
  "CPU": "string",
  "Memory": { "Used": 0, "Total": 0 },
  "Disk": { "Used": 0, "Total": 0 },
  "Shell": "string",
  "Terminal": "string",
  "Resolution": "string",
  "DE": "string",
  "WM": "string",
  "Theme": "string",
  "Icons": "string",
  "GPU": ["string"],
  "Network": [{ "Interface": "string", "IPv4": "string", "IPv6": "string", "MAC": "string" }],
  "Battery": { "Present": false, "Percentage": 0, "IsCharging": false, "TimeRemain": "string" },
  "LocalIP": "string",
  "PublicIP": "string"
}
```

Notes:
- Missing data may appear as empty strings, zero values, or empty arrays.
- Field names are stable and match `internal/collectors.SystemInfo`.
- JSON output does not include theme styling or ANSI colors.

## CLI Examples

```bash
bf --export json > system.json
bf --export yaml > system.yaml
bf --export text > system.txt
```

If you rely on this schema, pin a bubblefetch version in automation and watch `docs/CHANGELOG.md` for schema changes.
