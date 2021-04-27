"some universally used decorators will be here"
from django.core import management
from ship.models import Ship
from contextlib import contextmanager
import sys
import os
from rest_framework.response import Response
from rest_framework import status
import functools


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


def required_key(key):
    "decorator to ask for API_key in some request"
    def decorator_repeat(func):
        @functools.wraps(func)
        def wrapper_repeat(*args, **kwargs):

            if key != args[1].GET.get('api'):
                return Response(data={'error': 'wrong api key'},
                                status=status.HTTP_400_BAD_REQUEST)
            return func(*args, **kwargs)
        return wrapper_repeat
    return decorator_repeat
