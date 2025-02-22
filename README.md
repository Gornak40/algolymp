# algolymp
*Awesome collection of useful CLI tools for managing Polygon and Ejudge.*

Extended release notes can be found at [chat](https://t.me/algolymp).

## Workflow

| Tool | Description | Ejudge | Polygon | Status |
| --- | --- | :---: | :---: | :---: |
| [baron](#baron) | contest users manager | 🦍 | | ✅ |
| [blanka](#blanka) | create contest | 🦍 | | ✅ |
| [boban](#boban) | filter runs | 🦍 | | ✅ |
| [casper](#casper) | change contest visibility | 🦍 | | ✅ |
| [ejik](#ejik) | commit + check + reload | 🦍 | | ✅ |
| [fara](#fara) | powerful serve.cfg explorer | 🦍 | | ✅ |
| [gibon](#gibon) | api multitool | | 🦍 | ✅ |
| [postyk](#postyk) | print submits 🖨️ | 🦍 | | ✅ |
| [pepel](#pepel) | generate hasher solution | | | ✅ |
| [ripper](#ripper) | change runs status | 🦍 | | ✅ |
| [scalp](#scalp) | incremental scoring | | 🦍 | ✅ |
| [shoga](#shoga) | dump contest tables | 🦍 | | ✅ |
| [valeria](#valeria) | valuer.cfg + tex scoring | | 🦍 | ✅ |
| [vydra](#vydra) | upload package | | 🦍 | 🧪 |
| [wooda](#wooda) | glob problem files upload | | 🦍 | ✅ |
| ⚙️ | move json config to ini | | | 🤔 |
| 👻 | set good random group scores | | 🦍 | 🤔 |
| 👻 | algolymp config manager | | | 🤔 |
| 👻 | import polygon problem | 🦍 | 🦍 | 🤔 |
| 👻 | autogen static problem | 🦍 | | 🤔 |
| 👻 | zip extractor for websites | | | 🤔 |

### Icons

- ✅ Done
- 🧪 Testing
- 🧑‍💻 In progress
- 🤔 To do
- 👻 Name placeholder
- ⚙️ Refactor task
- 🦍 Engines usage

## Build

Download and install the latest version of [Go](https://go.dev/doc/install).

```bash
make
export PATH=$(pwd)/bin:$PATH
```

**Warning:** Some tools may not work on Windows. Feel free to open a PR to fix this.

## Config

Each tool specifies mandatory config variables. It is not necessary to fill in the others.

Put your config file in OS specific config directory:

- `~/.config/algolymp/config.json` on Linux;
- `~/Library/Application Support/algolymp/config.json` on Mac OS;
- `%APPDATA%/algolymp/config.json` on Windows.

Here is an example of a fully filled config:

```json
{
	"ejudge": {
		"url": "https://ejudge.algocourses.ru",
		"login": "<login>",
		"password": "<password>",
		"judgesDir": "/home/judges",
		"secret1": "<random>",
		"uprinter": "Lyceum 9, 3 floor"
	},
	"polygon": {
		"url": "https://polygon.codeforces.com",
		"apiKey": "<key>",
		"apiSecret": "<secret>"
	},
	"system": {
		"editor": "nano",
		"printer": "Samsung SCX-4200"
	}
}
```

**Tip:** You will probably need different configs. It's good practice to name them `config.json.lksh`, `config.json.tbank`, etc. and create a symlink to `config.json`.

## baron
*Ejudge contest users manager.*

### About

Read user ids from `stdin` and flip (inverse) their status in Ejudge contest.

This tool has two major disadvantages:

- It works with ids, not logins;
- It does not support *set* operation, only *flip*.

This is a compromise to not use the [Ejudge API](https://ejudge.ru/swagger/index.html) for as long as possible based on the idea of compatibility.

#### Supported modes

- `invis` - flip invisible
- `ban` - flip ban
- `lock` - flip lock
- `incom` - flip incomplete
- `priv` - flip privilleged

### Flags
- `-i` - contest id (required)
- `-f` - flip mode (required, `invis|ban|lock|incom|priv`)

### Config
- `ejudge.url`
- `ejudge.login`
- `ejudge.password`

### Examples
```bash
baron --help
cat prizewin.csv | baron -i 48501 -f invis
cat cheaters.csv | baron -i 48600 -f ban
baron -i 56001 -f lock # read from stdin
baron -i 46104 -f incom # read from stdin
baron -i 59000 -f priv <<< admin
```

![baron logo](https://algolymp.ru/static/img/baron.png)

## blanka
*Create Ejuge contest from template.*

### About

1. Create contest with id from template;
2. Commit changes;
3. *(Optional)* Open contest xml config for editing.
4. *(Optional)* Create problems directory symlink.

Useful before running [polygon-to-ejudge](https://github.com/grphil/polygon-to-ejudge).

### Flags
- `-i` - new contest id (required)
- `-t` - template contest id (required)
- `-e` - open contest xml config
- `-l` - create problems directory symlink

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

### About

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
*Change Ejudge contest visibility by ids.*

### About

Read contest ids from `stdin` and make them invisible or visible.

Useful with bash `seq`. Logs into Ejuge only once since `v0.14.2`.

### Flags
- `-m` - invisible or visible (required, `hide|show`)

### Config
- `ejudge.url`
- `ejudge.login`
- `ejudge.password`

### Examples
```bash
casper --help
echo 41014 | casper -m hide
seq 40301 40315 | casper -m show
casper -m hide # read from stdin
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

## fara
*Explorer for serve.cfg with mass modify.*

### About

Fara provides a custom selection language for `serve.cfg`.

Queries must follow the following structure:

- `.<field>` for the root section;
- `@<section>:<id>.<field>` for any other section.

The `<field>` parameter is the name of a configuration variable, such as `contest_time` in the global section or `time_limit` in the problem section.

The `<section>` parameter is the name of a section, such as `problem` or `language`.

The `<id>` parameter is the index (starting from 1) of the object in the specified section, e.g. `1` for the first `problem` or `3` for the third `language`.

Parameters `<field>` and `<id>` are optional. You can also pass multiple fields or ids, separating them with commas.

If you do not pass `-d`, `-s` or `-u` flags, fara will output the selected fields. Otherwise it will change them and output the resulting `serve.cfg`.

Some tips for you:
- Use `-q` **or** `-q` and `-d` **or** `-q` and `-u` **or** `-q` and `-s` **or** `-q` and `-u` and `-s`;
- Select sections in `-s` mode, selecting fields may end up with strange result;
- Check the selected fields with `-q` before changing them;
- Do not redirect fara output to the same file as input;
- Check out the examples to learn how best to use this tool.


### Flags
- `-q` - select queries (required)
- `-d` - delete selected fields
- `-u` - update selected fields, delete if `-` passed
- `-s` - field to init/overwrite with `-u` value in selected objects

### Config

No config needed.

### Examples

```bash
fara --help
fara -f /home/judges/048025/conf/serve.cfg -q .score_system,virtual,contest_time
fara -f /home/judges/048025/conf/serve.cfg -q @problem.id,short_name,long_name
fara -f /home/judges/049013/conf/serve.cfg -q @problem.use_stdin,use_stdout -d
fara -f /home/judges/050016/conf/serve.cfg -q @language:2 -d | fara -q @problem:3,4.time_limit -u 15 | bat -l ini
fara -f /home/judges/051009/conf/serve.cfg -q @problem:1,4,6 -s use_ac_not_ok | fara -q @problem:1,4,6 -s ignore_prev_ac > /home/judges/051009/conf/serve.cfg
fara -f serve.cfg -q @problem.id && fara -f serve.cfg -q @problem.id -s max_vm_size -u 512M | fara -q @problem.id -s max_stack_size -u 512M > serve.cfg.new
```

![fara logo](https://algolymp.ru/static/img/fara.png)

## gibon
*Polygon API methods multitool.*

### About

The tool is designed for Polygon API methods outside of the [wooda](#wooda) ideology.

Useful when dealing with large size problems, as API methods do not timeout.

The method `contest` is useful when using [scalp](#scalp) or other [gibon](#gibon) methods.

#### Supported methods

- `contest` - print problem ids in specified contest
- `commit` - commit changes with empty message without email notification
- `download` - download the latest (problem revision) linux package
- `package` - build full package with verification
- `update` - update working copy

### Flags
- `-i` - problem/contest id (required)
- `-m` - method (required)

All methods except `contest` accept a problem id with `-i` flag.

### Config
- `polygon.url`
- `polygon.apiKey`
- `polygon.apiSecret`

### Examples

```bash
gibon --help
gibon -i 42619 -m contest
gibon -i 363802 -m commit
gibon -i 363802 -m download
gibon -i 363802 -m package
gibon -i 363802 -m update
for i in $(gibon -i 42619 -m contest); do gibon -i $i -m commit && gibon -i $i -m package; done
```

![gibon logo](https://algolymp.ru/static/img/gibon.png)

## pepel
*Generate hasher solution based on a/ans/out files.*

### About

Print `python` solution that outputs correct answer for each passed input file and failes on any other.

Useful with Polygon to upload a problem without main correct solution.

**Make sure your input files has `\r\n` line endings (use `unix2dos`), because Polygon works in Windows.**

It's ready to work with any input/output files, encoding and escape sequences don't matter.

Works great with [wooda](#wooda).

Please, add a note to the solution in Polygon (e.g. `Generated by algolymp/pepel`). This will help other problemsetters to avoid misunderstanding.

### Flags
- `-i` - input files glob (required)
- `-a` - answer files glob (required)
- `-z` - zlib compression

You should know your shell and probably use `-i "<glob>"`, not `-i <glob>`.

For large answers, it's strictly recommended to use zlib compression.

### Config

No config needed.

### Examples

```bash
pepel --help
pepel -i "K/*.in.*" -a "K/*.out.*" > pepel.py
pepel -i "tests/*[^.a]" -a "tests/*.a" | bat -l python
pepel -i "*.in.*" -a "*.out.*" -z > pepel-mini.py
```

![pepel logo](https://algolymp.ru/static/img/pepel.png)

## postyk
*Service for printing Ejudge submits.*

### About

Daemon that sends user-requested Ejudge submits to the printer. Useful in team contests in ICPC format.

1. Enable [print_just_copy](https://ejudge.ru/wiki/index.php/Serve.cfg:global:print_just_copy) and [enable_printing](https://ejudge.ru/wiki/index.php/Serve.cfg:global:enable_printing) options in your contest.
2. Share `/var/lib/ejudge/cwork` directory via url `/print/<secret>`. Check out the [nginx](https://nginx.org) configuration example below.
3. Set the `ejudge.secret1` config parameter to `<secret>`.
4. Set the `system.printer` config parameter to hardware printer name. Leave it empty if you don't need printing.
5. Set the `ejudge.uprinter` config parameter to [printer](https://ejudge.ru/wiki/index.php/Userdb:users). Leave it empty if you use only one printer.

You can find out the printer name by using [lsusb](https://manpages.org/lsusb/8) on Linux or [Get-Printer](https://learn.microsoft.com/en-us/powershell/module/printmanagement/get-printer) on Windows.

### Flags
- `-i` - contest id (required)
- `-t` - refresh timeout in seconds (default: 20)

### Config
- `ejudge.url`
- `ejudge.secret1`
- `ejudge.uprinter`
- `system.printer`

### Nginx

```nginx
location /print/<secret> {
	alias /var/lib/ejudge/cwork;
	autoindex on;
}
```

### Examples

```bash
postyk --help
postyk -i 1010
postyk -i 54007 -t 300
```

![postyk logo](https://algolymp.ru/static/img/postyk.png)

## ripper
*Change Ejudge runs status.*

### About

Change runs status. Designed to work with [boban](#boban) or with raw ids from `stdin`.

**Be careful** using it, double check the [parameters](https://ejudge.ru/wiki/index.php/%D0%92%D0%B5%D1%80%D0%B4%D0%B8%D0%BA%D1%82%D1%8B_%D1%82%D0%B5%D1%81%D1%82%D0%B8%D1%80%D0%BE%D0%B2%D0%B0%D0%BD%D0%B8%D1%8F).

`RJ` is reject, not rejudge. Use `rejudge` status for rejudge.

### Flags
- `-i` - contest id (required)
- `-s` - new status (required, `DQ|IG|OK|PR|RJ|SM|SV|rejudge`)
- `-c` - send run comment (if empty, do not send)

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
echo 8 | ripper -i 46512 -s RJ -c "remove #define int long long" # send run comment
```

![ripper logo](https://algolymp.ru/static/img/ripper.png)

## scalp
*Set incremental problem scoring using Polygon API.*

### About

1. Enable problem points;
2. Enable problem groups;
3. Load tests metainfo;
4. Store incremental scoring (0 0 5 5 ... 5 6 6 ... 6) with sum of 100.

Very useful for dumb problems, prepared in ICPC style.

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

## shoga
*Dump Ejudge contest tables.*

### About

Print Ejudges contest tables (CSV format). Various modes are supported.

**Tip:** You can use some custom CSV toolkits, like [xsv](https://github.com/BurntSushi/xsv.git) or [qsv](https://github.com/jqnatividad/qsv.git) to process the output. But I prefer to use vanilla [awk](https://manpages.org/awk) or [cut](https://manpages.org/cut).

#### Supported modes

- `usr` - registered users
- `run` - contest runs
- `stn` - contest standings
- `prb` - contest problems
- `reg` - registration passwords
- `ips` - user IPs

### Flags
- `-i` - contest id (required)
- `-m` - dump mode (required, `usr|run|stn|prb|reg|ips`)

### Config
- `ejudge.url`
- `ejudge.login`
- `ejudge.password`

### Examples
```bash
shoga --help
shoga -i 59000 -m usr # registered users
shoga -i 59000 -m usr | awk '{split($0,a,";"); print a[2]}' # just registered logins
shoga -i 60705 -m usr | cut -d ';' -f 2 | tail -n +2 | sort # just registered logins
shoga -i 55000 -m run # contest runs
shoga -i 436 -m stn # full standings
shoga -i 436 -m stn | cut -d ";" -f 1,2,9,10 | head -n -3 # 6 problems acm contest standings
shoga -i 48005 -m prb # contest problems
shoga -i 51000 -m reg # registration passwords
shoga -i 51000 -m reg | grep myav | cut -d ';' -f 3,6 # specified password
shoga -i 434 -m ips # user IPs
```

![shoga logo](https://algolymp.ru/static/img/shoga.png)

## valeria
*Build valuer + textable using Polygon API.*

### About

1. Get problem tests and groups;
2. Build and commit `valuer.cfg` (in Ejudge format);
3. Build and print `scoring.tex`.

~~Not very fast now, waiting for `absentInput` parameter in Polygon API.~~

Thanks to Mike, it's been working fast since 30.01.2024.

Valeria supports several textable types.

- `universal` - Moscow summer olympiad school format. Works both in PDF and HTML.
- `moscow` - Moscow olympiads format. Works both in PDF and HTML if no variables are passed, otherwise works only in PDF.
- `kaluga` - Kaluga town olympiad format. Works both in PDF and HTML.

### Flags
- `-i` - problem id (required)
- `-v` - print valuer.cfg in stderr
- `-t` - textable type (universal | moscow | kaluga, default: universal)
- `-c` - variables list, useful for some textables (default: nil)

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
valeria -i 285375 -t moscow
valeria -i 285375 -t moscow -c n -c m -c k
valeria -i 396578 -t kaluga
```

![valeria logo](https://algolymp.ru/static/img/valeria.png)

## vydra
*Upload full problem package to Polygon using API.*

### About

**This tool is in beta right now.**

This tool uses `problem.xml` for uploading all package content.

Useful for migration between `polygon.lksh.ru` and `polygon.codeforces.com`.

Designed as an alternative to [polygon-cli](https://github.com/kunyavskiy/polygon-cli).

**Ensure that the problem you are uploading the package into is empty.**

### Known issues

- If problem has testsets other than `tests`, you should create them manually, [issue](https://github.com/Codeforces/polygon-issue-tracking/issues/549);
- If problem is interactive, set `Is problem interactive` checkbox manually;
- If problem has statement resources, upload them manually;
- If problem has custom input/output, set it manually;
- If problem has [FreeMaker](https://freemarker.apache.org) generator, it will expand;
- If problem has stresses, unpload them manually;
- If checker is custom, it's recommended to set `Auto-update` checkbox for `testlib.h`.

### Flags
- `-i` - problem id (required)
- `-p` - problem directory (default: `.`)

### Config
- `polygon.url`
- `polygon.apiKey`
- `polygon.apiSecret`

### Examples

```bash
vydra --help
vydra -i 364022
```

![vydra logo](https://algolymp.ru/static/img/vydra.png)

## wooda
*Upload problem files filtered by glob to Polygon using API.*

### About

Match all files in directory with glob pattern. Upload recognized files to Polygon.

Matching uses natural order (`test.1.in`, `test.2.in`, ..., `test.10.in`, ...).

#### Supported modes

- `test` - test (append, not replace)
- `tags` - tags (each tag is on a new line)
- `val` - validator
- `check` - checker
- `inter` - interactor
- `main` - main solution
- `ok` - correct solution
- `incor` - incorrect solution
- `sample` - sample (append, not replace)
- `image` - statement resource (likely image)

### Flags
- `-i` - problem id (required)
- `-m` - uploading mode (required)
- `-g` - problem files glob (required)

You should know your shell and probably use `-g "<glob>"`, not `-g <glob>`.

### Config
- `polygon.url`
- `polygon.apiKey`
- `polygon.apiSecret`

### Examples

```bash
wooda --help
wooda -i 337320 -m test -g "tests/*[^.a]" # exclude output
wooda -i 337320 -m tags -g tags
wooda -i 337320 -m val -g files/val*.cpp
wooda -i 337320 -m check -g check.cpp
wooda -i 337320 -m inter -g interactor.cpp
wooda -i 337320 -m main -g solutions/main.cpp # Main solution
wooda -i 337320 -m ok -g solutions/sol_apachee.cpp # OK solution
wooda -i 337320 -m incor -g solutions/brute.py # TL solution
wooda -i 337320 -m sample -g "statements/russian/example.[0-9][0-9]"
```

![wooda logo](https://algolymp.ru/static/img/wooda.png)
