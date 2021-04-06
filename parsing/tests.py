"checking parsers work"
from django.test import TestCase
from django.conf import settings
from .apps import on_server_start

# Create your tests here.
class TestParsers(TestCase):
    "tests to check main program work, I guess that admin interface just logs for now"

    def test_main_url(self):
        "check parser work"
        settings.DARK_PARSE=True
        settings.PATHS.redefine_folder("dark_copy")
        on_server_start()
