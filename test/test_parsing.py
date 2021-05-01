"checking parsers work"
from parsing.database import DbHandler
import pytest
# Create your tests here.


@pytest.mark.django_db
def test_to_check_parser(tables):
    DbHandler().parser_and_transfer()

    for table in tables:
        assert len(table.objects.all()) != 0
