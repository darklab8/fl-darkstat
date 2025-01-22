![how it looks](docs/v1.30_all_routes.png)

# Description

online version of the [flstat](https://discoverygc.com/forums/showthread.php?tid=115254) to navigate game data of [the game Freelancer](https://youtu.be/RHlH_qOH5zc). You can see data about Bases, Guns, Ships and multiple other stuff.

See demos:

- [Staging version](https://darklab8.github.io/fl-darkstat/)
- [Freelancer Discovery version](https://darklab8.github.io/fl-data-discovery/) [(action)](https://github.com/darklab8/fl-data-discovery/actions/workflows/publish.yaml)
- [Freelancer Vanilla version](https://darklab8.github.io/fl-data-vanilla/) [(action)](https://github.com/darklab8/fl-data-vanilla/actions/workflows/publish.yaml)
- [Freelancer Sirius Revival](https://darklab8.github.io/fl-data-flsr/) [(action)](https://github.com/darklab8/fl-data-flsr/actions/workflows/publish.yaml)

# Support

- It was made in mind with supporting [Freelancer Discovery](https://discoverygc.com/) as first order.
- Support is extended to Vanilla version.
- Any other mode will be supported on request, see contacts to get in touch.

# Development setup

- git clone https://github.com/darklab8/fl-configs repository for game configs scan, download it to same parent folder as this repository
- install golang of project version or higher (potentially will work anyway).

  - See current golang version [in CI workflow](.github/workflows/deploy.yml)
- install [templ](https://templ.guide/quick-start/installation)

  - go install github.com/a-h/templ/cmd/templ@latest
  - check specific version in [go.mod](./go.mod)
  - In case of emergency we could use vendored in version perhaps
- check [environment variables to set](.vscode/enverant.json)

  - set your own environment variable FREELANCER_FOLDER to Freelancer Folder
  - ensure it was set. `echo $FREELANCER_FOLDER` at Linux or `echo %FREELANCER_FOLDER%` at windows
    - optionally is enough to change value in [enverant.json](.vscode/enverant.json) for that
  - Check to have set other values from [enverant.json](.vscode/enverant.json) ! Some options make development way more pleasant by speeding up rerender by disabling unnecessary features!
- install [Taskfile](https://taskfile.dev/usage/) and check [commands to run](Taskfile.yml)

  - run some command, for example `task web`
- if u wish access to `task dev:watch` that reloads running web server on file changes, then install `pip install watchdog[watchmedo]` and ensure `watchmedo` binary is available to `task dev:watch` command written [in Taskfile](Taskfile.yml)
- If u wish making changes fl-configs and having them right away reflected to fl-darkstat

  - `go work init ; go work use . ; go work use ../fl-configs`
  - initialize Go workspaces, and provide relative path to fl-configs
  - [go workspaces]([https://go.dev/doc/tutorial/workspaces](https://go.dev/doc/tutorial/workspaces)) allow developing libraries code with real time update of usage to another repository
- All dependencies are vendored in with [go mod vendor](https://go.dev/ref/mod#go-mod-vendor) to [vendor folder](https://go.dev/ref/mod#go-mod-vendor) for long term maintanance purposes. We need to to run `go mod vendor` command after library updates for auto refreshing them. Vendored dependencies serve as backup in case of some libs dissapearings.

If u have problems with configuring development environment, then seek my contacts below to help you through it ^_^

# Features

- Long term maintance support for dozen of years. Minimum dependencies software with Golang and Htmx.
  - for this purpose everything is [go mod vendored in](https://go.dev/ref/mod#go-mod-vendor)
- full GitOps. On commit push to redeploy it automatically
  - See example in [fl-data-discovery repo](https://github.com/darklab8/fl-data-discovery). It contains .github/workflows + game data
- scans Freelancer folder and builds to static assets (html/css/js) deployable to Github pages or any other static assets serving place.
- Usable locally for Linux and Windows.
- Only Freelancer Discovery mod and Vanilla are supported at the moment

# Darkstat has API Swagger documentation

- API has swagger documentation accessable from its interface by button "API" at the top of menu

Hint:
- You could generate entire API Client out of openapi with commands like
  - `wget https://repo1.maven.org/maven2/org/openapitools/openapi-generator-cli/7.10.0/openapi-generator-cli-7.10.0.jar -O openapi-generator-cli.jar`
  - `sdk install java 11.0.25-sem` (assuming sdkman is installed)
  - `sdk use java 11.0.25-sem`
  - `java -jar openapi-generator-cli.jar generate -g csharp -i https://darkstat.dd84ai.com/swagger/doc.json -o ./generated_csharp` (example on Discovery deployment)
  - `java -jar openapi-generator-cli.jar generate -g csharp -i ./docs/swagger.json -o ./generated_csharp` (if locally)
  - `java -jar openapi-generator-cli.jar list` to get list of possible API client outputs

# What makes different from regular flstat

- Obviously online
- i also added at last Commodities view with prices per volume ^_^ better reflecting situation for Freelancer Discovery.
- It is interesting to see in Ship details exact Hp Types of equipment you can install onto ship. Other tabs like Guns, Shields, Engines show those Hp Type, so u could find equipment exactly supported for your Light Fighter, Heavy Fighter, Gunboat, Cruser or whatever (u can sort by column to find all such equipment)
- Tractors tab has info regarding Discovery IDs and where to buy them ^_^
- other extra tabs like Engines, CMs added
- Tabs for different equipment could be showing more full list of equipment in "Show all" mode.
- Has searching/filtering options with multiple matching items shown
- You can pin items for comparison
- For Discovery Freelancer, u can select ID/Tractor and having guns/ships etc filtered/shown according to what your ID can use without power core regeneration debuffs. Shows ID compatibility (75% ID compatibility at any equiped item will mean having only 75% of Power core regeneration)

# Usage

## Local usage

- download latest from https://github.com/darklab8/fl-darkstat/releases , they are autobuilt from CI, so they will be always there.
  - optionally build latest [according to instruction](https://github.com/darklab8/fl-darkstat/blob/550b40a49ec4f5dd1113457e4c96eee161296b7b/.github/actions/build/action.yml#L25) on your own if desired
- put file into root of Freelancer folder and start
  - optionally launch from anywhere, just add env variable FREELANCER_FOLDER with location to freelancer folder root.
- visit http://localhost:8000/ as printed in console to see web site locally
- Launching from `cmd` or any other console at Freelancer Discovery folder path is preferable. Because u will see detailed log output.

P.S. The tool uses lazy filesystem approach by grabbing first file with matching name. I did not use full paths.
So don't have folder "DATA2" duplicating all files in same FreelancerDiscovery folder

## Docker usage

- [Docker releases](https://hub.docker.com/r/darkwind8/darkstat) are available too. tag `production` is latest stable and running in prod.
- Configuration for its running check in terraform infra code of [module darkstat](./tf/modules/darkstat)
  - you need to point at least volume -v /data:/path_to_frelancer_folder
  - and point required environment variables [as described there](./tf/modules/darkstat/variables.tf)
  - docker images are built for amd64 and arm64 :)

# Acknowledments

- The tool was strongly inspired by [flstat](https://discoverygc.com/forums/showthread.php?tid=115254) originally written by Dan Tascau
  - regretfully original code was not found
  - some things helped from [patch written in Assembly to flstat](http://adoxa.altervista.org/freelancer/tools.html) by Adoxa
- In general a lot of stuff was checked from [Starport wiki](https://the-starport.com/wiki/)
  - [Bribing probabilities](https://the-starport.com/forums/topic/5372/bribe-probabilities/6?topic_id=5565) were inspired by Adoxa conversation at Starport in 2014
  - Also stuff like [market stuff](https://the-starport.com/wiki/ini-editing/typed-inis/markets/) page helped too
- Formulas for angular stuff were found in [flint](https://github.com/biqqles/flint/blob/master/flint/entities/ship.py#L82)
- [Discord Community in starport](https://discord.gg/freelancer-galactic-community-638984923591737355) also answered multiple questions
  - as well as Freelancer Discovery dev community
- Also thanks to The Alex (From Freelancer Discovery) for getting me [Python script for reading dlls](https://github.com/darklab8/fl-configs/blob/master/docs/inspiration/dll_reading/alex_py/main.py)
  - That helped rewriting it in go for [fl-configs lib](https://github.com/darklab8/fl-configs)
- Honorary mentions for very active moral support and extra ideas by
  - IrateRedKite (from starport Discord)
  - Bolte (from starport Discord)

<!--- 
- In case it will be ever needed, [just in case linking flcompanion](<https://github.com/Corran-Raisu/FLCompanion>)
- check Selfpatch for fl-data-discovery later https://github.com/Lazrius/DSLauncher/tree/default/Self%20Patch
-->

# Contacts

- discord DM: darkwind8
- discord server lab: https://discord.gg/zFzSs82y3W
- or open [Pull Request here](https://github.com/darklab8/fl-darkstat/issues) and write there
- [see up to date contacts here](https://darklab8.github.io/blog/index.html#contacts)

See anouncements at [Discovery Freelancer forum thread](https://discoverygc.com/forums/showthread.php?tid=187294)

# License

fl-darkstat was originally created by Andrei Novoselov (aka darkwind, aka dd84ai)
The work is released under AGPL license, free to modify, copy and etc. as long as you keep code open source and mentioned original author.
See [LICENSE](./LICENSE) file for details.

