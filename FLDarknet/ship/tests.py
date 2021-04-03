from django.test import TestCase
from .models import Ship
from main.decorators import loaded_db
import sys
from django.test import Client

# Create your tests here.

class Test_ShipModel(TestCase):
    
    @loaded_db
    def setUp(self):
        pass

    def tearDown(self):
        pass

    def test_validator_not_empty(self):
        count = len(Ship.objects.all())
        print("Ships =", count)
        self.assertIs(count != 0, True)

    def test_main_url(self):
        self.client = Client()
        resp = self.client.get('/', follow=True)
        self.assertEqual(resp.status_code, 200)

    def test_admin_url(self):
        self.client = Client()
        resp = self.client.get('/admin/', follow=True)
        self.assertEqual(resp.status_code, 200)

    def test_ship_url(self):
        self.client = Client()
        resp = self.client.get('/admin/ship/', follow=True)
        self.assertEqual(resp.status_code, 200)

    def test_ship_ship_url(self):
        self.client = Client()
        resp = self.client.get('/admin/ship/ship/', follow=True)
        self.assertEqual(resp.status_code, 200)

    def test_ship_ship_change_url(self):
        self.client = Client()
        resp = self.client.get('/admin/ship/ship/1/change/', follow=True)
        self.assertEqual(resp.status_code, 200)


