import pytest

import os
import sys
import django
from django.core import management
sys.path.insert(0, os.path.abspath('..'))
os.environ['DJANGO_SETTINGS_MODULE'] = 'core.settings'
django.setup()


from ship.models import Ship
from commodity.models import Commodity


@pytest.fixture
def loaded_dump():
    if len(Ship.objects.all()) == 0:
        management.call_command('loaddata', 'dump.json')


@pytest.fixture
def tables():
    return [Commodity, Ship]
