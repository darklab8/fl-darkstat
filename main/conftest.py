import pytest
from django.core import management
from ship.models import Ship
from commodity.models import Commodity


@pytest.fixture
def loaded_dump():
    if len(Ship.objects.all()) == 0:
        management.call_command('loaddata', 'dump.json')


@pytest.fixture
def tables():
    return [Commodity, Ship]
