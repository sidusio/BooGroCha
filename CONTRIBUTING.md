# Contributing
This is an open project so feel free to contribute with code and/or ideas in issues and pull requests.

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

### Structure
In this project we have decided on [a specific project structure](https://github.com/golang-standards/project-layout) and would like you to follow this structure as well when contributing code to the project.
