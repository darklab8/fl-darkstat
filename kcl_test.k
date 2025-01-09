import infra.kcls.models.github_action as ga
import infra.kcls.models.github_workflow as gw
import infra.kcls.common.github as ghc
import infra.kcls.actions.checkout_freelancer as cf
import infra.kcls.actions.install_go as go
import infra.kcls.actions.install_taskfile as tf
import infra.kcls.actions.install_templ as templ
import infra.kcls.actions.patch_disco as disco

import yaml

_freelancer_folder = r"${{ github.workspace }}/fl-data"

yaml.dump_to_file(gw.Workflow {
    name = "Test build"
    on = {
        push.branches = ["master"]
        schedule: [{cron = "0 9 * * *"}]
    }
    jobs = {
        test = {
            name = "Test build"
            steps = [
                ghc.CheckoutRepo
                cf.CheckoutFreelancer {
                    with = {
                        "freelancer-mod": "discovery"
                        "freelancer-folder": _freelancer_folder
                        "ssh-key-base64-discovery": r"${{ secrets.ID_RSA_FILES_FREELANCER_DISCOVERY }}"
                    }
                }
                ga.Step {run = "ls ./fl-data"}
                go.InstallGo {}
                tf.InstallTaskfile {}
                templ.InstallTempl {}
                templ.GenerateTempl {}
                ga.Step {
                    name = "add versions"
                    run = "task build-version"
                }
                disco.PatchDisco {
                    with = {
                        "freelancer-folder" = _freelancer_folder
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
                        SITE_ROOT = "/fl-darkstat/"
                        FREELANCER_FOLDER = _freelancer_folder
                        FLDARKSTAT_HEADING = '<span style="font-weight:1000;">DEV ENV</span> <a href="https://github.com/darklab8/fl-darkstat">fl-darkstat</a> for <a href="https://github.com/darklab8/fl-data-discovery">Freelancer Discovery</a>'
                        DARKSTAT_DETAILED = "true"
                        RELAY_HOST = "https://darkrelay.dd84ai.com"
                    }
                }
            ]
        }
    }
}, ".github/workflows/test.yml", ignore_private=True, ignore_none=True)
