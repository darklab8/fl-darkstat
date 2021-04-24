"Testing user things are going to be here"
from django.test import TestCase
from django.test import Client


class TestUrls(TestCase):
    """tests to check main program work,
    I guess that admin interface just logs for now"""

    def setUp(self):
        self.client = Client()

    def test_main_url(self):
        "Checking what will happen after index rendering"
        resp = self.client.get('/', follow=True)
        self.assertEqual(resp.status_code, 200)

    def test_admin_url(self):
        "Checking admin interface rendering"
        resp = self.client.get('/admin/', follow=True)
        self.assertEqual(resp.status_code, 200)
