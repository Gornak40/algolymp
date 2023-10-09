# algolymp
*A collection of useful CLI tools for managing Polygon and Ejudge.*

## Build
```bash
sudo apt install go
make build
```

## Config
```json
{
	"ejudge": {
		"url": "https://ejudge.algocourses.ru",
		"login": "<login>",
		"password": "<password>"
	}
}
```

# ejik
*Refresh Ejudge contest by id.*

## About

1. Commit changes;
2. Check contest settings;
3. Reload config files.

Useful after running `polygon-to-ejudge`.

## Flags
- `-c` - path to config file (requires: `ejudge`)
- `-v` - extended output from Ejudge responses.
