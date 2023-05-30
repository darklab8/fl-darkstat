""""Section for regular shop commodities"""
from django.db import models
from parsing.extracting import view_wrapper_with_infocard, add_to_model
from django.db.models import (
    CharField,
    IntegerField,
    FloatField,
)


class Commodity(models.Model):
    """"Model for freelancer ingame regular commodities in the shop"""

    class Meta:
        verbose_name_plural = "commodities"

    # str
    nickname = CharField(
        max_length=50, db_index=True, blank=True, null=True)
    f"{nickname.__repr__()}"
    loot_appearance = CharField(
        max_length=50, db_index=True, blank=True, null=True
    )
    pod_appearance = CharField(
        max_length=50, db_index=True, blank=True, null=True
    )

    # floats
    volume = FloatField(blank=True, null=True)

    # int
    ids_info = IntegerField(db_index=True, blank=True, null=True)

    units_per_container = IntegerField(
        db_index=True, blank=True, null=True)
    decay_per_second = IntegerField(blank=True, null=True)
    hit_pts = IntegerField(blank=True, null=True)

    # SPECIAL
    ids_name = IntegerField(db_index=True, blank=True, null=True)
    name = CharField(
        max_length=50, db_index=True, blank=True, null=True)

    @classmethod
    def fill_table(cls, commodities, infocards, database_name):
        """Filling our commodity database section with data"""
        for obj in commodities:
            kwg = {}

            add_to_model(
                kwg,
                obj,
                str,
                (
                    "nickname",
                    "loot_appearance",
                    "pod_appearance",
                ),
            )

            add_to_model(kwg, obj, float, ("volume",))

            add_to_model(
                kwg,
                obj,
                int,
                (
                    "ids_info",
                    "units_per_container",
                    "decay_per_second",
                    "hit_pts",
                ),
            )

            view_wrapper_with_infocard(
                infocards, kwg, obj, int, "ids_name", "name")

            db_data = cls(**kwg)
            db_data.save(using=database_name)
