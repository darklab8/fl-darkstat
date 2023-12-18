import time
from django.core import management
from django.conf import settings
from .parser import main_parse


class DbHandler:
    "class for everything related to databases"

    def __init__(self):
        pass

    def load_from_dump(self, database):
        "load database from dump file"
        self.make_empty(database)

        management.call_command("loaddata", "dump.json",
                                f"--database={database}")

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
            '--exclude=contenttypes',
            '--exclude=auth',
            natural_foreign=True,
            natural_primary=True,
            indent=2,
            output="dump.json",
        )

    def parse_to(self, database_name):
        "parser freelancer into db"
        self.make_empty(database_name)

        from commodity.models import Commodity
        from ship.models import Ship

        parsed = main_parse()

        Commodity.fill_table(
            commodities=parsed.equipment.select_equip.commodity,
            infocards=parsed.infocards,
            database_name=database_name)

        Ship.fill_table(
            ships=parsed.ships.shiparch.ship,
            infocards=parsed.infocards,
            goods_by_ship={
                item['ship'][0]: item
                for item in parsed.equipment.goods.good
                if 'ship' in item and 'shiphull' in item["category"][0]
            },
            goods_by_hull={
                item['hull'][0]: item
                for item in parsed.equipment.goods.good
                if 'hull' in item and 'ship' in item["category"][0]
            },
            power={
                item['nickname'][0]: item
                for item in parsed.equipment.misc_equip.power
                if 'nickname' in item
            },
            engines={
                item['nickname'][0]: item
                for item in parsed.equipment.engine_equip.engine
                if 'nickname' in item
            },
            database_name=database_name)

    def parser_and_transfer(self):
        "parse to RAM memory and transfer to default database"
        self.parse_to("memory")
        self.save_to_dump("memory")
        self.load_from_dump("default")

    def daemon_database_update(self):
        "background forever process for auto updates"

        self.make_empty("default")

        sleeping_seconds = settings.TIMEOUT_BETWEEN_PARSE
        while True:
            self.parser_and_transfer()
            print(f"sleeping for {sleeping_seconds} seconds")
            time.sleep(sleeping_seconds)
