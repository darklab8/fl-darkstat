
import click
from scripts.dock import dock
from scripts.test import test
from scripts.celery import celery
from scripts.git import git
from scripts.sphinx import sphinx
from scripts.universal import say, PROJECT_MANAGE


@click.group()
def root():
    "fldarknet commands"
    pass


@root.command()
@click.option('--debug', '-d',
              is_flag=True,
              help="enables debug",
              default=False)
@click.option('--folder', '-f', 'freelancer_folder',
              default="dark_copy",
              help="sets path to freelancer folder for parsing in background, "
              "default='dark_copy'")
@click.option('--timeout', '-t',
              type=int,
              default=1000,
              help="sets timeout between parsing loops")
def run(debug, freelancer_folder, timeout):
    "launch server"
    launcher = []

    launcher.append(f"export debug={debug}; ")
    launcher.append(f"export freelancer_folder={freelancer_folder}; ")
    launcher.append(f"export timeout={timeout}; ")

    launcher.append(f"{PROJECT_MANAGE} runserver")

    if not debug:
        launcher.append(" --noreload --insecure")

    full_command = "".join(launcher)
    say(full_command)


@root.command()
def shell():
    say(f"{PROJECT_MANAGE} shell")


@root.command()
def check():
    say(f"{PROJECT_MANAGE} check --deploy")


root.add_command(dock)
root.add_command(celery)
root.add_command(test)
root.add_command(git)
root.add_command(sphinx)

if __name__ == '__main__':
    root()
