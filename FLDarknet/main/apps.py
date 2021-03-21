from django.apps import AppConfig
from django.conf import settings
import re

data = {}
infocards = {}

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

class MainConfig(AppConfig):
    name = 'main'
    def ready(self):
        #import flint
        #flint.paths.set_install_path('Freelancer')
        #comms= flint.get_commodities()
        from commodities.models import Commodity
        Commodity.objects.all().delete()
        
        global data
        #breakpoint()
        from os import walk
        import os
        for (dirpath, dirnames, filenames) in walk(settings.EQUIPMENT_DIR):
            for filename in filenames:
                try:
                    data[filename] = parse_file(os.path.join(settings.EQUIPMENT_DIR,filename))
                except:
                    print('ERROR in ', filename)
            break

        global infocards
        infocards = parse_infocards(settings.INFOCARDS_PATH)

        goods = data['select_equip.ini']
        arr = goods['[commodity]'].copy()
        for obj in arr:
            try:
                if 'ids_name' in obj.keys():
                    obj['ids_name'] = int(obj['ids_name'][0])
                    if obj['ids_name'] in infocards:
                        obj['name'] = (infocards[obj['ids_name']][1])

                if 'ids_info' in obj.keys():
                    obj['ids_info'] = int(obj['ids_info'][0])
                    if obj['ids_info'] in infocards:
                        obj['infocard'] = (infocards[obj['ids_info']][1])

                if 'units_per_container' in obj.keys():
                    obj['units_per_container'] = int(obj['units_per_container'][0])

                if 'decay_per_second' in obj.keys():
                    obj['decay_per_second'] = int(obj['decay_per_second'][0])

                if 'hit_pts' in obj.keys():
                    obj['hit_pts'] = int(obj['hit_pts'][0])

                if 'pod_appearance' in obj.keys():
                    obj['pod_appearance'] = str(obj['pod_appearance'][0])

                if 'loot_appearance' in obj.keys():
                    obj['loot_appearance'] = str(obj['loot_appearance'][0])

                if 'nickname' in obj.keys():
                    obj['nickname'] = str(obj['nickname'][0])

                if 'volume' in obj.keys():
                    obj['volume'] = float(obj['volume'][0])

                c = Commodity(
                    **obj
                )
                c.save()
            except:
                print("ERR in filling commodities", obj)

        # for filename in data.keys():
        #     for header in data[filename].keys():
        #         for obj in data[filename][header]:
        #             #print(obj)
        #             #breakpoint()
        #             break
            
        
        # goods = parse_file(settings.GOODS_DIR)
        # select_equip = parse_file(settings.SEL_EQUIP_DIR)
        # market_commodities = parse_file(settings.MARKET_DIR)

# test = set()
# goods = data['select_equip.ini']
# arr = goods['[commodity]']
# for obj in arr:
#     for key in obj.keys():
#         if key not in test:
#             print(key, " = ", obj[key])
#             test.add(key)

# goods = data['select_equip.ini']
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