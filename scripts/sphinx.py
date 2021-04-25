import click
from .universal import say


@click.group()
def sphinx():
    "sphinx commands"
    pass


@sphinx.command()
def build():
    "link with pylint"
    say("sphinx-build -b html docs/source docs/build")
