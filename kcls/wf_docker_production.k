import infra.kcls.models.github_action as ga
import infra.kcls.models.github_workflow as gw
import infra.kcls.common.github as ghc
import infra.kcls.actions.install_go as go
import infra.kcls.actions.install_taskfile as tf
import infra.kcls.actions.install_autogit as autogit

_workflow_name = "Deploy docker Production"
_workflow = gw.Workflow {
    _filename = "docker-production.yml"
    name = _workflow_name
    on = {
        workflow_dispatch: {}
        push.tags: ["*"]
    }
    env = ghc.DockerHostEnv
    jobs = {
        build = gw.JobDispatch {
            name = "Build docker image and save to latest"
            uses = "darklab8/fl-darkstat/.github/workflows/deploy-docker-staging.yml@master"
        }
        job = {
            name = _workflow_name
            needs = ["build"]
            steps = [
                ghc.CheckoutRepoWithCommits
                go.InstallGo {}
                tf.InstallTaskfile {}
                autogit.InstallAutogit {}
                ghc.GetVersionFromTag
                ghc.DockerLogin
                ghc.InstallDarklabSshKey
                ga.Step {
                    name = "Push image"
                    run = """set -ex
docker pull darkwind8/darkstat:latest
sleep 5
docker tag darkwind8/darkstat:latest darkwind8/darkstat:${ghc.GetVersionOutput}
docker tag darkwind8/darkstat:latest darkwind8/darkstat:production
docker push darkwind8/darkstat:${ghc.GetVersionOutput}
sleep 10
docker push darkwind8/darkstat:production
"""
                }
                ga.Step {
                    name = "Switch service image to new one"
                    run = """set -ex
docker pull darkwind8/darkstat:production
sleep 5
docker service update --image darkwind8/darkstat:production darkstat-production
"""
                }
            ]
        }
    }
}
