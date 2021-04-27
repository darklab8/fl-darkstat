import os
import click
from .universal import say


@click.group()
def sphinx():
    "sphinx commands"
    pass


@sphinx.command()
def build():
    "build sphinx documentation"
    say(f"sphinx-build -b html {os.path.join('sphinx','source')} docs")
