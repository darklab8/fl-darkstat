from django.apps import AppConfig
from django.conf import settings
import re
from os import walk
import os
from django.core import management

equipment = {}
universe = {}
infocards = {}
ships = {}
hp_type = {}

def strip_from_rn(a):
    return a.replace("\r","").replace("\n","")

def parse_infocards(filename):
    import codecs
    f = codecs.open( filename, "r", "utf-8" )
    #f = open(filename)

    content = f.readlines()
    
    regex_numbers = "^\d+\r|^\d+\n"
    output = {}
    line_count = len(content)
    #breakpoint()
    for i in range(line_count):
        #print(content[i])
        if (re.search(regex_numbers, content[i]) is not None):
            try:
                output[int(strip_from_rn(content[i]))] = [strip_from_rn(content[i+1]), strip_from_rn(content[i+2])]
            except:
                print('ERR in infocards parser')
        
    i+=1
    return output

def parse_file(filename):
    #breakpoint()
    f = open(filename)
    content = f.readlines()
    
    output = {}
    regex_for_headers = "(\[)\w+(\])"

    line_count = len(content)
    for i in range(line_count):
        if (re.search(regex_for_headers, content[i]) is not None):
            header = content[i].lower().replace("\n","")

            if header not in output.keys():
                output[header] = []

            i+=1
            obj = {}
            while (
                (re.search(regex_for_headers, content[i+1]) is None) and 
                (i < (line_count - 2))
                ):

                if (content[i] == '\n'):
                    i+=1
                    continue
                
                if (re.search("^(;)", content[i]) is not None):
                    i+=1
                    continue

                splitted = content[i].replace(" ","").replace("\n","").split('=')

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

        i+=1
    return output

def view_wrapper(kwg, obj, cl, name):
    if name in obj.keys():
        kwg[name] = cl(obj[name][0])

def view_wrapper_with_infocard(kwg, obj, cl, name, infoname):
    if name in obj.keys():
        kwg[name] = cl(obj[name][0])
        if kwg[name] in infocards:
            kwg[infoname] = (infocards[kwg[name]][1])

def fill_commodity_table(Commodity):
    #COMMODITY TABLE
    goods = equipment['select_equip.ini']
    arr = goods['[commodity]'].copy()
    for obj in arr:
        try:
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

            c = Commodity(
                **kwg
            )
            c.save()
        except:
            print("ERR in filling commodities", obj)

def fill_ship_table(Ship):
    #COMMODITY TABLE
    goods = ships['shiparch.ini']
    arr = goods['[ship]'].copy()
    for obj in arr:
        try:
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
            view_wrapper(kwg, obj, float, 'camera_turn_look_ahead_slerp_amount')
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

            global hp_type
            if 'nickname' in obj.keys() and 'hp_type' in obj.keys():
                hp_type[obj['nickname'][0]] = obj['hp_type']

            c = Ship(
                **kwg
            )
            c.save()
        except:
            print("ERR in filling ships", obj)

def RecursiveReading(folderpath):
    
    dictpath = {}
    for (dirpath, dirnames, filenames) in walk(folderpath):
        #1 Level
        for filename in filenames:
            try:
                #dictpath[filename] = 1
                dictpath[filename] = parse_file(os.path.join(dirpath,filename))
            except:
                print('ERROR in ', filename)

        for dirname in dirnames:
            dictpath[dirname] = RecursiveReading(os.path.join(dirpath,dirname))

        break
    
    return dictpath

def folder_reading(folderpath):
    dictpath = {}
    for (dirpath, dirnames, filenames) in walk(folderpath):
        for filename in filenames:
            try:
                dictpath[filename] = parse_file(os.path.join(folderpath,filename))
            except:
                print('ERROR in ', filename)
        break
    return dictpath

class MainConfig(AppConfig):
    name = 'main'
    def ready(self):
        if os.environ.get('RUN_MAIN', None) == 'true':
            return

        management.call_command('flush', '--noinput')
        management.call_command('migrate')

        #import flint
        #flint.paths.set_install_path('Freelancer')
        #comms= flint.get_commodities()
        from commodities.models import Commodity
        from ship.models import Ship

        global equipment
        equipment = folder_reading(settings.EQUIPMENT_DIR)

        global infocards
        infocards = parse_infocards(settings.INFOCARDS_PATH)

        fill_commodity_table(Commodity)

        global universe
        universe = RecursiveReading(settings.UNIVERSE_DIR)

        global ships
        ships = folder_reading(settings.SHIPS_DIR)

        fill_ship_table(Ship)
        #breakpoint()123

        # for filename in equipment.keys():
        #     for header in equipment[filename].keys():
        #         for obj in equipment[filename][header]:
        #             #print(obj)
        #             #breakpoint()
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
#     #if 'category' in goods[key].keys() and goods[key]['category'] == 'equipment':
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
#             if ('addon' not in subkey):
#                 print('ERR not addon')

                    

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