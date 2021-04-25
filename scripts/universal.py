import click
import os

PROJECT_NAME = "fldarknet"
PROJECT_CORE = "core"
PROJECT_MANAGE = "python manage.py"


def say(phrase):
    click.echo(phrase)
    os.system(phrase)
