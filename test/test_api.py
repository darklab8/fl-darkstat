"Testing is going to be here in order to be DRY"
from django.conf import settings
from rest_framework.test import APIClient
import pytest
from main.decorators import record_to_docs


@pytest.mark.django_db
@pytest.mark.parametrize("app", settings.ADDED_APPS)
def test_api_to_retrieve_tables(loaded_dump, app):
    resp = APIClient().get(f"/{app}/list", format='json')

    assert (len(resp.json())) > 0, app

    record_to_docs(app, "list.json", resp.json())
