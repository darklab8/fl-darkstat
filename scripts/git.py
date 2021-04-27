import click
from .universal import say


@click.group()
def git():
    "git commands"
    pass


@git.command()
def creds():
    say('git config user.email "dd84ai@gmail.com"')
    say('git config user.name "dd84ai"')
