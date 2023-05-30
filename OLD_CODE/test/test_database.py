"Testing is going to be here in order to be DRY"
import pytest
pytestmark = pytest.mark.django_db


def test_for_not_empty_table(loaded_dump, tables):
    for table in tables:
        assert len(table.objects.all()) != 0
