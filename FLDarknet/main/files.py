from django.conf import settings
from shutil import copy2 as copy
import functools
import os
import pathlib

def clean_folder_from_files(dir_to_search):
    for dirpath, dirnames, filenames in os.walk(dir_to_search):

        for filename in filenames:
            try:
                os.remove(os.path.join(dirpath,filename))
            except OSError as ex:
                print(ex)

def create_nested_folder(folderpath):
    try:
        os.makedirs(folderpath)
    except FileExistsError:
        # directory already exists
        pass

def save_to_dark_copy(func):
    @functools.wraps(func)
    def wrapper_do_twice(*args, **kwargs):

        if settings.DARK_COPY:
            filename = args[0]
            targetname = filename.replace(settings.FREELANCER_FOLDER,settings.DARK_COPY_NAME)

            create_nested_folder(os.path.dirname(targetname))

            copy(
                filename, 
                targetname
                )

        return func(*args, **kwargs)
    return wrapper_do_twice

@save_to_dark_copy
def read_regular_file(filename):
    f = open(filename)
    content = f.readlines()
    return content

@save_to_dark_copy
def read_utf8_file(filename):
    import codecs
    f = codecs.open( filename, "r", "utf-8" )
    content = f.readlines()
    return content