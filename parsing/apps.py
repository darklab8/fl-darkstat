"module to parse data"
import os
import re
from os import walk
import xmltodict
from django.apps import AppConfig
from django.conf import settings
from django.core import management
from .files import (
    read_regular_file,
    read_utf8_file,
    clean_folder_from_files,
)


class Universe:
    """"This class will have parsed data from files
    stored in preparation to be sent to database"""
    equipment = {}
    universe = {}
    infocards = {}
    ships = {}
    hp_type = {}

    goods_by_nickname = {}
    goods_by_ship = {}
    goods_by_hull = {}




u = Universe()


def strip_from_rn(line):
    """Strips string from \r or \n trash"""
    return line.replace("\r", "").replace("\n", "")


def parse_infocards(filename):
    """"Parses infocard file into dictionary"""
    content = read_utf8_file(filename)

    regex_numbers = r"^\d+\r|^\d+\n"
    output = {}
    line_count = len(content)
    for i in range(line_count):
        if re.search(regex_numbers, content[i]) is not None:
            output[int(strip_from_rn(content[i]))] = [
                strip_from_rn(content[i+1]), strip_from_rn(content[i+2])]

    i += 1
    return output


def parse_file(filename):
    """Parses file into dictionary"""
    content = read_regular_file(filename)

    output = {}
    regex_for_headers = r"(\[)\w+(\])"

    line_count = len(content)
    for i in range(line_count):
        if re.search(regex_for_headers, content[i]) is not None:
            header = content[i].lower().replace("\n", "")

            if header not in output.keys():
                output[header] = []

            i += 1
            obj = {}
            while (
                (re.search(regex_for_headers, content[i+1]) is None) and
                (i < (line_count - 2))
            ):

                if content[i] == '\n':
                    i += 1
                    continue

                if re.search("^(;)", content[i]) is not None:
                    i += 1
                    continue

                splitted = content[i].replace(
                    " ", "").replace("\n", "").split('=')

                if len(splitted) == 2:
                    if splitted[0] not in obj.keys():
                        obj[splitted[0]] = []

                    if ',' not in splitted[1]:
                        obj[splitted[0]].append(splitted[1])
                    else:
                        obj[splitted[0]].append(splitted[1].split(','))
                else:
                    print('ERR splitted')
                i += 1

            output[header].append(obj)

        i += 1
    return output


def view_wrapper(kwg, obj, data_type, name):
    """"Function which prepares one value to be inserted into database"""
    if name in obj.keys():
        try:
            kwg[name] = data_type(obj[name][0])
        except ValueError as value_error_1:
            if data_type is int:
                if ";" in obj[name][0]:
                    splitted = obj[name][0].split(";")[0]
                    kwg[name] = data_type(splitted)
                elif '.' in obj[name][0]:
                    kwg[name] = int(float(obj[name][0]))
                else:
                    raise ValueError from value_error_1
            else:
                raise ValueError from value_error_1


def view_wrapper_with_infocard(kwg, obj, data_type, name, infoname):
    """Function that prepares two values to be inserted into database
    with getting extra one in infocard"""
    if name in obj.keys():
        try:
            kwg[name] = data_type(obj[name][0])
        except ValueError as value_error_1:
            if data_type is int and ";" in obj[name][0]:
                splitted = obj[name][0].split(";")[0]
                kwg[name] = data_type(splitted)
            else:
                raise ValueError from value_error_1
        if kwg[name] in u.infocards:
            kwg[infoname] = (u.infocards[kwg[name]][1])



def fill_commodity_table(model_obj):
    """Filling our commodity database section with data"""
    goods = u.equipment['select_equip.ini']
    arr = goods['[commodity]'].copy()
    for obj in arr:
        kwg = {}
        view_wrapper_with_infocard(kwg, obj, int, 'ids_name', 'name')
        view_wrapper(kwg, obj, int, 'ids_info')

        view_wrapper(kwg, obj, int, 'units_per_container')
        view_wrapper(kwg, obj, int, 'decay_per_second')
        view_wrapper(kwg, obj, int, 'hit_pts')

        view_wrapper(kwg, obj, str, 'pod_appearance')
        view_wrapper(kwg, obj, str, 'loot_appearance')
        view_wrapper(kwg, obj, str, 'nickname')

        view_wrapper(kwg, obj, float, 'volume')

        db_data = model_obj(
            **kwg
        )
        db_data.save()

def fill_ship_table(model_obj):
    """Filling ship database with data from universe"""
    goods = u.ships['shiparch.ini']
    arr = goods['[ship]'].copy()
    for i, obj in enumerate(arr):

        kwg = {}
        view_wrapper(kwg, obj, str, 'nickname')
        view_wrapper_with_infocard(kwg, obj, int, 'ids_name', 'name')
        view_wrapper(kwg, obj, int, 'ids_info')
        view_wrapper(kwg, obj, float, 'mass')
        view_wrapper(kwg, obj, int, 'hold_size')
        view_wrapper(kwg, obj, float, 'linear_drag')
        view_wrapper(kwg, obj, int, 'max_bank_angle')
        view_wrapper(kwg, obj, float, 'camera_angular_acceleration')
        view_wrapper(kwg, obj, int, 'camera_horizontal_turn_angle')
        view_wrapper(kwg, obj, int, 'camera_vertical_turn_up_angle')
        view_wrapper(kwg, obj, int, 'camera_vertical_turn_down_angle')
        view_wrapper(kwg, obj, float,
                        'camera_turn_look_ahead_slerp_amount')
        view_wrapper(kwg, obj, int, 'hit_pts')
        view_wrapper(kwg, obj, float, 'nudge_force')
        view_wrapper(kwg, obj, int, 'strafe_force')
        view_wrapper(kwg, obj, int, 'strafe_power_usage')
        view_wrapper(kwg, obj, float, 'explosion_resistance')
        view_wrapper(kwg, obj, int, 'ids_info1')
        view_wrapper(kwg, obj, int, 'ids_info2')
        view_wrapper(kwg, obj, int, 'ids_info3')
        view_wrapper(kwg, obj, int, 'ship_class')
        view_wrapper(kwg, obj, int, 'nanobot_limit')
        view_wrapper(kwg, obj, int, 'shield_battery_limit')

        if 'type' in obj.keys():
            kwg['typeof'] = str(obj['type'][0])

        if 'nickname' in obj.keys() and 'hp_type' in obj.keys():
            u.hp_type[obj['nickname'][0]] = obj['hp_type']

        try:
            dic = xmltodict.parse(u.infocards[kwg['ids_info']][1])[
                'RDL']['TEXT']
            if not dic[0]:
                dic = xmltodict.parse(u.infocards[kwg['ids_info1']][1])[
                    'RDL']['TEXT']
            kwg['info_name'] = dic[0]
        except KeyError:
            print("ERR not able to find infocard for ship object #", i,
            " ship nickname", kwg.get('nickname', 'no nickname'),
            "ids_info", kwg.get('ids_info', 'no ids_info')
            )
        except xmltodict.expat.ExpatError:
            print("ERR xmltodict.expat.ExpatError, can't parse infocard xml #", i,
            " ship nickname", kwg.get('nickname', 'no nickname'),
            "ids_info", kwg.get('ids_info', 'no ids_info')
            )

        if kwg['nickname'] in u.goods_by_ship['shiphull']:
            hull = u.goods_by_ship['shiphull'][kwg['nickname']
                                                ]['nickname'][0]
            try:
                ship = u.goods_by_hull['ship'][hull]
            except KeyError:
                print("ERR no package in goods.ini for ship hull =", hull)
            # print('123')
            # for add in ship['addon']
            # i f add[0] in u.equipment['misc_equip.ini']['[power]'].keys()
            # TODO find in addons powercore st_equip
            # and perhaps engine in engine_equip

        db_data = model_obj(
            **kwg
        )
        db_data.save()
        #except Exception as error:
            #print("ERR in filling ship #", i)


def recursive_reading(folderpath):
    """"Function to read all files from Universe folder resursively"""
    dictpath = {}
    for (dirpath, dirnames, filenames) in walk(folderpath):
        # 1 Level
        for filename in filenames:
            try:
                # dictpath[filename] = 1
                dictpath[filename] = parse_file(
                    os.path.join(dirpath, filename))
            except IndexError:
                print('ERR IndexError in ', filename)
            except UnicodeDecodeError:
                print("ERR UnicodeDecodeError in ", filename)

        for dirname in dirnames:
            dictpath[dirname] = recursive_reading(
                os.path.join(dirpath, dirname))

        break

    return dictpath


def folder_reading(folderpath):
    """Fuction to parse all files in one folder"""
    dictpath = {}
    for (__, __, filenames) in walk(folderpath):
        for filename in filenames:
            try:
                dictpath[filename] = parse_file(
                    os.path.join(folderpath, filename))
            except IndexError:
                print('ERR IndexError in ', filename)
            except UnicodeDecodeError:
                print("ERR UnicodeDecodeError in ", filename)
        break
    return dictpath


def split_goods(dic, key):
    """"Converts parsed data from list into being accessable by chosen hash key"""
    goods = u.equipment['goods.ini']['[good]']
    for obj in goods:
        if obj['category'][0] not in dic:
            dic[obj['category'][0]] = {}

        if key == 'shiphull':
            if key == obj['category'][0] and 'ship' in obj.keys():
                dic[obj['category'][0]][obj['ship'][0]] = obj
        elif key == 'ship':
            if key == obj['category'][0] and 'hull' in obj.keys():
                dic[obj['category'][0]][obj['hull'][0]] = obj
        else:
            if key in obj:
                dic[obj['category'][0]][obj[key][0]] = obj
        # print(f"ERR in goods_by_{key} #", i)


def on_server_start():
    "main function to load database on a server run"
    if os.environ.get('RUN_MAIN', None) == 'true':
        return

    if settings.DARK_COPY:
        clean_folder_from_files(settings.PATHS.dark_copy_dir)
    #     import stat
    #     if not os.access(settings.dark_copy_dir, os.W_OK):
    #         # Is the error an access error ?
    #         os.chmod(settings.dark_copy_dir, stat.S_IWUSR)
    #     else:
    #         pass
    #     os.remove(settings.dark_copy_dir)

    if settings.DARK_LOAD:
        management.call_command('flush', '--noinput')
        management.call_command('migrate')
        management.call_command('loaddata', 'dump.json')
        return

    if not settings.DARK_PARSE:
        return

    # import flint
    # flint.paths.set_install_path('Freelancer')
    # comms = flint.get_commodities()

    management.call_command('flush', '--noinput')
    management.call_command('migrate')

    from commodities.models import Commodity
    from ship.models import Ship

    u.equipment = folder_reading(settings.PATHS.equipment_dir)
    u.infocards = parse_infocards(settings.PATHS.infocards_path)
    u.universe = recursive_reading(settings.PATHS.universe_dir)
    u.ships = folder_reading(settings.PATHS.ships_dir)

    split_goods(u.goods_by_nickname, 'nickname')
    split_goods(u.goods_by_ship, 'shiphull')
    split_goods(u.goods_by_hull, 'ship')

    fill_commodity_table(Commodity)
    fill_ship_table(Ship)

    if settings.DARK_SAVE:
        management.call_command(
            'dumpdata',
            natural_foreign=True,
            natural_primary=True,
            indent=2,
            output="dump.json"
            )


class ParsingConfig(AppConfig):
    "main class to launch server in a different configuration to parse/save/load data"

    name = 'parsing'

    def ready(self):
        on_server_start()


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
