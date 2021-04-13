""""Section for regular shop commodities"""
from django.db import models
from parsing.extracting import view_wrapper_with_infocard, add_to_model


class Commodity(models.Model):
    """"Model for freelancer ingame regular commodities in the shop"""

    class Meta:
        verbose_name_plural = "commodities"

    # str
    nickname = models.CharField(max_length=50, db_index=True, blank=True, null=True)
    loot_appearance = models.CharField(
        max_length=50, db_index=True, blank=True, null=True
    )
    pod_appearance = models.CharField(
        max_length=50, db_index=True, blank=True, null=True
    )

    # floats
    volume = models.FloatField(blank=True, null=True)

    # int
    ids_info = models.IntegerField(db_index=True, blank=True, null=True)

    units_per_container = models.IntegerField(db_index=True, blank=True, null=True)
    decay_per_second = models.IntegerField(blank=True, null=True)
    hit_pts = models.IntegerField(blank=True, null=True)

    # SPECIAL
    ids_name = models.IntegerField(db_index=True, blank=True, null=True)
    name = models.CharField(max_length=50, db_index=True, blank=True, null=True)


def fill_commodity_table(dicty, database):
    """Filling our commodity database section with data"""
    goods = dicty.equipment["select_equip.ini"]
    arr = goods["[commodity]"].copy()
    for obj in arr:
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

        view_wrapper_with_infocard(dicty, kwg, obj, int, "ids_name", "name")

        db_data = Commodity(**kwg)
        db_data.save(using=database)