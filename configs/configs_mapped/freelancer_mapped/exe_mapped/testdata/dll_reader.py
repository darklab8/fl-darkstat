# type: ignore
"""
By Alex (of Discovery GC and The Starport), 31st May 2015 & 27th June 2016
Based heavily on PHP work done for the Discovery GC wiki by
Christopher Shake (cshake@gmail.com) in January 2010
"""

import os, struct, sys, io

def ReadText(fh, count):
    strout = b''
    for j in range(0, count):
        if j == 0:
            h = fh.read(2)
            if h == b"\xff\xfe":
                continue # strip BOM
            strout += h
        else:
            strout += fh.read(2)

    windows_decoded = strout.decode('windows-1252')
    windows_decoded[::2]
    return strout.decode('windows-1252')[::2].encode('utf-8')

def parseDLL(fh: io.BufferedReader, out, global_offset):
    # Header stuff, most of it is just read and ignored but we need a few addresses from it.
    fh.seek(0x3C)
    PE_sig_loc, = struct.unpack('B', fh.read(1)) # get PE Sig location
    fh.seek(PE_sig_loc + 4) # goto COFF header (after sig)
    COFF_Head_Machine, = struct.unpack('h', fh.read(2)) # 014c - i386 or compatible
    COFF_Head_NumberOfSections, = struct.unpack('h', fh.read(2))
    COFF_Head_TimeDateStamp, = struct.unpack('=l', fh.read(4))
    COFF_Head_PointerToSymbolTable, = struct.unpack('=l', fh.read(4))
    COFF_Head_NumberOfSymbols, = struct.unpack('=l', fh.read(4))

    data = fh.read(2)
    COFF_Head_SizeOfOptionalHeader, = struct.unpack('h', data)
    COFF_Head_Characteristics, = struct.unpack('h', fh.read(2)) # 210e
    OPT_Head_Start = fh.tell()
    if COFF_Head_SizeOfOptionalHeader != 0: # image header exists
        OPT_Head_Magic, = struct.unpack('h', fh.read(2))
        OPT_Head_MajorLinkerVers, = struct.unpack('c', fh.read(1))
        OPT_Head_MinorLinkerVers, = struct.unpack('c', fh.read(1))
        OPT_Head_SizeOfCode, = struct.unpack('=l', fh.read(4))
        OPT_Head_SizeOfInitializedData, = struct.unpack('=l', fh.read(4))
        OPT_Head_SizeOfUninitializedData, = struct.unpack('=l', fh.read(4))
        OPT_Head_AddressOfEntryPoint, = struct.unpack('=l', fh.read(4))
        OPT_Head_BaseOfCode, = struct.unpack('=l', fh.read(4))
        if OPT_Head_Magic == 0x20B: # if it's 64-bit
            OPT_Head_ImageBase, = struct.unpack('q', fh.read(8))
        else:
            OPT_Head_BaseOfData, = struct.unpack('=l', fh.read(4))
            OPT_Head_ImageBase, = struct.unpack('=l', fh.read(4))

        SectionAlignment = fh.read(4)
        FileAlignment = fh.read(4)
        MajorOperatingSystemVersion = fh.read(2)
        MinorOperatingSystemVersion = fh.read(2)
        MajorImageVersion = fh.read(2)
        MinorImageVersion = fh.read(2)
        MajorSubsystemVersion = fh.read(2)
        MinorSubsystemVersion = fh.read(2)
        Win32VersionValue = fh.read(4)
        SizeOfImage = fh.read(4)
        SizeOfHeaders = fh.read(4)
    # Get the section header info, we only care about ".rsrc" though
    fh.seek(OPT_Head_Start + COFF_Head_SizeOfOptionalHeader)
    DLL_Sections = {}
    for i in range(0, COFF_Head_NumberOfSections):
        nt = fh.read(8)
        name = nt.decode('utf-8').strip("\x00") # TODO: There was much more complex code for this in PHP, but the input format is completely different. Like different order and format different.
        DLL_Sections[name] = {}
        DLL_Sections[name]['VirtualSize'], = struct.unpack('=l', fh.read(4))
        DLL_Sections[name]['VirtualAddress'], = struct.unpack('=l', fh.read(4))
        DLL_Sections[name]['SizeOfRawData'], = struct.unpack('=l', fh.read(4))
        DLL_Sections[name]['PointerToRawData'], = struct.unpack('=l', fh.read(4))
        DLL_Sections[name]['PointerToRelocations'], = struct.unpack('=l', fh.read(4))
        DLL_Sections[name]['PointerToLinenumbers'], = struct.unpack('=l', fh.read(4))
        DLL_Sections[name]['NumberOfRelocations'], = struct.unpack('h', fh.read(2))
        DLL_Sections[name]['NumberOfLinenumbers'], = struct.unpack('h', fh.read(2))
        DLL_Sections[name]['Characteristics'], = struct.unpack('=l', fh.read(4))

    rsrcstart = DLL_Sections['.rsrc']['PointerToRawData']
    fh.seek(rsrcstart + 14) # go to start of .rsrc
    numentries, = struct.unpack('h', fh.read(2))
    datatypes = []
    # get the data types stored in the resource section
    for i in range(0, numentries):
        dataType, = struct.unpack('=l', fh.read(4))

        doi = fh.read(2)
        doj = fh.read(1)
        dataOffset, = struct.unpack('<i', doi + doj + '\x00'.encode('utf-8'))

        datatypes.append({'type': dataType, 'offset': dataOffset})
        fh.seek(1, os.SEEK_CUR) # jump ahead 1 byte

    # each different data type is stored in a block, loop through each
    for i in range(0, len(datatypes)):
        fh.seek(datatypes[i]['offset'] + rsrcstart)
        name = fh.read(8)
        fh.seek(6, os.SEEK_CUR)
        numentries, = struct.unpack('h', fh.read(2))
        sectionstart = fh.tell() # remember where we are here
        for entry in range(0, numentries):
            # get the id number and location of this entry
            idnum, = struct.unpack('i', fh.read(4))

            doi = fh.read(2)
            doj = fh.read(1)
            nameloc, = struct.unpack('<i', doi + doj + '\x00'.encode('utf-8'))

            brk = fh.read(1)
            backto = fh.tell() # remember where we were in the list of entries

            fh.seek(rsrcstart + nameloc) # jump to the entry
            name = fh.read(8) # get the name
            fh.seek(8, os.SEEK_CUR)
            lang = fh.read(4) # language for this resource

            someinfoloc, = struct.unpack('i', fh.read(4)) # location of the real location of the entry....


            fh.seek(rsrcstart + someinfoloc) # jump there
            absloc, = struct.unpack('i', fh.read(4)) # get the real location
            datalength, = struct.unpack('i', fh.read(4)) # entry length in bytes

            # now that we've got absolute location of each resource, get it!
            fh.seek(absloc)
            if datatypes[i]['type'] == 0x06: # string table
                for strindex in range(0, 16): # each string table has up to 16 entries
                    tableLen, = struct.unpack('h', fh.read(2))
                    if not tableLen:
                        continue # drop completely empty strings
                    ids_text = ReadText(fh, tableLen)
                    ids_index = (idnum - 1)*16 + strindex + global_offset
                    out[ids_index] = ids_text
            elif datatypes[i]['type'] == 0x17: # html
                ids_index = idnum + global_offset

                # py this script works great for this index
                # but not got script
                if 500904 == ids_index:
                    print(123)

                if datalength % 2:
                    datalength -= 1 # if odd length, ignore the last byte (UTF-16 is 2 bytes per character...)
                ids_text = ReadText(fh, datalength // 2).rstrip()
                out[ids_index] = ids_text

            # go back and get the next one
            fh.seek(backto)

    fh.close()

def parseDLLs(dll_fnames):
    out = {}
    for idx, name in enumerate(dll_fnames):
#        print("Now parsing DLL {} ({})".format(name, idx + 1))
        with open(name, 'rb') as fh:
            if fh:
                fh.tell
                global_offset = (2 ** 16) * (idx) # the ids_number for index 0 in this file
                parseDLL(fh, out, global_offset)
            else:
                sys.stderr.write("Could not open {}\n".format(name))

    return out

from pathlib import Path
def getAllInfocards():
    exe_path = Path("/home/naa/windows10shared/fl-files-discovery/EXE")
    # Number of element i simportant here apperently
    # as resources.dll is supposed to be always 0 indexed
    # and the rest should get their index number from freelancer.ini
    dllPaths = [
        str(exe_path/"resources.dll"),
        str(exe_path/"InfoCards.dll"),
        str(exe_path/"MiscText.dll"),
        # str(exe_path/"nameresources.dll"),
        # str(exe_path/"equipresources.dll"),
        # str(exe_path/"offerbriberesources.dll"),
        # str(exe_path/"misctextinfo2.dll"),
        # str(exe_path/"Discovery.dll"),
        # str(exe_path/"DsyAddition.dll"),
    ]
    ids = parseDLLs(dllPaths)
#    print("Done reading DLLs!")
    return ids

import json
import base64
if __name__=="__main__":
    ids:dict[str,bytes] = getAllInfocards()
    # pprint(ids)
    with open("parsed.json","w") as file:
        file.write(json.dumps(
            {key: base64.b64encode(value).decode('utf-8') for key, value in ids.items()}
        ))

    with open("parsed.json","r") as file:
        data =file.read()
        dicty: dict[str,str] = json.loads(data)

        for i, (key, value) in enumerate(dicty.items()):
            print(f"{i=} {key=}, value={base64.b64decode(value.encode('utf-8'))!r}")

            if i > 10:
                break
        
        # listed = list(dicty.keys())
        # print(listed[0])
        # print(listed[-1])

        # print(base64.b64decode(dicty["33389"]))

            
