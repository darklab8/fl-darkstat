import click
from .universal import say, PROJECT_NAME


@click.group()
def dock():
    "docker commands"
    pass


@dock.command()
def build():
    say(f"git pull && docker build -t {PROJECT_NAME}:latest .")


@dock.command()
def run():
    say(f"docker run --name {PROJECT_NAME} -t "
        f"-d -p 80:8000 --rm {PROJECT_NAME}:latest")


@dock.command()
def stop():
    say('docker stop $(docker ps -a -q --filter="'
        f"name={PROJECT_NAME}"
        '")'
        )


@dock.command()
def clean():
    say('docker rmi $(docker images -a -q)')


@dock.command()
@click.option("--port", default=80, help="Port number")
def test(port):
    "launches container without daemonization"
    say(f"docker run --name {PROJECT_NAME} -t -p "
        f"{port}:8000 --rm {PROJECT_NAME}:latest")
