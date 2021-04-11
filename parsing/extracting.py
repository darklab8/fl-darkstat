import xmltodict


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
                elif "." in obj[name][0]:
                    kwg[name] = int(float(obj[name][0]))
                else:
                    raise ValueError from value_error_1
            else:
                raise ValueError from value_error_1


def view_wrapper_with_infocard(dicty, kwg, obj, data_type, name, infoname):
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
        if kwg[name] in dicty.infocards:
            kwg[infoname] = dicty.infocards[kwg[name]][1]
