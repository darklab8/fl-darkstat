> [!NOTE]  
> **The project is still in early development stage and haven't reached MVP (Minimal Viable Product)**

# Description

Remake in golang of online navmap for [Freelancer Discovery](<https://discoverygc.com/>). Newer version should be easier to maintain, and should be way more rapid due to design being a static site generator powered by htmx. Potentially also implementing different new features :]

# Support

- It will be made in mind with supporting [Freelancer Discovery](https://discoverygc.com/) as first order.
- Support will be extended to Vanilla version.
- Any other mode will be supported on request, see contacts to get in touch.

# Development setup

- git clone https://github.com/darklab8/fl-configs repository for game configs scan, download it to same parent folder as this repository
- install golang of project version or higher (potentially will work anyway).
  - See current golang version [in CI workflow](.github/workflows/deploy.yml)
- install [templ](https://templ.guide/quick-start/installation)
  - go install github.com/a-h/templ/cmd/templ@latest
  - check specific version in [go.mod](./go.mod)
- check [environment variables to set](.vscode/enverant.json)
  - set your own environment variable FREELANCER_FOLDER to Freelancer Folder
  - ensure it was set. `echo $FREELANCER_FOLDER` at Linux or `echo %FREELANCER_FOLDER%` at windows
    - optionally is enough to change value in [enverant.json](.vscode/enverant.json) for that
  - Check to have set other values from [enverant.json](.vscode/enverant.json) ! Some options make development way more pleasant by speeding up rerender by disabling unnecessary features!
- install [Taskfile](https://taskfile.dev/usage/) and check [commands to run](Taskfile.yml)
  - run some command, for example `task web`
- if u wish access to `task dev:watch` that reloads running web server on file changes, then install `pip install watchdog[watchmedo]` and ensure `watchmedo` binary is available to `task dev:watch` command written [in Taskfile](Taskfile.yml)
- If u wish making changes fl-configs and having them right away reflected to fl-darkstat (same for fl-darkcore)
  - `go work init ; go work use . ; go work use ../fl-configs`
  - initialize Go workspaces, and provide relative path to fl-configs
  - [go workspaces]([https://go.dev/doc/tutorial/workspaces](https://go.dev/doc/tutorial/workspaces)) allow developing libraries code with real time update of usage to another repository

If u have problems with configuring development environment, then seek my contacts below to help you through it ^_^

# Usage locally

To be written (copy paste from darkstat ;)

# Features

To be writtem...

# Acknowledments

- Inspired by old navmap made by @AudunVN
    - https://discoverygc.com/forums/showthread.php?tid=132266
    - https://github.com/AudunVN/Navmap

# Contacts

- discord DM: darkwind8
- discord server lab: https://discord.gg/aukHmTK82J
- or open [Pull Request here](https://github.com/darklab8/fl-darmap/issues) and write there
