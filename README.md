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
	},
	"polygon": {
		"url": "https://polygon.codeforces.com",
		"apiKey": "<key>",
		"apiSecret": "<secret>"
	}
}
```

## blanka
*Create Ejuge contest from template.*

### About

1. Create contest with id from template;
2. Commit changes;
3. *(Optional)* Open contest xml config for editing.

Useful before running `polygon-to-ejudge`.

### Flags
- `-i` - new contest id (required)
- `-t` - template contest id (required)

### Config
- `ejudge.url`
- `ejudge.login`
- `ejudge.password`

### Examples
* `blanka -i 51011 -t 51000`


## ejik
*Refresh Ejudge contest by id.*

### About

1. Commit changes;
2. Check contest settings;
3. Reload config files.

Useful after running `polygon-to-ejudge`.

### Flags
- `-i` - contest id (required)
- `-v` - extended output from Ejudge responses.

### Config
- `ejudge.url`
- `ejudge.login`
- `ejudge.password`

### Examples
* `ejik --help`
* `ejik -v -i 47103`
* `ejik -i 40507`


## valeria
*Build valuer + scorer using Polygon API.*

### About

1. Get problem tests and groups;
2. Build and commit `valuer.cfg` (in Ejudge format);
3. Build and print `scoring.tex` (in Moscow summer olympiad school format).

Not very fast now, waiting for `absentInput` parameter in Polygon API.

### Flags
- `-i` - problem id (required)

### Config
- `polygon.url`
- `polygon.apiKey`
- `polygon.apiSecret`

### Examples
* `valeria --help`
* `valeria -i 288808`
* `valeria -i 318511 > scoring.tex`
* `valeria -i 318882 | bat -l tex`
