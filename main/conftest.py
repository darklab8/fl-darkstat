import pytest
from django.core import management
from ship.models import Ship


@pytest.fixture
def loaded_dump():
    if len(Ship.objects.all()) == 0:
        management.call_command('loaddata', 'dump.json')
