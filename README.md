# algolymp
*A collection of useful CLI tools for managing Polygon and Ejudge.*

## Build
```bash
sudo apt install go
make build
```

## Config

Put your config file in `~/.config/algolymp/config.json`.

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
- `-i` - contest id (required)
- `-c` - path to config file (sections: `ejudge`)
- `-v` - extended output from Ejudge responses.

## Examples
* `ejik --help`
* `ejik -c config.json -v -i 47103`
* `ejik -i 40507`
