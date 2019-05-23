[![Go Report Card](https://goreportcard.com/badge/github.com/williamleven/BooGroCha)](https://goreportcard.com/report/github.com/williamleven/BooGroCha)
[![Build Status](https://travis-ci.com/williamleven/BooGroCha.svg?branch=master)](https://travis-ci.com/williamleven/BooGroCha)
# BooGroCha
A CLI client for chalmers' group room booking system.

## Features
- [x] Book group rooms
- [x] Integration Chalmers booking system
- [x] Save personal settings in config file
- [x] List available bookings for selected time
- [x] Delete bookings
- [x] List bookings
- [x] Relative times for bookings
- [ ] [AUR package](https://github.com/williamleven/BooGroCha/issues/4)
- [ ] [Rooms sorted base on user preferences](https://github.com/williamleven/BooGroCha/issues/5)
- [ ] [Option to only show rooms from preferred campus](https://github.com/williamleven/BooGroCha/issues/6)
- [ ] [Integrate with Chalmers Library booking system](https://github.com/williamleven/BooGroCha/issues/7)
- [ ] [Integrate with Johanneberg Student Union building booking system](https://github.com/williamleven/BooGroCha/issues/8)
- [ ] [Prompt for password if not set](https://github.com/williamleven/BooGroCha/issues/9)


## Installation

### Compile from source

#### Requirements
* go (>= 1.12)
* $GOPATH/bin in $PATH

#### Steps
```bash
$ git clone https://github.com/williamleven/BooGroCha
$ cd BooGroCha/cmd/bgc
$ go install
```

### AUR package
TBA

## Usage

### Book a room
Allows for booking rooms by showing the available rooms for the given date and time.

```bash
$ bgc book <date> <time>
```
* **\<date\>** can be either an absolute date (`YYYYMMDD`, `YYMMDD`) or a relative date (`MMDD`, `DD`, `D`, `today`, `tomorrow`, `monday`, `tuesday` ...)
* **\<time\>** can be either (`HH:mm-HH:mm`, `HH-HH`, `H-H`, `H-HH`, `HH-H`) or aliases like `lunch`

### List booked rooms

```bash
$ bgc list
```

### Delete booked rooms

```bash
$ bgc delete
```

### Configuration
Allows the user to set parameters in a config file.

#### Setting a variable
```bash
$ bgc config set <variable> <value>
```
* **\<variable\>** can for example be `cid` or `pass`
* **\<value\>** should be the value that you want to set the variable to (NOTE: when setting the password you will be prompted for input instead of setting it directly)


#### Showing a variable
```bash
$ bgc config get <variable>
```
* **\<variable\>** can for example be `cid` or `campus` (NOTE: You cannot show the `pass` variable this way)

## Development

```bash
$ git clone https://github.com/williamleven/BooGroCha
$ cd BooGroCha/cmd/bgc
$ go build
$ ./bgc ...
```

### GoLand
At the moment of writing (version 2018.3.5) GoLand doesn't support go modules in a project by default.
Instead you have to tick the `Enable Go Modules (vgo) integration` under `Settings -> Go -> Go Modules (vgo)` and set the proxy setting to `direct`.

### Contribution
This is a open project so feel free to contribute with code and/or ideas in issues and pull requests.

#### Structure
In this project we have decided on [a specific project structure](https://github.com/golang-standards/project-layout) and would like you to follow this structure as well when contributing code to the project.
