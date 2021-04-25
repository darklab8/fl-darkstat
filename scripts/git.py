import click
from .universal import say


@click.group()
def git():
    "git commands"
    pass


# cool hook against commits with errors
hook_install = 'flake8 --install-hook git'
hook_on = 'git config --bool flake8.strict true'
hook_off = 'git config --bool flake8.strict false'


@git.command()
def hook():
    say('git config user.email "dd84ai@gmail.com"')
    say('git config user.name "dd84ai"')
    say(hook_on)


@git.command()
def unhook():
    say(hook_off)
