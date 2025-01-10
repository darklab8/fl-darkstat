import .wf_test_build
import .wf_pages_staging
import .wf_pages_production
import .wf_docker_staging
import .wf_docker_production

import yaml

_workflows = [
    wf_test_build._workflow
    wf_pages_staging._workflow
    wf_pages_production._workflow
    wf_docker_staging._workflow
    wf_docker_production._workflow
]

_run = [yaml.dump_to_file(a, ".github/workflows/${a._filename}", ignore_private=True, ignore_none=True) for a in _workflows]
