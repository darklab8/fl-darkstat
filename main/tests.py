"Testing user things are going to be here"
import json
import os
from django.test import TestCase
from django.test import Client
from main.decorators import loaded_db
from django.conf import settings
from rest_framework.test import APIClient
from commodity.models import Commodity
from ship.models import Ship


class TestAdminUrl(TestCase):
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


class TestTemplateUrl(TestCase):
    """Tests to check db model commodity"""

    @loaded_db
    def setUp(self):
        pass

    def test_validator_not_empty(self):
        def check_validator_not_empty(app, class_):
            """Checking if objects were able to load into db"""
            count = len(class_.objects.all())
            assert count != 0, app

        check_validator_not_empty("commodity", Commodity)
        check_validator_not_empty("ship", Ship)

    def test_app_url(self):
        def check_app_url(app):
            """Checking main section url loading"""
            resp = Client().get(f"/admin/{app}/", follow=True)
            assert resp.status_code == 200, app

        for app in settings.ADDED_APPS:
            check_app_url(app)

    def test_app_app_url(self):
        def check_app_app_url(app):
            """Checking if table is able to load for view"""
            resp = Client().get(f"/admin/{app}/{app}/", follow=True)
            assert resp.status_code == 200, app

        for app in settings.ADDED_APPS:
            check_app_app_url(app)

    def test_app_app_change_url(self):
        def check_app_app_change_url(app):
            """"Checking if inline data is loading correctly"""
            resp = Client().get(
                f"/admin/{app}/{app}/1/change/", follow=True)
            assert resp.status_code == 200, app

        for app in settings.ADDED_APPS:
            check_app_app_change_url(app)


class TestTemplateAPI(TestCase):
    """Tests to check db model"""

    @loaded_db
    def setUp(self):
        pass

    def test_json_response_is_not_empty(self):
        def check_json_response_is_not_empty(app):
            resp = APIClient().get(f"/{app}/list", format='json')
            assert (len(resp.json())) > 0, app

            if settings.REFRESH_EXAMPLES:
                with open(os.path.join(
                    'sphinx',
                    'source',
                    f'{app}',
                    'write',
                    'list.json'
                ), 'w') as file_:
                    file_.write(json.dumps(resp.json(), indent=2))

        for app in settings.ADDED_APPS:
            check_json_response_is_not_empty(app)
