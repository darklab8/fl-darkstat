
from types import MappingProxyType
from collections import namedtuple
from types import SimpleNamespace

# Changeable dict
dict_ = {
    "data": 456,
}
sdict = SimpleNamespace(**dict_)
print(sdict.data)
sdict.data = 123
sdict.new = 456
print(sdict.data)
print('sdict new = ' + str(sdict.new))
breakpoint()

# read only tuple from dict
dict_ = {
    "data": 456,
    "test": 123
}
Custom = namedtuple('custom', dict_)
ndict = Custom(**dict_)
ndict.data
print(ndict.data)

# read only dict
dict_ = {
    "data": 456,
}
rdict_ = MappingProxyType(dict_)
print("r = " + str(rdict_['data']))
rdict_['data'] = 123
