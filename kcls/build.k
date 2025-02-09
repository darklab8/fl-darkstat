import .wf_test_build
import .wf_pages_staging
import .wf_pages_production
import infra.kcls.common.workflow

import yaml

_image_name = "darkwind8/darkstat"
_workflows = [
    wf_test_build._workflow
    wf_pages_staging._workflow
    wf_pages_production._workflow
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

_run = [yaml.dump_to_file(a, ".github/workflows/${a._filename}", ignore_private=True, ignore_none=True) for a in _workflows]
