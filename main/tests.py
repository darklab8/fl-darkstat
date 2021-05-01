"Testing is going to be here in order to be DRY"
from django.test import Client
from django.conf import settings
from rest_framework.test import APIClient
import pytest
from .decorators import record_to_docs


@pytest.mark.django_db
@pytest.mark.parametrize("url", ['/', '/admin/'])
def test_main_url_render(loaded_dump, url):
    assert Client().get(url, follow=True).status_code == 200


@pytest.mark.django_db
@pytest.mark.parametrize("app", settings.ADDED_APPS)
def test_admin_section_render(loaded_dump, app):
    assert Client().get(f"/admin/{app}/", follow=True).status_code == 200


@pytest.mark.django_db
@pytest.mark.parametrize("app", settings.ADDED_APPS)
def test_admin_table_render(loaded_dump, app):
    assert Client().get(f"/admin/{app}/{app}/", follow=True).status_code == 200


@pytest.mark.django_db
@pytest.mark.parametrize("app", settings.ADDED_APPS)
def test_admin_table_change_render(loaded_dump, app):
    assert Client().get(
        f"/admin/{app}/{app}/1/change/", follow=True).status_code == 200


@pytest.mark.django_db
def test_for_not_empty_table(loaded_dump, tables):
    for table in tables:
        assert len(table.objects.all()) != 0


@pytest.mark.django_db
@pytest.mark.parametrize("app", settings.ADDED_APPS)
def test_api_to_retrieve_tables(loaded_dump, app):
    resp = APIClient().get(f"/{app}/list", format='json')

    assert (len(resp.json())) > 0, app

    record_to_docs(app, "list.json", resp.json())
