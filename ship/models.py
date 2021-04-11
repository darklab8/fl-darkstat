"model for table about space ships"
import xmltodict
from django.db import models
from parsing.extracting import view_wrapper, view_wrapper_with_infocard
# Create your models here.


class Ship(models.Model):
    "Any related stats about space ship information"
    class Meta:
        verbose_name_plural = "ships"

    nickname = models.CharField(
        max_length=50, db_index=True, blank=True, null=True)

    ids_name = models.IntegerField(db_index=True, blank=True, null=True)
    name = models.CharField(
        max_length=50, db_index=True, blank=True, null=True)

    ids_info = models.IntegerField(db_index=True, blank=True, null=True)

    info_name = models.CharField(
        max_length=100, db_index=True, blank=True, null=True)

    mass = models.FloatField(blank=True, null=True)
    hold_size = models.IntegerField(blank=True, null=True)

    capacity = models.IntegerField(blank=True, null=True, default=None)
    charge_rate = models.IntegerField(blank=True, null=True, default=None)


    linear_drag = models.FloatField(blank=True, null=True)
    max_bank_angle = models.IntegerField(db_index=True, blank=True, null=True)
    camera_angular_acceleration = models.FloatField(
        blank=True, null=True, verbose_name='ang speed')

    camera_horizontal_turn_angle = models.IntegerField(
        blank=True, null=True, verbose_name='hor ang')

    camera_vertical_turn_up_angle = models.IntegerField(
        blank=True, null=True, verbose_name='turn up')

    camera_vertical_turn_down_angle = models.IntegerField(
        blank=True, null=True, verbose_name='turn down')

    camera_turn_look_ahead_slerp_amount = models.FloatField(
        blank=True, null=True, verbose_name='look ahead')

    hit_pts = models.IntegerField(blank=True, null=True)
    nudge_force = models.FloatField(blank=True, null=True)
    strafe_force = models.IntegerField(blank=True, null=True)
    strafe_power_usage = models.IntegerField(
        blank=True, null=True, verbose_name='strafe usage')
    explosion_resistance = models.FloatField(
        blank=True, null=True, verbose_name='exp res')

    ids_info1 = models.IntegerField(blank=True, null=True)
    ids_info2 = models.IntegerField(blank=True, null=True)
    ids_info3 = models.IntegerField(blank=True, null=True)
    ship_class = models.IntegerField(blank=True, null=True)
    nanobot_limit = models.IntegerField(
        blank=True, null=True, verbose_name='nanobots')
    shield_battery_limit = models.IntegerField(
        blank=True, null=True, verbose_name='batteries')

    # different original name: type
    typeof = models.CharField(
        max_length=50, db_index=True, blank=True, null=True)


def fill_ship_table(dicty, database):
    """Filling ship database with data from universe"""
    goods = dicty.ships["shiparch.ini"]
    arr = goods["[ship]"].copy()
    for i, obj in enumerate(arr):

        kwg = {}
        view_wrapper(kwg, obj, str, "nickname")
        view_wrapper_with_infocard(dicty, kwg, obj, int, "ids_name", "name")
        view_wrapper(kwg, obj, int, "ids_info")
        view_wrapper(kwg, obj, float, "mass")
        view_wrapper(kwg, obj, int, "hold_size")
        view_wrapper(kwg, obj, float, "linear_drag")
        view_wrapper(kwg, obj, int, "max_bank_angle")
        view_wrapper(kwg, obj, float, "camera_angular_acceleration")
        view_wrapper(kwg, obj, int, "camera_horizontal_turn_angle")
        view_wrapper(kwg, obj, int, "camera_vertical_turn_up_angle")
        view_wrapper(kwg, obj, int, "camera_vertical_turn_down_angle")
        view_wrapper(kwg, obj, float, "camera_turn_look_ahead_slerp_amount")
        view_wrapper(kwg, obj, int, "hit_pts")
        view_wrapper(kwg, obj, float, "nudge_force")
        view_wrapper(kwg, obj, int, "strafe_force")
        view_wrapper(kwg, obj, int, "strafe_power_usage")
        view_wrapper(kwg, obj, float, "explosion_resistance")
        view_wrapper(kwg, obj, int, "ids_info1")
        view_wrapper(kwg, obj, int, "ids_info2")
        view_wrapper(kwg, obj, int, "ids_info3")
        view_wrapper(kwg, obj, int, "ship_class")
        view_wrapper(kwg, obj, int, "nanobot_limit")
        view_wrapper(kwg, obj, int, "shield_battery_limit")

        if "type" in obj.keys():
            kwg["typeof"] = str(obj["type"][0])

        if "nickname" in obj.keys() and "hp_type" in obj.keys():
            dicty.hp_type[obj["nickname"][0]] = obj["hp_type"]

        try:
            dic = xmltodict.parse(dicty.infocards[kwg["ids_info"]][1])["RDL"]["TEXT"]
            if not dic[0]:
                dic = xmltodict.parse(dicty.infocards[kwg["ids_info1"]][1])["RDL"]["TEXT"]
            kwg["info_name"] = dic[0]
        except KeyError:
            print(
                "ERR not able to find infocard for ship object #",
                i,
                " ship nickname",
                kwg.get("nickname", "no nickname"),
                "ids_info",
                kwg.get("ids_info", "no ids_info"),
            )
        except xmltodict.expat.ExpatError:
            print(
                "ERR xmltodict.expat.ExpatError, can't parse infocard xml #",
                i,
                " ship nickname",
                kwg.get("nickname", "no nickname"),
                "ids_info",
                kwg.get("ids_info", "no ids_info"),
            )

        if kwg["nickname"] in dicty.goods_by_ship["shiphull"]:
            hull = dicty.goods_by_ship["shiphull"][kwg["nickname"]]["nickname"][0]
            try:
                ship = dicty.goods_by_hull["ship"][hull]
            except KeyError:
                print("ERR no package in goods.ini for ship hull =", hull)
            for addon in ship['addon']:
                power_nickname = addon[0]
                if power_nickname in dicty.misc_equip_power_by_nickname:

                    powercore = dicty.misc_equip_power_by_nickname[power_nickname]

                    kwg["capacity"] = int(powercore['capacity'][0])
                    kwg["charge_rate"]= int(powercore['charge_rate'][0])
            
            # print('123')
            # for add in ship['addon']
            # i f add[0] in dicty.equipment['misc_equip.ini']['[power]'].keys()
            # TODO find in addons powercore st_equip
            # and perhaps engine in engine_equip

        db_data = Ship(**kwg)
        db_data.save(using=database)
        # except Exception as error:
        # print("ERR in filling ship #", i)