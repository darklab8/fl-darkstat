import click
from .universal import say


@click.group()
def celery():
    "celery commands"
    pass


@celery.command()
def start():
    say("celery -A core beat -l INFO")


@celery.command()
def show():
    say("celery -A core worker -l INFO")
