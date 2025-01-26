# Provided by Alex from Freelancer Discovery

createIDTable = None
def FLFactionHash(nickName):
        FLFACTIONHASH_POLYNOMIAL = 0x1021

        # Build the crc lookup table if it hasn't been created
        global createIDTable
        if createIDTable is None:
                createIDTable = []
                for i in range(0, 2**8):
                        x = i << 8
                        for bit in range(0, 8):
                                x <<= 1
                                if x & 0x10000:
                                        x = (x ^ FLFACTIONHASH_POLYNOMIAL) & 0xFFFF
                        createIDTable.append(x)

        # Calculate the hash.
        hash = 0xFFFF
        for c in bytearray(nickName.lower().encode('utf-8')):
                hash = (hash >> 8) ^ createIDTable[(hash & 0xFF) ^ c]
        return hash

if __name__ == "__main__":
        import sys
        for n in sys.argv[1:]:
                h = FLFactionHash(n)
                print(n, h, '0x{:04x}'.format(h))