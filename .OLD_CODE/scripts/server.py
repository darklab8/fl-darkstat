import click
from .universal import say


@click.group()
def server():
    "server commands"
    pass


@server.command()
def connect():
    say("ssh 80.78.247.245")
