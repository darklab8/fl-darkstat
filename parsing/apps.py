"module to parse data"
import os
from threading import Thread
from django.apps import AppConfig
from django.conf import settings
from .database import DbHandler


class ParsingConfig(AppConfig):
    """main class to launch server in a
    different configuration to parse/save/load data"""

    name = "parsing"

    def ready(self):
        db_handler = DbHandler()
        if (os.environ.get("RUN_MAIN", None) != "true"
                and settings.BACKGROUND_PROCESS):
            thread_name = "django_1UGbackground_worker3"

            thread = Thread(
                name=thread_name,
                target=db_handler.daemon_database_update,
                daemon=True,
                args=(),
            )
            thread.start()
