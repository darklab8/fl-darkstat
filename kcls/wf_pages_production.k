import infra.kcls.models.github_workflow as gw
import infra.kcls.common.github as ghc
import infra.kcls.models.github_action as ga
import infra.kcls.actions.install_go as go
import infra.kcls.actions.install_taskfile as tf
import infra.kcls.actions.install_templ as templ
import infra.kcls.actions.install_autogit as autogit

import .common as c

_workflow = gw.Workflow {
    _filename = "pages-production.yml"
    name = "Deploy pages production"
    on = {
        workflow_dispatch: {}
        push.tags = ["*"]
    }
    jobs = {
        job = _pages_production_job
    }
}
_pages_production_job: gw.Job = {
    name = "Github Pages Production"
    steps = [
        autogit.InstallAutogit {}
        ghc.CheckoutRepoWithCommits
        go.InstallGo {}
        tf.InstallTaskfile {}
        templ.InstallTempl {}
        templ.GenerateTempl {}
        c.TaskfileBuildVersion
        ga.Step {
            name: "Build Binaries"
            run: "task build:bins"
        }
        autogit.GetChangelogStep {}
        ghc.UploadBinariesStep(ghc.UploadBinariesStepOpts {
            body = autogit.GetChangelogOutput
            files = [
                "dist/darkstat-linux-amd64"
                "dist/darkstat-linux-arm64"
                "dist/darkstat-macos-amd64"
                "dist/darkstat-macos-arm64"
                "dist/darkstat-windows-amd64.exe"
                "dist/darkstat-windows-arm64.exe"
            ]
        })
        ghc.TriggerWorkflowDispatch(ghc.TriggerWorkflowDispatchOpts {
            repository = "darklab8/fl-data-vanilla"
        }) {name = "Repository Dispatch for Vanilla"}
        ghc.TriggerWorkflowDispatch(ghc.TriggerWorkflowDispatchOpts {
            repository = "darklab8/fl-data-discovery"
        }) {name = "Repository Dispatch for Discovery"}
        ghc.TriggerWorkflowDispatch(ghc.TriggerWorkflowDispatchOpts {
            repository = "darklab8/fl-data-flsr"
        }) {name = "Repository Dispatch for FLSR"}
    ]
}
