"module to work with files and folders, consisting modified module os operations"
import functools
import os
import re
from os import walk
import codecs
import xmltodict
from shutil import copy2 as copy
from django.conf import settings


def clean_folder_from_files(dir_to_search):
    "delete all files in folde recursively"
    for dirpath, __, filenames in os.walk(dir_to_search):

        for filename in filenames:
            try:
                os.remove(os.path.join(dirpath, filename))
            except OSError as ex:
                print(ex)


def create_nested_folder(folderpath):
    "if our folder path is missing any folders, create them"
    try:
        os.makedirs(folderpath)
    except FileExistsError:
        # directory already exists
        pass


def save_to_dark_copy(func):
    "sneaky file copy to dark_copy folder for workflow unit tests"
    @functools.wraps(func)
    def wrapper_do_twice(*args, **kwargs):

        # active = args[1]
        # if active:
        filename = args[0]
        targetname = filename.replace(
            settings.PATHS.freelancer_folder, settings.PATHS.dark_copy_name)

        create_nested_folder(os.path.dirname(targetname))

        copy(
            filename,
            targetname
        )

        return func(*args, **kwargs)
    return wrapper_do_twice

#@save_to_dark_copy
def read_regular_file(filename):
    "get content of regular file"
    with open(filename) as filelink:
        return filelink.readlines()


#@save_to_dark_copy
def read_utf8_file(filename):
    "get content of utf8 encoded file"
    with codecs.open(filename, "r", "utf-8") as filelink:
        return filelink.readlines()


def recursive_reading(folderpath):
    """"Function to read all files from Universe folder resursively"""
    dictpath = {}
    for (dirpath, dirnames, filenames) in walk(folderpath):
        # 1 Level
        for filename in filenames:
            try:
                # dictpath[filename] = 1
                dictpath[filename] = parse_file(os.path.join(dirpath, filename))
            except IndexError:
                print("ERR IndexError in ", filename)
            except UnicodeDecodeError:
                print("ERR UnicodeDecodeError in ", filename)

        for dirname in dirnames:
            dictpath[dirname] = recursive_reading(os.path.join(dirpath, dirname))

        break

    return dictpath


def folder_reading(folderpath):
    """Fuction to parse all files in one folder"""
    dictpath = {}
    for (__, __, filenames) in walk(folderpath):
        for filename in filenames:
            try:
                dictpath[filename] = parse_file(os.path.join(folderpath, filename))
            except IndexError:
                print("ERR IndexError in ", filename)
            except UnicodeDecodeError:
                print("ERR UnicodeDecodeError in ", filename)
        break
    return dictpath


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
                strip_from_rn(content[i + 1]),
                strip_from_rn(content[i + 2]),
            ]

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
            while (re.search(regex_for_headers, content[i + 1]) is None) and (
                i < (line_count - 2)
            ):

                if content[i] == "\n":
                    i += 1
                    continue

                if re.search("^(;)", content[i]) is not None:
                    i += 1
                    continue

                splitted = content[i].replace(" ", "").replace("\n", "").split("=")

                if len(splitted) == 2:
                    if splitted[0] not in obj.keys():
                        obj[splitted[0]] = []

                    if "," not in splitted[1]:
                        obj[splitted[0]].append(splitted[1])
                    else:
                        obj[splitted[0]].append(splitted[1].split(","))
                else:
                    print("ERR splitted")
                i += 1

            output[header].append(obj)

        i += 1
    return output


def split_goods(goods, dic, key):
    """"Converts parsed data from list into being accessable by chosen hash key"""
    for obj in goods:
        if obj["category"][0] not in dic:
            dic[obj["category"][0]] = {}

        if key == "shiphull":
            if key == obj["category"][0] and "ship" in obj.keys():
                dic[obj["category"][0]][obj["ship"][0]] = obj
        elif key == "ship":
            if key == obj["category"][0] and "hull" in obj.keys():
                dic[obj["category"][0]][obj["hull"][0]] = obj
        else:
            if key in obj:
                dic[obj["category"][0]][obj[key][0]] = obj
        # print(f"ERR in goods_by_{key} #", i)


def rearrange_array_to_dict_by_keys(arr, key):
    """"Converts parsed data from list into being accessable by chosen hash key"""
    dic = {}
    for elem in arr:
        if key in elem:
            dic[elem[key][0]] = elem
    return dic

class Dicts:
    """ "This class will have parsed data from files
    stored in preparation to be sent to database"""

    equipment = {}
    universe = {}
    infocards = {}
    ships = {}
    hp_type = {}

    goods_by_nickname = {}
    goods_by_ship = {}
    goods_by_hull = {}

    misc_equip_power_by_nickname = {}


def main_parse():
    dicty = Dicts()
    dicty.equipment = folder_reading(settings.PATHS.equipment_dir)
    dicty.infocards = parse_infocards(settings.PATHS.infocards_path)
    dicty.universe = recursive_reading(settings.PATHS.universe_dir)
    dicty.ships = folder_reading(settings.PATHS.ships_dir)

    goods = dicty.equipment["goods.ini"]["[good]"]
    split_goods(goods, dicty.goods_by_nickname, "nickname")
    split_goods(goods, dicty.goods_by_ship, "shiphull")
    split_goods(goods, dicty.goods_by_hull, "ship")

    dicty.misc_equip_power_by_nickname = rearrange_array_to_dict_by_keys \
        (dicty.equipment['misc_equip.ini']['[power]'], "nickname")
    return dicty