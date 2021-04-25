"checking parsers work"
from django.test import TestCase
from .database import DbHandler
from main.decorators import suppress_stdout

# Create your tests here.


class TestParsers(TestCase):
    """tests to check main program work,
    I guess that admin interface just logs for now"""

    def test_main_url(self):
        "check parser work"
        with suppress_stdout():
            db_handler = DbHandler()
            db_handler.parse_to("default")
            db_handler.save_to_dump("default")
            db_handler.load_from_dump("default")
