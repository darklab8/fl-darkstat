import click
from .universal import say, PROJECT_NAME


@click.group()
def dock():
    "docker commands"
    pass


@dock.command()
def logs():
    say(f"docker logs -t {PROJECT_NAME}:latest")


def builder():
    say("python scripts.py manage -b static")
    say(f"git pull && docker build -t {PROJECT_NAME}:latest .")


@dock.command()
def build():
    builder()


def runner(port):
    say(f"docker run --name {PROJECT_NAME} -t "
        f"-d -p {port}:8000 --rm {PROJECT_NAME}:latest")


@dock.command()
@click.option('--port', '-p', 'port',
              type=int,
              default=8000,
              help="sets docker redirect port")
def run(port):
    runner(port)


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
    say(f"docker rmi $(docker images '{PROJECT_NAME}' -a -q)")


@dock.command()
def clean():
    cleaner()


@dock.command()
@click.option('--port', '-p', 'port',
              type=int,
              default=8000,
              help="sets docker redirect port")
def deploy(port):
    "command to deploy/or redeploy from zero"
    cleaner()
    stopper()
    builder()
    runner(port)


@dock.command()
@click.option("--port", default=80, help="Port number")
def test(port):
    "launches container without daemonization"
    say(f"docker run --name {PROJECT_NAME} -t -p "
        f"{port}:8000 --rm {PROJECT_NAME}:latest")
