"""module to work with files and folders,
consisting modified module os operations"""
import functools
import os
import re
from os import walk
import codecs
from shutil import copy2 as copy
from django.conf import settings

from types import MappingProxyType
# from collections import namedtuple
from types import SimpleNamespace


def clean_folder_from_files(dir_to_search: str) -> None:
    "delete all files in folde recursively"
    for dirpath, __, filenames in os.walk(dir_to_search):

        for filename in filenames:
            try:
                os.remove(os.path.join(dirpath, filename))
            except OSError as ex:
                print(ex)


def create_nested_folder(folderpath: str) -> None:
    "if our folder path is missing any folders, create them"
    try:
        os.makedirs(folderpath)
    except FileExistsError:
        pass


def save_to_dark_copy(func: callable) -> callable:
    "sneaky file copy to dark_copy folder for workflow unit tests"

    @functools.wraps(func)
    def wrapper_do_twice(*args, **kwargs):

        filename = args[0]
        targetname = filename.replace(settings.PATHS.freelancer_folder,
                                      settings.PATHS.dark_copy_name)

        create_nested_folder(os.path.dirname(targetname))

        copy(filename, targetname)

        return func(*args, **kwargs)

    return wrapper_do_twice


def recursive_reading(folderpath: str) -> SimpleNamespace:
    """"Function to read all files from Universe folder resursively"""
    def prepapre_simple_name(filename: str) -> str:
        return filename.replace(".ini", "").lower()

    dictpath = SimpleNamespace()
    for (dirpath, dirnames, filenames) in walk(folderpath):
        # 1 Level
        for filename in filenames:
            try:
                # dictpath[filename] = 1
                setattr(dictpath, prepapre_simple_name(filename),
                        parse_file(os.path.join(dirpath, filename)))
            except IndexError:
                print("ERR IndexError in ", filename)
            except UnicodeDecodeError:
                print("ERR UnicodeDecodeError in ", filename)

        for dirname in dirnames:
            setattr(dictpath, prepapre_simple_name(dirname),
                    recursive_reading(os.path.join(dirpath, dirname)))

        break

    return dictpath


def parse_infocards(filename: str) -> MappingProxyType:
    """"Parses infocard file into dictionary"""
    def strip_from_rn(line: str) -> str:
        """Strips string from \r or \n trash"""
        return line.replace("\r", "").replace("\n", "")

    def read_utf8_file(filename: str) -> list:
        "get content of utf8 encoded file"
        with codecs.open(filename, "r", "utf-8") as filelink:
            return filelink.readlines()

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
    return MappingProxyType(output)


def parse_file(filename: str) -> SimpleNamespace:
    """Parses file into read only dictionary"""
    def read_regular_file(filename: str) -> list:
        "get content of regular file"
        with open(filename) as filelink:
            return filelink.readlines()

    content = read_regular_file(filename)

    output = SimpleNamespace()
    regex_for_headers = r"(\[)\w+(\])"

    line_count = len(content)
    for i in range(line_count):
        if re.search(regex_for_headers, content[i]) is not None:
            header = content[i].lower().replace("\n", "").replace(";", "")\
                .replace("[", "").replace("]", "")

            if not hasattr(output, header):
                setattr(output, header, [])

            i += 1
            obj = {}
            while (re.search(regex_for_headers, content[i + 1]) is
                   None) and (i < (line_count - 2)):

                if content[i] == "\n":
                    i += 1
                    continue

                if re.search("^(;)", content[i]) is not None:
                    i += 1
                    continue

                splitted = content[i].replace(" ", "").replace("\n",
                                                               "").split("=")

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

            getattr(output, header).append(obj)

        i += 1
    return output


def main_parse() -> SimpleNamespace:
    parsed = SimpleNamespace()
    parsed.equipment = recursive_reading(settings.PATHS.equipment_dir)
    parsed.infocards = parse_infocards(settings.PATHS.infocards_path)
    parsed.universe = recursive_reading(settings.PATHS.universe_dir)
    parsed.ships = recursive_reading(settings.PATHS.ships_dir)

    return parsed
