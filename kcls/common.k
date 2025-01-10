import infra.kcls.models.github_action as ga
import infra.kcls.actions.install_go as go
import infra.kcls.actions.install_taskfile as tf
import infra.kcls.actions.install_templ as templ
import infra.kcls.actions.patch_disco as disco

_freelancer_folder = r"${{ github.workspace }}/fl-data"
_image_name = "darkwind8/darkstat"

schema BuildArgs:
    site_root: str
    freelancer_folder: str
    heading: str
    relay_host: str = ""
    is_detailed: bool = False

TaskfileBuildVersion = ga.Step {
    name = "add versions"
    run = "task build-version"
}

DarkstatBuildAssetsSteps = lambda a: BuildArgs {
    # return
    [
        ga.Step {run = "ls ./fl-data"}
        go.InstallGo {}
        tf.InstallTaskfile {}
        templ.InstallTempl {}
        templ.GenerateTempl {}
        TaskfileBuildVersion
        disco.PatchDisco {
            with = {
                "freelancer-folder" = a.freelancer_folder
            }
        }
        ga.Step {
            name = "Test things"
            run = "task test -- -v"
            env = {FLDARKDATA_LOG_LEVEL = "DEBUG"}
        }
        ga.Step {
            name = "build"
            run = "task build"
            env = {
                SITE_ROOT = a.site_root
                FREELANCER_FOLDER = a.freelancer_folder
                FLDARKSTAT_HEADING = a.heading
                if a.is_detailed:
                    DARKSTAT_DETAILED = "true"
                else:
                    DARKSTAT_DETAILED = "false"
                if a.relay_host != "":
                    RELAY_HOST = a.relay_host
            }
        }
    ]
}
