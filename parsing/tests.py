"checking parsers work"
import os
from parsing.database import DbHandler
import pytest
pytestmark = pytest.mark.django_db

@pytest.mark.skipif(not bool(os.environ.get("pipline")), reason="long test")
def test_to_check_parser(tables):
    DbHandler().parser_and_transfer()

    for table in tables:
        assert len(table.objects.all()) != 0
