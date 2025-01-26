// Given by Laz

public static uint CreateId(string nickname)
{
    if (StringUtils.createIDTable == null)
    {
    StringUtils.createIDTable = new uint[256];
    for (uint index1 = 0; index1 < 256U; ++index1)
    {
        uint num = index1;
        for (uint index2 = 0; index2 < 8U; ++index2)
        num = ((int) num & 1) == 1 ? num >> 1 ^ 671105024U : num >> 1;
        StringUtils.createIDTable[(int) index1] = num;
    }
    }
    byte[] bytes = Encoding.ASCII.GetBytes(nickname.ToLowerInvariant());
    uint num1 = 0;
    for (int index = 0; index < bytes.Length; ++index)
    num1 = num1 >> 8 ^ StringUtils.createIDTable[(int) (byte) num1 ^ (int) bytes[index]];
    return (uint) (((int) (num1 >> 24) | (int) (num1 >> 8) & 65280 | (int) num1 << 8 & 16711680 | (int) num1 << 24) >>> 2 | int.MinValue);
}

public static uint CreateFactionId(string nickname)
{
    if (StringUtils.createFactionIDTable == null)
    {
    StringUtils.createFactionIDTable = new uint[256];
    for (uint index1 = 0; index1 < 256U; ++index1)
    {
        uint num = index1 << 8;
        for (uint index2 = 0; index2 < 8U; ++index2)
        num = (((int) num & 32768) == 32768 ? (uint) ((int) num << 1 ^ 4129) : num << 1) & (uint) ushort.MaxValue;
        StringUtils.createFactionIDTable[(int) index1] = num;
    }
    }
    byte[] bytes = Encoding.ASCII.GetBytes(nickname.ToLowerInvariant());
    uint factionId = (uint) ushort.MaxValue;
    for (uint index = 0; (long) index < (long) bytes.Length; ++index)
    factionId = (factionId & 65280U) >> 8 ^ StringUtils.createFactionIDTable[(int) factionId & (int) byte.MaxValue ^ (int) bytes[(int) index]];
    return factionId;
}
