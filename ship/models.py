"model for table about space ships"
import xmltodict
from django.db import models
from parsing.extracting import view_wrapper_with_infocard, add_to_model
from django.db.models import (
    CharField,
    IntegerField,
    FloatField,
    OneToOneField,
    CASCADE,
)
from types import SimpleNamespace
# Create your models here.


class Piece(models.Model):
    pass


class ShipStrFloats(Piece):
    ship_float_piece = OneToOneField(Piece,
                                     on_delete=CASCADE,
                                     parent_link=True)

    # str
    nickname: IntegerField = CharField(max_length=50,
                                       db_index=True,
                                       blank=True,
                                       null=True)

    # floats
    mass = FloatField(blank=True, null=True)
    linear_drag = FloatField(blank=True, null=True)
    camera_angular_acceleration = FloatField(blank=True,
                                             null=True,
                                             verbose_name="ang speed")
    camera_turn_look_ahead_slerp_amount = FloatField(blank=True,
                                                     null=True,
                                                     verbose_name="look ahead")
    nudge_force = FloatField(blank=True, null=True)
    explosion_resistance = FloatField(blank=True,
                                      null=True,
                                      verbose_name="exp res")


class ShipIntegers(Piece):
    ship_integer_piece = OneToOneField(Piece,
                                       on_delete=CASCADE,
                                       parent_link=True)

    # int
    ids_info = IntegerField(db_index=True, blank=True, null=True)
    hold_size = IntegerField(blank=True, null=True)
    max_bank_angle = IntegerField(db_index=True, blank=True, null=True)
    camera_horizontal_turn_angle = IntegerField(blank=True,
                                                null=True,
                                                verbose_name="hor ang")

    camera_vertical_turn_up_angle = IntegerField(blank=True,
                                                 null=True,
                                                 verbose_name="turn up")

    camera_vertical_turn_down_angle = IntegerField(blank=True,
                                                   null=True,
                                                   verbose_name="turn down")
    hit_pts = IntegerField(blank=True, null=True)
    strafe_force = IntegerField(blank=True, null=True)
    strafe_power_usage = IntegerField(blank=True,
                                      null=True,
                                      verbose_name="strafe usage")

    ids_info1 = IntegerField(blank=True, null=True)
    ids_info2 = IntegerField(blank=True, null=True)
    ids_info3 = IntegerField(blank=True, null=True)
    ship_class = IntegerField(blank=True, null=True)
    nanobot_limit = IntegerField(blank=True,
                                 null=True,
                                 verbose_name="nanobots")
    shield_battery_limit = IntegerField(blank=True,
                                        null=True,
                                        verbose_name="batteries")


class ShipSpecial(Piece):
    ship_special_piece = OneToOneField(Piece,
                                       on_delete=CASCADE,
                                       parent_link=True)

    # SPECIAL
    ids_name = IntegerField(db_index=True, blank=True, null=True)
    name = CharField(max_length=50, db_index=True, blank=True, null=True)

    # different original name: type
    typeof = CharField(max_length=50, db_index=True, blank=True, null=True)

    info_name = CharField(max_length=100, db_index=True, blank=True, null=True)

    # powercore
    capacity = IntegerField(blank=True, null=True, default=None)
    charge_rate = IntegerField(blank=True, null=True, default=None)

    # engine
    cruise_speed = IntegerField(blank=True, null=True, default=None)
    impulse_speed = IntegerField(blank=True, null=True, default=None)


class Ship(ShipStrFloats, ShipIntegers, ShipSpecial):
    "Any related stats about space ship information"

    class Meta:
        verbose_name_plural = "ships"

    @classmethod
    def fill_table(cls, ships, infocards, good_original, power, engines,
                   database_name):
        """Filling ship database with data from universe"""
        # arranging goods
        goods = SimpleNamespace(
            by_ship={
                item['ship'][0]: item
                for item in good_original
                if 'ship' in item and 'shiphull' in item["category"][0]
            },
            by_hull={
                item['hull'][0]: item
                for item in good_original
                if 'hull' in item and 'ship' in item["category"][0]
            })

        for i, obj in enumerate(ships):

            kwg = {}

            add_to_model(kwg, obj, str, ("nickname", ))

            add_to_model(
                kwg,
                obj,
                float,
                (
                    "mass",
                    "linear_drag",
                    "camera_angular_acceleration",
                    "camera_turn_look_ahead_slerp_amount",
                    "nudge_force",
                    "explosion_resistance",
                ),
            )

            add_to_model(
                kwg,
                obj,
                int,
                (
                    "ship_class",
                    "hold_size",
                    "nanobot_limit",
                    "hit_pts",
                    "shield_battery_limit",
                    "ids_info",
                    "max_bank_angle",
                    "camera_horizontal_turn_angle",
                    "camera_vertical_turn_up_angle",
                    "camera_vertical_turn_down_angle",
                    "strafe_force",
                    "strafe_power_usage",
                    "ids_info1",
                    "ids_info2",
                    "ids_info3",
                ),
            )

            # ids_name + name from infocards
            view_wrapper_with_infocard(infocards, kwg, obj, int, "ids_name",
                                       "name")

            # add typeoff
            if "type" in obj:
                kwg["typeof"] = str(obj["type"][0])

            # This is ship gun/equipment slots! Add them to some database!
            # if "nickname" in obj and "hp_type" in obj:
            #     hp_type[obj["nickname"][0]] = obj["hp_type"]
            # TODO move to parser and parse later for available equip slots

            # add name of the ship from infocard's beginning
            try:
                dic = xmltodict.parse(
                    infocards[kwg["ids_info"]][1])["RDL"]["TEXT"]
                if not dic[0]:
                    dic = xmltodict.parse(
                        infocards[kwg["ids_info1"]][1])["RDL"]["TEXT"]
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
                    "ERR xmltodict.expat.ExpatError, "
                    "can't parse infocard xml #",
                    i,
                    " ship nickname",
                    kwg.get("nickname", "no nickname"),
                    "ids_info",
                    kwg.get("ids_info", "no ids_info"),
                )

            # add powercore parameters and engine parameters
            if kwg["nickname"] in goods.by_ship:
                hull = goods.by_ship[kwg["nickname"]]["nickname"][0]
                try:
                    ship = goods.by_hull[hull]
                except KeyError:
                    print("ERR no package in goods.ini for ship hull =", hull)
                for addon in ship["addon"]:
                    addon_nickname = addon[0]
                    if addon_nickname in power.by_nickname:

                        powercore = \
                            power.by_nickname[addon_nickname]

                        kwg["capacity"] = int(powercore["capacity"][0])
                        kwg["charge_rate"] = int(powercore["charge_rate"][0])

                    elif addon_nickname in engines.by_nickname:

                        engine = engines.by_nickname[addon_nickname]

                        try:
                            kwg["cruise_speed"] = int(
                                engine["cruise_speed"][0])
                        except KeyError:
                            kwg["cruise_speed"] = 350

                        kwg["impulse_speed"] = int(
                            float(engine["max_force"][0]) /
                            float(engine["linear_drag"][0]))
                        # breakpoint()
                        # print("ERR no cruise_speed for \
                        #   ship hull = ", hull, " ", kwg['nickname'])

                # print('123')
                # for add in ship['addon']
                # i f add[0] in dicty.equipment['misc_equip.ini']['[power]']
                # TODO find in addons powercore st_equip
                # and perhaps engine in engine_equip

            db_data = cls(**kwg)
            db_data.save(using=database_name)
            # except Exception as error:
            # print("ERR in filling ship #", i)
