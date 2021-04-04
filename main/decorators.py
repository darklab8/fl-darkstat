"some universally used decorators will be here"
from django.core import management
from ship.models import Ship


def loaded_db(func):
    "useful decorator to preload stuff for unit tests"
    def decorator_function(*args, **kwargs):
        # print('executing '+func.__name__)

        if len(Ship.objects.all()) == 0:
            management.call_command('loaddata', 'dump.json')

        return func(*args, **kwargs)
    return decorator_function
