import time
from django.core import management
from django.conf import settings
from .files import main_parse


class DbHandler:
    "class for everything related to databases"
    def __init__(self):
        pass

    def load_from_dump(self, database):
        "load database from dump file"
        self.make_empty(database)

        management.call_command("loaddata", "dump.json", f"--database={database}")

    def make_empty(self, database):
        "clean database from data and apply migrations"
        management.call_command(
            "flush",
            "--noinput",
            f"--database={database}",
        )
        management.call_command("migrate", f"--database={database}")

    def save_to_dump(self, database):
        "save database to dump file"
        management.call_command(
            "dumpdata",
            f"--database={database}",
            natural_foreign=True,
            natural_primary=True,
            indent=2,
            output="dump.json",
        )

    def parse_to(self, database):
        "parser freelancer into db"
        self.make_empty(database)

        from commodities.models import fill_commodity_table
        from ship.models import fill_ship_table

        dicty = main_parse()

        fill_commodity_table(dicty, database)
        fill_ship_table(dicty, database)

    def parser_and_transfer(self):
        "parse to RAM memory and transfer to default database"
        self.parse_to("memory")
        self.save_to_dump("memory")
        self.load_from_dump("default")


    def daemon_database_update(self):
        "background forever process for auto updates"

        self.make_empty("default")
        # try:
        #     self.load_from_dump("default")
        # except Exception as error:
        #     print(error)
        #     breakpoint()

        sleeping_seconds = int(settings.TIMEOUT_BETWEEN_PARSE)
        while True:
            self.parser_and_transfer()
            print(f"sleeping for {sleeping_seconds} seconds")
            time.sleep(sleeping_seconds)