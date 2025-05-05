import yaml

import .workflows.test_build
import .workflows.pages_production
import .workflows.pages_staging
import infra.kcls.common.workflow

import .workflows.test_lint

_image_name = "darkwind8/darkstat"
_workflows = [
    test_lint._workflow
    test_build._workflow
    pages_staging._workflow
    pages_production._workflow
    workflow.DockerDeploySwarm(workflow.DockerDeploySwarmOpts {
        image_name = _image_name
        service_name = "darkstat"
        environment = "staging"
    })
    workflow.DockerDeploySwarm(workflow.DockerDeploySwarmOpts {
        image_name = _image_name
        service_name = "darkstat"
        environment = "production"
        extra_services = ["dev-darkstat-app", "vanilla-darkstat-app"]
    })
]

_run = [yaml.dump_to_file(a, "workflows/${a._filename}", ignore_private=True, ignore_none=True) for a in _workflows]
