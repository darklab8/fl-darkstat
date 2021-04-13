"Tests to cover anything related to ship table"
from django.test import TestCase
from django.test import Client
from main.decorators import loaded_db
from .models import Ship


class TestShipModel(TestCase):
    "Tests to check ship model"

    @loaded_db
    def setUp(self):
        self.client = Client()

    def test_validator_not_empty(self):
        "test to check, if objects were loaded to db"
        count = len(Ship.objects.all())
        print("Ships =", count)
        self.assertIs(count != 0, True)

    def test_ship_url(self):
        """test to check rendering of ship section"""
        resp = self.client.get("/admin/ship/", follow=True)
        self.assertEqual(resp.status_code, 200)

    def test_ship_ship_url(self):
        """Test to check rendering of main ship table"""
        resp = self.client.get("/admin/ship/ship/", follow=True)
        self.assertEqual(resp.status_code, 200)

    def test_ship_ship_change_url(self):
        """Test to check rendering of inline admin interface"""
        resp = self.client.get("/admin/ship/ship/1/change/", follow=True)
        self.assertEqual(resp.status_code, 200)


class TestShipAPI(TestCase):
    """Tests to check db model commodity"""

    @loaded_db
    def setUp(self):
        pass

    def test_check_json_response_is_not_empty(self):
        self.client = Client()
        resp = self.client.get("/api/ship/?format=json", follow=True)
        assert (len(resp.json())) > 0
