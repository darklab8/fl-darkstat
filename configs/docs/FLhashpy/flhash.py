# Provided by Alex from Freelancer Discovery

createIDTable = None
def FLHash(nickName):
        FLHASH_POLYNOMIAL = 0xA001
        LOGICAL_BITS = 30
        PHYSICAL_BITS = 32

        # Build the crc lookup table if it hasn't been created
        global createIDTable
        if createIDTable is None:
                createIDTable = []
                for i in range(0, 2**8):
                        x = i
                        for bit in range(0, 8):
                                if (x & 1) == 1:
                                        x = (x >> 1) ^ (FLHASH_POLYNOMIAL << (LOGICAL_BITS - 16))
                                else:
                                        x >>= 1
                        createIDTable.append(x)
                        #print("createIDTable " + str(i) + ": " + str(x))

        # Calculate the hash.
        hash = 0
        for c in bytearray(nickName.lower().encode('utf-8')):
                hash = (hash >> 8) ^ createIDTable[(hash & 0xFF) ^ c]
        # b0rken because byte swapping is not the same as bit reversing, but
        # that's just the way it is; two hash bits are shifted out and lost
        hash = (hash >> 24) | ((hash >> 8) & 0x0000FF00) | ((hash << 8) & 0x00FF0000) | ((hash & 0xFF) << 24)
        hash = (hash >> (PHYSICAL_BITS - LOGICAL_BITS)) | 0x80000000
        return hash

if __name__ == "__main__":
        import sys
        for n in sys.argv[1:]:
                h = FLHash(n)
                print(n, h, '0x{:08x}'.format(h))