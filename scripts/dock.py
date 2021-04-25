import click
from .universal import say, PROJECT_NAME


@click.group()
def dock():
    "docker commands"
    pass


def builder():
    say(f"git pull && docker build -t {PROJECT_NAME}:latest .")


@dock.command()
def build():
    builder()


def runner():
    say(f"docker run --name {PROJECT_NAME} -t "
        f"-d -p 80:8000 --rm {PROJECT_NAME}:latest")


@dock.command()
def run():
    runner()


def stopper():
    say('docker stop $(docker ps -a -q --filter="'
        f"name={PROJECT_NAME}"
        '")'
        )


@dock.command()
def stop():
    stopper()


def cleaner():
    "getting rid of already built docker layers"
    say('docker rmi $(docker images -a -q)')


@dock.command()
def clean():
    cleaner()


@dock.command()
def deploy():
    "command to deploy/or redeploy from zero"
    cleaner()
    stopper()
    builder()
    runner()


@dock.command()
@click.option("--port", default=80, help="Port number")
def test(port):
    "launches container without daemonization"
    say(f"docker run --name {PROJECT_NAME} -t -p "
        f"{port}:8000 --rm {PROJECT_NAME}:latest")
