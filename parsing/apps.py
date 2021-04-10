"module to parse data"
import os
import sys
from threading import Thread
from django.apps import AppConfig
from .database import DbHandler


class ParsingConfig(AppConfig):
    "main class to launch server in a different configuration to parse/save/load data"

    name = "parsing"

    def ready(self):
        db_handler = DbHandler()
        if os.environ.get("RUN_MAIN", None) != "true":
            if "runserver" in sys.argv:
                thread_name = "django_1UGbackground_worker3"

                thread = Thread(
                    name=thread_name, target=db_handler.daemon_database_update, daemon=True, args=()
                )
                thread.start()

        # for filename in equipment.keys():
        #     for header in equipment[filename].keys():
        #         for obj in equipment[filename][header]:
        #             # print(obj)
        #             # breakpoint()
        #             break

        # goods = parse_file(settings.GOODS_DIR)
        # select_equip = parse_file(settings.SEL_EQUIP_DIR)
        # market_commodities = parse_file(settings.MARKET_DIR)


# from main.apps import *
# test = set()
# goods = ships['shiparch.ini']
# arr = goods['[ship]']
# for obj in arr:
#     for key in obj.keys():
#         if key not in test:
#             if len(obj[key]) > 1:# and not isinstance(obj[key][0],list):
#                 print(key, " = ", obj[key])
#             test.add(key)

# goods = equipment['select_equip.ini']
# arr = goods['[commodity]']
# for obj in arr:
#     if 'ids_name' in obj.keys():
#         ids = int(obj['ids_name'][0])
#         if ids in infocards:
#             print(infocards[ids])
#         else:
#             break
#             breakpoint()
#             print('ERR in test 1')

# possible_keys = set()
# for key in goods.keys():
#     # if 'category' in goods[key].keys() and goods[key]['category'] == 'equipment':
#     for subkey in goods[key].keys():
#         if 'category' in subkey:
#             possible_keys.add(goods[key][subkey])
# print(possible_keys)

# for key in goods.keys():
#     if 'commodity' in key:
#         print(goods[key])

# for key in goods.keys():
#     for subkey in goods[key].keys():
#         if isinstance(goods[key][subkey], list):
#             if 'addon' not in subkey:
#                 print('ERR not addon')
#
#
# for i in range(10):
#     c = Commodity(
#         name = str(i)
#     )
#     c.save()

# for obj in comms:
#     c = Commodity(
#         name = obj.name(),
#         nickname = obj.nickname,
#         ids_name = obj.ids_name,
#         ids_info = obj.ids_info,
#         lootable = obj.lootable,
#         decay_per_second = obj.decay_per_second,
#         volume = obj.volume,
#     )
#     c.save()
