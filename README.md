![how it looks](docs/howitlooks.png)

# Description

online version of the [flstat](https://discoverygc.com/forums/showthread.php?tid=115254) to navigate game data of the game Freelancer. You can see data about Bases, Guns, Ships and multiple other stuff.

See demos:
- [development version](https://darklab8.github.io/fl-darkstat/)
- [staging version](https://darklab8.github.io/fl-data-discovery/)

# Support

It was made in mind with supporting [Freelancer Discovery](https://discoverygc.com/) as first order. Support will be extended to Vanilla. Any other mode will be supported on request, see contacts.

# Development setup

- git clone https://github.com/darklab8/fl-configs repository for game configs scan, download it to same parent folder as this repository
- go work init ; go work use . ; go work use ../fl-configs
    - initialize Go workspaces, and provide relative path to fl-configs
    - go workspaces allow developing libraries code with real time update of usage to another repository
- install [templ](https://templ.guide/quick-start/installation)
    - go install github.com/a-h/templ/cmd/templ@latest
    - check specific version in [go.mod](./go.mod)
    - In case of emergency we could use vendored in version perhaps
- check [environment variables to set](.vscode/settings.json)
    - set your own environment variable to Freelancer Folder
- install [Taskfile](https://taskfile.dev/usage/) and check [commands to run](Taskfile.yml)
    - run some command, for example task:web

# Features

- Long term maintance support for dozen of years. Minimum dependencies software with Golang and Htmx.
    - for this purpose everything is [go mod vendored in](https://go.dev/ref/mod#go-mod-vendor)
- full GitOps. On commit push to redeploy it automatically
    - See example in [fl-data-discovery repo](https://github.com/darklab8/fl-data-discovery). It contains .github/workflows + game data
- scans Freelancer folder and builds to static assets (html/css/js) deployable to Github pages or any other static assets serving place.
- Usable locally in Linux
  - (Windows support for local usage is coming soon)

# Contacts

discord DM: darkwind8
discord server lab: https://discord.gg/aukHmTK82J
or open [Pull Request here](https://github.com/darklab8/fl-darkstat/issues) and write info there
