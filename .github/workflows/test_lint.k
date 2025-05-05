import infra.kcls.models.github_workflow as gw
import infra.kcls.common.github as ghc
import infra.kcls.actions.checkout_freelancer as cf

import infra.kcls.models.github_action as ga
import infra.kcls.actions.install_go as go
import infra.kcls.actions.install_taskfile as tf
import infra.kcls.actions.install_templ as templ
import infra.kcls.actions.install_swag as swag

import infra.kcls.actions.patch_disco as disco

import .common as c

_workflow = gw.Workflow {
    _filename = "test_lint.yml"
    name = "Test lint"
    on = {
        push.branches = ["master"]
        schedule: [{cron = "0 9 * * *"}]
        pull_request: {}
    }
    jobs = {
        job = _test_job
    }
}
_test_job: gw.Job = {
    name = "Test build"
    steps = [
        ghc.CheckoutRepo
        cf.NewCheckoutFreelancer(cf.CheckoutArgs {
            freelancer_mod = cf._mod_discovery
            freelancer_folder = c._freelancer_folder
        })
        ga.Step {run = "ls ./fl-data"}
        go.InstallGo {}
        tf.InstallTaskfile {}
        templ.InstallTempl {}
        templ.GenerateTempl {}
        swag.SwagInstall {}
        swag.SwagBuild {}
        c.TaskfileBuildVersion
        disco.PatchDisco {
            with = {
                "freelancer-folder" = c._freelancer_folder
            }
        }
        ga.Step {
            name = "lint things"
            uses = "golangci/golangci-lint-action@v8"
        }
    ]
}
