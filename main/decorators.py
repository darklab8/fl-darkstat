"some universally used decorators will be here"
from django.core import management
from ship.models import Ship
from contextlib import contextmanager
import sys
import os


@contextmanager
def suppress_stdout(supress_errors=False):
    with open(os.devnull, "w") as devnull:
        old_stdout = sys.stdout
        sys.stdout = devnull

        if supress_errors:
            old_stderr = sys.stderr
            sys.stderr = devnull
        try:
            yield
        finally:
            sys.stdout = old_stdout

            if supress_errors:
                sys.stderr = old_stderr


def loaded_db(func):
    "useful decorator to preload stuff for unit tests"
    def decorator_function(*args, **kwargs):
        # print('executing '+func.__name__)

        with suppress_stdout():
            if len(Ship.objects.all()) == 0:
                management.call_command('loaddata', 'dump.json')

        return func(*args, **kwargs)
    return decorator_function
