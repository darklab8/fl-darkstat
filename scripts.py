
import click
from scripts.dock import dock
from scripts.test import test
from scripts.celery import celery
from scripts.git import git
from scripts.sphinx import sphinx
from scripts.django import django


@click.group()
@click.pass_context
def root(context):
    "root commands"
    pass


root.add_command(dock)
root.add_command(celery)
root.add_command(test)
root.add_command(git)
root.add_command(sphinx)
root.add_command(django)

if __name__ == '__main__':
    root(obj={})
