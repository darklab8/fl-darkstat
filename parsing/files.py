"module to work with files and folders, consisting modified module os operations"
import functools
import os
import codecs
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

        if settings.DARK_COPY:
            filename = args[0]
            targetname = filename.replace(
                settings.FREELANCER_FOLDER, settings.DARK_COPY_NAME)

            create_nested_folder(os.path.dirname(targetname))

            copy(
                filename,
                targetname
            )

        return func(*args, **kwargs)
    return wrapper_do_twice


@save_to_dark_copy
def read_regular_file(filename):
    "get content of regular file"
    with open(filename) as filelink:
        return filelink.readlines()


@save_to_dark_copy
def read_utf8_file(filename):
    "get content of utf8 encoded file"
    with codecs.open(filename, "r", "utf-8") as filelink:
        return filelink.readlines()
