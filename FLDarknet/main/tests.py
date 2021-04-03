from django.test import TestCase
import sys
from django.test import Client

# Create your tests here.


class Test_Main(TestCase):

    def setUp(self):
        self.client = Client()

    def test_main_url(self):
        resp = self.client.get('/', follow=True)
        self.assertEqual(resp.status_code, 200)

    def test_admin_url(self):
        resp = self.client.get('/admin/', follow=True)
        self.assertEqual(resp.status_code, 200)