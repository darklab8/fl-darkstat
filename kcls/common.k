import infra.kcls.models.github_action as ga
import infra.kcls.actions.install_go as go
import infra.kcls.actions.install_taskfile as tf
import infra.kcls.actions.install_templ as templ
import infra.kcls.actions.install_swag as swag

import infra.kcls.actions.patch_disco as disco
import infra.kcls.actions.darkstat

_freelancer_folder = r"${{ github.workspace }}/fl-data"

TaskfileBuildVersion = ga.Step {
    name = "add versions"
    run = "task build-version"
}

DarkstatBuildAssetsSteps = lambda a: darkstat.BuildArgs {
    # return
    [
        ga.Step {run = "ls ./fl-data"}
        go.InstallGo {}
        tf.InstallTaskfile {}
        templ.InstallTempl {}
        templ.GenerateTempl {}
        swag.SwagInstall {}
        swag.SwagBuild {}
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
                SITE_HOST = a.site_host
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
