# algolymp
*A collection of useful CLI tools for managing Polygon and Ejudge.*

## Workflow

| Tool | Description | Ejudge | Polygon | Status |
| --- | --- | :---: | :---: | :---: |
| [blanka](#blanka) | create contest | 🦍 | | ✅ |
| [boban](#boban) | filter runs | 🦍 | | ✅ |
| [casper](#casper) | change visibility | 🦍 | | ✅ |
| [ejik](#ejik) | commit + check + reload | 🦍 | | ✅ |
| [ripper](#ripper) | change runs status | 🦍 | | ✅ |
| [scalp](#scalp) | incremental scoring | | 🦍 | ✅ |
| [valeria](#valeria) | valuer.cfg + tex scoring | | 🦍 | ✅ |
| [wooda](#wooda) | glob problem files upload | | 🦍 | 🧑‍💻 |
| ⚙️ | move json config to ini | | | 🧑‍💻 |
| 👻 | list/commit problems | | 🦍 | 🤔 |
| 👻 | download/upload package | | 🦍 | 🤔 |
| 👻 | import polygon problem | 🦍 | 🦍 | 🤔 |
| 👻 | autogen static problem | 🦍 | | 🤔 |

### Icons

- ✅ Done
- 🧑‍💻 In progress
- 🤔 To do
- 👻 Name placeholder
- ⚙️ Refactor task
- 🦍 Engines usage

## Build
```bash
sudo apt install go
make
export PATH=$(pwd)/bin:$PATH
```

## Config

Put your config file in `~/.config/algolymp/config.json`.

```json
{
	"ejudge": {
		"url": "https://ejudge.algocourses.ru",
		"login": "<login>",
		"password": "<password>",
		"judgesDir": "/home/judges"
	},
	"polygon": {
		"url": "https://polygon.codeforces.com",
		"apiKey": "<key>",
		"apiSecret": "<secret>"
	},
	"system": {
		"editor": "nano"
	}
}
```

## blanka
*Create Ejuge contest from template.*

### About

1. Create contest with id from template;
2. Commit changes;
3. *(Optional)* Open contest xml config for editing.

Useful before running [polygon-to-ejudge](https://github.com/grphil/polygon-to-ejudge).

### Flags
- `-i` - new contest id (required)
- `-t` - template contest id (required)
- `-e` - open contest xml config

### Config
- `ejudge.url`
- `ejudge.login`
- `ejudge.password`
- `ejudge.judgesDir` (optional)
- `system.editor` (optional)

### Examples
```bash
blanka --help
blanka -i 51011 -t 51000
blanka -i 51013 -t 51000 -e
```

![blanka logo](https://algolymp.ru/static/img/blanka.png)

## boban
*Filter Ejudge runs.*

### Abount

Filter and print Ejudge runs ids.

### Flags
- `-i` - contest id (required)
- `-f` - filter expression (default: empty)
- `-c` - last runs count (default: 20)

### Config
- `ejudge.url`
- `ejudge.login`
- `ejudge.password`

### Examples
```bash
boban --help
boban -i 47106 -f "prob == 'A'" > runs.txt
boban -i 50014 -f "status == PR" -c 1000
boban -i 50014 -c 10000 2> /dev/null | wc -l
```

![boban logo](https://algolymp.ru/static/img/boban.png)

## casper
*Change Ejudge contest visibility by id.*

### About

- Make contest visible;
- Make contest invisible.

Useful with bash `for` loop at the end of the year.

### Flags
- `-i` - contest id (required)
- `-s` - make contest visible (invisible if not set)

### Config
- `ejudge.url`
- `ejudge.login`
- `ejudge.password`

### Examples
```bash
casper --help
casper -i 41014
casper -i 41014 -s
for i in {41014..41023}; do casper -i $i; done
```

![casper logo](https://algolymp.ru/static/img/casper.png)

## ejik
*Refresh Ejudge contest by id.*

### About

1. Commit changes;
2. Check contest settings;
3. Reload config files.

Useful after running [polygon-to-ejudge](https://github.com/grphil/polygon-to-ejudge).

Feel free to use it after every change

### Flags
- `-i` - contest id (required)
- `-v` - extended output from Ejudge responses

### Config
- `ejudge.url`
- `ejudge.login`
- `ejudge.password`

### Examples
```bash
ejik --help
ejik -v -i 47103
ejik -i 40507
```

![ejik logo](https://algolymp.ru/static/img/ejik.png)

## ripper
*Change Ejudge runs status.*

### About

Change runs status. Designed to work with [boban](#boban) or with raw ids from `stdin`.

**Be careful** using it, double check the [parameters](https://ejudge.ru/wiki/index.php/%D0%92%D0%B5%D1%80%D0%B4%D0%B8%D0%BA%D1%82%D1%8B_%D1%82%D0%B5%D1%81%D1%82%D0%B8%D1%80%D0%BE%D0%B2%D0%B0%D0%BD%D0%B8%D1%8F).

`RJ` is reject, not rejudge. Use `rejudge` status for rejudge.

### Flags
- `-i` - contest id (required)
- `-s` - new status (required, `DQ|IG|OK|PR|RJ|SM|SV|rejudge`)

### Config
- `ejudge.url`
- `ejudge.login`
- `ejudge.password`

### Examples

```bash
ripper --help
ripper -i 51023 -s RJ # read from stdin
cat banlist.txt | ripper -i 47110 -s DQ # ban submits with list
boban -i 52010 -f "prob == 'D' && score >= 50" -c 10000 | ripper -i 52010 -s rejudge # rejudge incorrect group
boban -i 50014 -f "login == 'barmaley' && status == OK" | ripper -i 50014 -s SM # torture a participant
boban -i 48001 -f "status == PR" -c 2000 | ripper -i 48001 -s OK # smart code-review
```

![ripper logo](https://algolymp.ru/static/img/ripper.png)

## scalp
*Set incremental problem scoring using Polygon API.*

### About

1. Enable problem points;
2. Enable problem groups;
3. Load tests metainfo;
4. Store incremental scoring (0 0 5 5 ... 5 6 6 ... 6) with sum of 100.

Very useful for dump problems, prepared in ICPC style.

### Flags
- `-i` - problem id (required)
- `-s` - mark samples as scored tests (put samples in group 0 with 0 score if not set)

### Config
- `polygon.url`
- `polygon.apiKey`
- `polygon.apiSecret`

### Examples
```bash
scalp --help
scalp -i 330352
scalp -i 330328 -s
```

![scalp logo](https://algolymp.ru/static/img/scalp.png)

## valeria
*Build valuer + scorer using Polygon API.*

### About

1. Get problem tests and groups;
2. Build and commit `valuer.cfg` (in Ejudge format);
3. Build and print `scoring.tex` (in Moscow summer olympiad school format).

~~Not very fast now, waiting for `absentInput` parameter in Polygon API.~~

Thanks to Mike, it's been working fast since 30.01.2024.

### Flags
- `-i` - problem id (required)
- `-v` - print valuer.cfg in stderr

### Config
- `polygon.url`
- `polygon.apiKey`
- `polygon.apiSecret`

### Examples
```bash
valeria --help
valeria -i 288808 -v
valeria -i 318511 > scoring.tex
valeria -i 318882 | bat -l tex
```

![valeria logo](https://algolymp.ru/static/img/valeria.png)

## wooda
*Upload problem files filtered by glob to Polygon using API.*

### About

**Now this is a proof of concept. Many more modes will be supported in the future.**

Match all files in directory with glob pattern. Upload recognized files to Polygon.

Supported modes:

- `test`

### Flags
- `-i` - problem id (required)
- `-m` - uploading mode (required)
- `-g` - problem files glob (required)

You should know your shell and probably pass `-g "<glob>"`, not `-g <glob>`

### Config
- `polygon.url`
- `polygon.apiKey`
- `polygon.apiSecret`

### Examples

```bash
wooda --help
wooda -i 337320 -m test -g "tests/*[^.a]"
```

![wooda logo](https://algolymp.ru/static/img/wooda.png)
