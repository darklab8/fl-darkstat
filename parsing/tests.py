"checking parsers work"
from django.test import TestCase
from django.conf import settings
from .apps import on_server_start

# Create your tests here.
class TestParsers(TestCase):
    "tests to check main program work, I guess that admin interface just logs for now"

    def setUp(self):
        settings.DARK_PARSE=True
        settings.FREELANCER_FOLDER="dark_copy"

    def test_main_url(self):
        "check parser work"
        on_server_start()
