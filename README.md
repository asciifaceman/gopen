# GOPen
It's fairly common for (non-godlike) pen testers to write plenty of scripts in bash to automate their daily tasks, such as nmap and parsing its output. There are tools in place for parsing this output, but in a discussion with a friend I mentioned to him how he could write his own suite of tools for his daily bore-fest tasks. This is just meant to be a quick and dirty example of that.

# Requirements
* Go 1.11.2+
* [go dep](https://github.com/golang/dep)
* nmap installed

# Usage
Usage is decently simple, although some scans still require you to run this program with sudo.

* Dev
  * `go run main.go nmap run [flags]`
* Compiled
  * `gopen nmap run [flags]`

* Flags
  * `--target`, `-t`
    * Can be used multiple times
    * Defines a host target
  * `--pool`, `-p`
    * Number of workers to launch for scanning / parsing
  * `--flags`, `-f`
    * Can be used multiple times
    * Custom nmap flags
    * One at a time, in order
    * Can have up to one whitespace, if wrapped in `"`.
      * ex. `-f "p 1-65535"`

If no flags are defined, it defaults to the flag set `-v -A`.

```
13:41 $ sudo go run main.go nmap run -t scanme.nmap.org -p 3 -f "p 1-65535" -f sV -f sS -f T4
2019-04-25T13:42:12.108-0700    INFO    cmd/run.go:51   Starting gopen...
2019-04-25T13:42:12.108-0700    INFO    nmap/nmap.go:128        has joined the game     {"scanner": 3}
2019-04-25T13:42:12.108-0700    INFO    nmap/nmap.go:138        received task   {"scanner": 3, "target": "scanme.nmap.org"}
2019-04-25T13:42:12.108-0700    INFO    nmap/nmap.go:128        has joined the game     {"scanner": 2}
2019-04-25T13:42:12.108-0700    INFO    nmap/nmap.go:212        has served its time     {"scanner": 2}
2019-04-25T13:42:12.108-0700    INFO    nmap/nmap.go:128        has joined the game     {"scanner": 1}
2019-04-25T13:42:12.108-0700    INFO    nmap/nmap.go:212        has served its time     {"scanner": 1}
2019-04-25T13:43:09.382-0700    INFO    nmap/nmap.go:179        writing results to file {"target": "scanme.nmap.org", "filename": "scanme.nmap.org-12.json"}
2019-04-25T13:43:09.384-0700    INFO    nmap/nmap.go:205        completed task  {"scanner": 3, "target": "scanme.nmap.org", "elapsed": "57.276032427s"}
2019-04-25T13:43:09.384-0700    INFO    nmap/nmap.go:212        has served its time     {"scanner": 3}
```

Outputs `hostname-second.json`

# Development / Compilation
* Dependencies
  * Should be vendored
  * Can run `dep ensure`