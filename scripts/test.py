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
    say("".join(("flake8 --exclude .git,venv,*/migrations/*,.tox",
                " --ignore=A002,W503,W504 .")))


@test.command()
@click.option('--refresh', '-r',
              is_flag=True,
              help="enables refresh of data examples",
              default=False)
@click.option('--cover', '-c', 'cover',
              is_flag=True,
              help="shows coverage",
              default=False)
@click.option('--pipline', '-p', 'pipline',
              is_flag=True,
              help="enables long tests",
              default=False)
def unit(refresh, cover, pipline):
    "get unit tests"
    os.environ['refresh'] = bool_to_env(refresh)
    os.environ['pipline'] = bool_to_env(pipline)
    launcher = ["pytest -n 6"]
    if cover:
        launcher.append("--cov=.")
    say(" ".join(launcher))


@test.command()
def tox():
    "full test run to be done between commits"
    say("tox -r")


@test.command()
def mypy():
    "type hinting checker"
    say("mypy .")

