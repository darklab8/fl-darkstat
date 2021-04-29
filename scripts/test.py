import os
import click
from .universal import say, PROJECT_CORE, bool_to_env


@click.group()
def test():
    "testing commands"
    pass


@test.command()
def pylint():
    "link with pylint"
    say("pylint --load-plugins pylint_django --django-settings-module="
        '"' + f"{PROJECT_CORE}" + '.settings"'
        " --disable=django-not-configured --exit-zero `ls -d */`"
        )


@test.command()
def flake():
    "lint with flake8"
    say("flake8 --exclude .git,venv,*/migrations/*,.tox .")


@test.command()
@click.option('--refresh', '-r',
              is_flag=True,
              help="enables refresh of data examples",
              default=False)
@click.option('--app', '-a', 'app',
              default="",
              help="choose to test particular app")
def unit(refresh, app):
    "get unit tests"
    os.environ['refresh'] = bool_to_env(refresh)
    say(
        "coverage run --omit 'venv/*,.tox/*'"
        f" --source='.' manage.py test {app}")


@test.command()
def matrix():
    "full test run to be done between commits"
    say("tox -r")


@test.command()
def cover():
    """get coverage of unit tests
    it should be used only after 'unit' command"""
    say("coverage report")
