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
    launcher = []
    launcher.append(f"export debug={debug}; ")
    launcher.append(f"export background={background}; ")
    launcher.append(f"export freelancer_folder={freelancer_folder}; ")
    launcher.append(f"export timeout={timeout}; ")
    context.obj['launcher'] = launcher
    context.obj['debug'] = debug
    pass


@manage.command()
@click.pass_context
def run(context):
    "launch server"
    context.obj['launcher'].append(f"{PROJECT_MANAGE} runserver")

    if not context.obj['debug']:
        context.obj['launcher'].append(" --noreload --insecure")

    say("".join(context.obj['launcher']))


@manage.command()
@click.pass_context
def shell(context):
    context.obj['launcher'].append(f"{PROJECT_MANAGE} shell")
    say("".join(context.obj['launcher']))


@manage.command()
@click.pass_context
def check(context):
    context.obj['launcher'].append(f"{PROJECT_MANAGE} check --deploy")
    say("".join(context.obj['launcher']))
