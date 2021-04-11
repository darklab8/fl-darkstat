"""Module to test model commodity"""
import json
from django.test import (
    TestCase,
    Client,
)
from main.decorators import loaded_db
from .models import Commodity
# Create your tests here.

class TestCommodityUrls(TestCase):
    """Tests to check db model commodity"""

    @loaded_db
    def setUp(self):
        self.client = Client()

    def test_validator_not_empty(self):
        """Checking if objects were able to load into db"""
        count = len(Commodity.objects.all())
        print("Commodity =", count)
        self.assertIs(count != 0, True)

    def test_commodity_url(self):
        """Checking main section url loading"""
        resp = self.client.get('/admin/commodity/', follow=True)
        self.assertEqual(resp.status_code, 200)

    def test_commodity_commodity_url(self):
        """Checking if table is able to load for view"""
        resp = self.client.get('/admin/commodity/commodity/', follow=True)
        self.assertEqual(resp.status_code, 200)

    def test_ship_ship_change_url(self):
        """"Checking if inline data is loading correctly"""
        resp = self.client.get('/admin/commodity/commodity/1/change/', follow=True)
        self.assertEqual(resp.status_code, 200)

class TestCommodityModel(TestCase):
    """Tests to check db model commodity"""

    @loaded_db
    def setUp(self):
        pass

    def test_check_json_response_is_not_empty(self):
        self.client = Client()
        resp = self.client.get('/api/commodity/?format=json', follow=True)
        assert (len(resp.json())) > 0

