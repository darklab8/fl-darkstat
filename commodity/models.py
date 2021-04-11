""""Section for regular shop commodities"""
from django.db import models
from parsing.extracting import view_wrapper, view_wrapper_with_infocard


class Commodity(models.Model):
    """"Model for freelancer ingame regular commodities in the shop"""

    class Meta:
        verbose_name_plural = "commodities"

    ids_name = models.IntegerField(db_index=True, blank=True, null=True)
    name = models.CharField(
        max_length=50, db_index=True, blank=True, null=True)

    ids_info = models.IntegerField(db_index=True, blank=True, null=True)

    units_per_container = models.IntegerField(
        db_index=True, blank=True, null=True)
    decay_per_second = models.IntegerField(blank=True, null=True)
    hit_pts = models.IntegerField(blank=True, null=True)

    volume = models.FloatField(blank=True, null=True)

    loot_appearance = models.CharField(
        max_length=50, db_index=True, blank=True, null=True)
    pod_appearance = models.CharField(
        max_length=50, db_index=True, blank=True, null=True)

    nickname = models.CharField(
        max_length=50, db_index=True, blank=True, null=True)


def fill_commodity_table(dicty, database):
    """Filling our commodity database section with data"""
    goods = dicty.equipment["select_equip.ini"]
    arr = goods["[commodity]"].copy()
    for obj in arr:
        kwg = {}
        view_wrapper_with_infocard(dicty, kwg, obj, int, "ids_name", "name")
        view_wrapper(kwg, obj, int, "ids_info")

        view_wrapper(kwg, obj, int, "units_per_container")
        view_wrapper(kwg, obj, int, "decay_per_second")
        view_wrapper(kwg, obj, int, "hit_pts")

        view_wrapper(kwg, obj, str, "pod_appearance")
        view_wrapper(kwg, obj, str, "loot_appearance")
        view_wrapper(kwg, obj, str, "nickname")

        view_wrapper(kwg, obj, float, "volume")

        db_data = Commodity(**kwg)
        db_data.save(using=database)