"some universally used decorators will be here"
from django.core import management
from ship.models import Ship
from contextlib import contextmanager
import sys
import os
from rest_framework.response import Response
from rest_framework import status
import functools
from django.conf import settings
import json


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


def assert_variables(*types_):
    "decorator for fun type checking, but better to use mypy"
    def decorator_repeat(func):
        @functools.wraps(func)
        def wrapper_repeat(*args, **kwargs):

            if settings.DEBUG:
                zip_obj = zip(types_, args)
                for type_, arg in zip_obj:
                    assert isinstance(arg, type_), str(
                        " ".join(
                            (
                                str(arg),
                                'is not',
                                str(type_),
                                'in func',
                                str(repr(func)),
                            )
                        )
                    )

            return func(*args, **kwargs)
        return wrapper_repeat
    return decorator_repeat


def record_to_docs(app, filename, dict_):
    if settings.REFRESH_EXAMPLES:
        with open(os.path.join(
            'sphinx',
            'source',
            f'{app}',
            'write',
            f'{filename}'
        ), 'w') as file_:
            file_.write(json.dumps(dict_, indent=2))
