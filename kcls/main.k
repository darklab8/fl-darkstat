import .workflow_test as workflow_test
import yaml

_workflows = [
    workflow_test._workflow_test
]

_run = [yaml.dump_to_file(a, ".github/workflows/${a._filename}", ignore_private=True, ignore_none=True) for a in _workflows]
