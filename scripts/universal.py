import click
import os

PROJECT_NAME = os.path.basename(os.getcwd())
print(f"project_name = {PROJECT_NAME}")
PROJECT_CORE = "core"
PROJECT_MANAGE = "python manage.py"


def say(phrase):
    click.echo(phrase)
    os.system(phrase)


def bool_to_env(boolean):
    if boolean:
        return 'true'

    return ''
