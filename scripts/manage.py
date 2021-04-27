import os
import click
from .universal import say, PROJECT_MANAGE


@click.group()
@click.option('--debug', '-d',
              is_flag=True,
              help="enables debug",
              default=False)
@click.option('--background', '-b',
              is_flag=True,
              help="disables daemon runing in background for parsing",
              default=True)
@click.option('--folder', '-f', 'freelancer_folder',
              default="dark_copy",
              help="sets path to freelancer folder for parsing in background, "
              "default='dark_copy'")
@click.option('--timeout', '-t',
              type=int,
              default=1000,
              help="sets timeout between parsing loops")
@click.pass_context
def manage(context, debug, background, freelancer_folder, timeout):
    "manage commands"
    context.obj['debug'] = debug

    os.environ['debug'] = str(debug)
    os.environ['background'] = str(background)
    os.environ['freelancer_folder'] = str(freelancer_folder)
    os.environ['timeout'] = str(timeout)


@manage.command()
@click.pass_context
def run(context):
    "launch server"
    launcher = []
    launcher.append(f"{PROJECT_MANAGE} runserver")

    if not context.obj['debug']:
        launcher.append(" --noreload --insecure")

    say("".join(launcher))


@manage.command()
def shell():
    say(f"{PROJECT_MANAGE} shell")


@manage.command()
def check():
    say(f"{PROJECT_MANAGE} check --deploy")
