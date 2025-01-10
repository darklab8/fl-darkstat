import infra.kcls.models.github_action as ga
import infra.kcls.models.github_workflow as gw
import infra.kcls.common.github as ghc
import infra.kcls.actions.install_go as go
import infra.kcls.actions.install_taskfile as tf
import infra.kcls.actions.install_autogit as autogit

_workflow_name = "Deploy docker staging"
_workflow = gw.Workflow {
    _filename = "docker-staging.yml"
    name = _workflow_name
    on = {
        push.branches = ["master"]
        workflow_dispatch: {}
        workflow_call: {}
    }
    env = ghc.DockerHostEnv
    jobs = {
        job = {
            name = _workflow_name
            steps = [
                ghc.CheckoutRepoWithCommits
                go.InstallGo {}
                tf.InstallTaskfile {}
                autogit.InstallAutogit {}
                autogit.GetVersionStep
                ghc.DockerLogin
                ghc.InstallDarklabSshKey
                ga.Step {
                    name = "Build image"
                    run = """docker build --build-arg "BUILD_VERSION=${autogit.GetVersionOutput}" --tag darkwind8/darkstat:${autogit.GetVersionOutput} ."""
                }
                ga.Step {
                    name = "Push image"
                    run = """set -ex
docker tag darkwind8/darkstat:${autogit.GetVersionOutput} darkwind8/darkstat:staging
docker tag darkwind8/darkstat:${autogit.GetVersionOutput} darkwind8/darkstat:latest
docker push darkwind8/darkstat:staging
sleep 5
docker push darkwind8/darkstat:latest
"""
                }
                ga.Step {
                    name = "Switch service image to new one"
                    run = """set -ex
docker pull darkwind8/darkstat:staging
sleep 5
docker service update --image darkwind8/darkstat:staging darkstat-staging
"""
                }
            ]
        }
    }
}
