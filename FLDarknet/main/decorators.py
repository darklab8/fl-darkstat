from django.contrib import messages
from django.core import management
from django.conf import settings


def loaded_db(f):
    def decorator_function(*args, **kwargs):
        #print('executing '+f.__name__)

        from ship.models import Ship
        if len(Ship.objects.all()) == 0:
            management.call_command('loaddata', 'dump.json')

        return f(*args, **kwargs)
    return decorator_function
