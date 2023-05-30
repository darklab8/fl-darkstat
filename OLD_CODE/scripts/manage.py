import os
import click
from .universal import say, PROJECT_MANAGE, bool_to_env


@click.group()
@click.option('--debug',
              '-d',
              is_flag=True,
              help="enables debug",
              default=False)
@click.option('--background',
              '-b',
              'enable_background',
              is_flag=True,
              help="disables daemon runing in background for parsing",
              default=False)
@click.option('--folder',
              '-f',
              'freelancer_folder',
              default="dark_copy",
              help="sets path to freelancer folder for parsing in background, "
              "default='dark_copy'")
@click.option('--timeout',
              '-t',
              type=int,
              default=1000,
              help="sets timeout between parsing loops")
@click.pass_context
def manage(context, debug, enable_background, freelancer_folder, timeout):
    "manage commands"
    context.obj['debug'] = debug
    os.environ['debug'] = bool_to_env(debug)
    os.environ['enable_background'] = bool_to_env(enable_background)
    os.environ['freelancer_folder'] = str(freelancer_folder)
    os.environ['timeout'] = str(timeout)


@manage.command()
@click.option('--ip-port',
              '-p',
              'address',
              default="0.0.0.0:8000",
              help="set ip address and port")
@click.pass_context
def run(context, address, venv):
    "launch server"
    if context.obj['debug']:
        say(f"{PROJECT_MANAGE} runserver {address}")
    else:
        say("venv/bin/python scripts.py manage -b static")
        say(f"venv/bin/gunicorn core.wsgi -b {address}")


def staticer():
    say("mkdir static")
    say(f"{PROJECT_MANAGE} collectstatic -c --noinput")


@manage.command()
def static():
    staticer()


@manage.command()
def shell():
    say(f"{PROJECT_MANAGE} shell")


@manage.command()
def check():
    say(f"{PROJECT_MANAGE} check --deploy")
