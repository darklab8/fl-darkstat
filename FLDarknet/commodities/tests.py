from django.test import TestCase
from .models import Commodity
from main.decorators import loaded_db
import sys
from django.test import Client

# Create your tests here.


class Test_CommodityModel(TestCase):

    @loaded_db
    def setUp(self):
        self.client = Client()

    def test_validator_not_empty(self):
        count = len(Commodity.objects.all())
        print("Commodity =", count)
        self.assertIs(count != 0, True)

    def test_commodity_url(self):
        resp = self.client.get('/admin/commodities/', follow=True)
        self.assertEqual(resp.status_code, 200)

    def test_commodity_commodity_url(self):
        resp = self.client.get('/admin/commodities/commodity/', follow=True)
        self.assertEqual(resp.status_code, 200)

    def test_ship_ship_change_url(self):
        resp = self.client.get('/admin/commodities/commodity/1/change/', follow=True)
        self.assertEqual(resp.status_code, 200)
