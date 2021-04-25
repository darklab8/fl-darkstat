
import click
from scripts.dock import dock
from scripts.test import test
from scripts.celery import celery
from scripts.universal import say, PROJECT_MANAGE


@click.group()
def root():
    "fldarknet commands"
    pass


@root.command()
@click.option('--debug/--no-debug', default=True)
def run(debug):
    "launch server"
    if debug:
        say(f"export DEBUG=true; {PROJECT_MANAGE} runserver")
    else:
        say(f"{PROJECT_MANAGE} runserver --noreload --insecure")


@root.command()
def shell():
    say(f"{PROJECT_MANAGE} shell")


@root.command()
def check():
    say(f"{PROJECT_MANAGE} check --deploy")


root.add_command(dock)
root.add_command(celery)
root.add_command(test)

if __name__ == '__main__':
    root()
