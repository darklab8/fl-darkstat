import BufferView, { BufferObject } from "../utility/bufferview";

/**
 * Localized text strings and info cards (a weird rich text format stored as XML) are stored as STRING and HTML
 * resources in DLL files.
 * 
 * Editor supports only reading these resources by parsing relevant parts of DLL file to get to data chunks.
 * Writing DLL resources is outside the scope of the editor.
 * If someone wants to create functionality to write resources they're more than welcome to try.
 * Generally it would involve nodejs, node-ffi-napi and using UpdateResource function from kernel32.dll.
 */

/** Each DLL can contain up to 65535 resource entries. */
const MAX_RESOURCES = 0xFFFF;
const OFFSET_MASK = 0x7FFFFFFF;

const HEADER_SIGNATURE = 0x5A4D; // MZ
const HEADER_PORTABLE_SIZE = 0x18;
const HEADER_PORTABLE_SIGNATURE = 0x00004550;
const SECTION_SIZE = 0x28;
const SECTION_RESOURCES = ".rsrc";
const STRING_BLOCK_SIZE = 16;

enum ResourceType {
    STRING = 6,
    HTML = 23,
}

class PortableHeader implements BufferObject {
    readonly byteLength = HEADER_PORTABLE_SIZE;

    public signature = 0;
    public machine = 0;
    public numberOfSections = 0;
    public timeDateStamp = 0;
    public pointerToSymbolsTable = 0;
    public numberOfSymbols = 0;
    public sizeOfOptionalHeader = 0;
    public characteristics = 0;

    readBuffer(view: BufferView) {
        this.signature = view.readUint32();
        this.machine = view.readUint16();
        this.numberOfSections = view.readUint16();
        this.timeDateStamp = view.readUint32();
        this.pointerToSymbolsTable = view.readUint32();
        this.numberOfSymbols = view.readUint32();
        this.sizeOfOptionalHeader = view.readUint16();
        this.characteristics = view.readUint16();
    }

    writeBuffer() {
        throw new Error("Method not implemented.");
    }
}

class Section implements BufferObject {
    readonly byteLength = SECTION_SIZE;

    public name = "";
    public virtualSize = 0;
    public virtualAddress = 0;
    public sizeOfRawData = 0;
    public pointerToRawData = 0;
    public pointerToRelocations = 0;
    public pointerToLineNumbers = 0;
    public numberOfRelocations = 0;
    public numberOfLineNumbers = 0;
    public characteristics = 0;

    readBuffer(view: BufferView) {
        this.name = view.readString(8).replaceAll("\0", "").trim();
        this.virtualSize = view.readUint32();
        this.virtualAddress = view.readUint32();
        this.sizeOfRawData = view.readUint32();
        this.pointerToRawData = view.readUint32();
        this.pointerToRelocations = view.readUint32();
        this.pointerToLineNumbers = view.readUint32();
        this.numberOfRelocations = view.readUint16();
        this.numberOfLineNumbers = view.readUint16();
        this.characteristics = view.readUint32();
    }

    writeBuffer() {
        throw new Error("Method not implemented.");
    }
}

type DirectoryEntry = {
    name: number,
    offsetToData: number,
}

class Directory implements BufferObject {

    readonly byteLength = 0;

    public namedEntries: DirectoryEntry[] = [];
    public indexEntries: DirectoryEntry[] = [];

    public characteristics = 0;
    public timeDateStamp = 0;
    public majorVersion = 0;
    public minorVersion = 0;
    public numberOfNamedEntries = 0;
    public numberOfIdEntries = 0;

    readBuffer(view: BufferView) {
        this.characteristics = view.readUint32();
        this.timeDateStamp = view.readUint32();
        this.majorVersion = view.readUint16();
        this.minorVersion = view.readUint16();
        this.numberOfNamedEntries = view.readUint16();
        this.numberOfIdEntries = view.readUint16();
    
        for (let i = 0; i < this.numberOfNamedEntries; i++)
            this.namedEntries.push({ name: view.readUint32(), offsetToData: view.readUint32() });
        
        for (let i = 0; i < this.numberOfIdEntries; i++)
            this.indexEntries.push({ name: view.readUint32(), offsetToData: view.readUint32() });
    }

    writeBuffer() {
        throw new Error("Method not implemented.");
    }
}

type StringEntry = {
    name: number,
    offsetToString: number,
    size: number,
    codePage: number,
}

type LanguageData = {
    offsetToData: number,
    size: number,
    codePage: number,
    reserved: number,
}

/**
 * Read data languages.
 * @param view
 * @param offset
 */
function readLanguages(view: BufferView) {
    const directory = new Directory();
    directory.readBuffer(view);

    const result: LanguageData[] = [];

    for (const { offsetToData: offsetToLanguageData } of directory.indexEntries) {
        view.offset = offsetToLanguageData & OFFSET_MASK;

        result.push({
            offsetToData: view.readUint32(),
            size: view.readUint32(),
            codePage: view.readUint32(),
            reserved: view.readUint32(),
        });
    }

    return result;
}

/**
 * Read string entries from directory at offset.
 * Directory contains language blocks which in turn point to actual string data.
 * Only first codepage block is used.
 * @param view
 * @param offset
 */
function readStrings(view: BufferView) {
    const directory = new Directory();
    directory.readBuffer(view);

    const strings: StringEntry[] = [];

    let offsetToString, size, codePage;

    for (const { name, offsetToData } of directory.indexEntries) {
        view.offset = offsetToData & OFFSET_MASK;

        // Read first language entry.
        // TODO: Filter by codepage which should be read from editor configuration file.
        [ { offsetToData: offsetToString = 0, size = 0, codePage = 0 } ] = readLanguages(view);
        if (! size) continue;
        
        strings.push({ name, offsetToString, size, codePage });
    }

    return strings;
}

export class InfoResource {

    public filename = "";

    constructor(
        readonly strings: string[] = [],
        readonly cards: string[] = [],
    ) {}

    async loadFile(file: File) {
        this.cards.length = this.strings.length = 0;
        this.filename = file.name;

        // Load whole file in a go to parse it quickly and slow reads. 
        const view = await BufferView.fromBlob(file);
        if (view.readUint16() !== HEADER_SIGNATURE) throw new TypeError("Invalid DLL header signature.");

        view.offset = view.getUint32(0x3C, true);

        const header = new PortableHeader();
        header.readBuffer(view);

        if (header.signature !== HEADER_PORTABLE_SIGNATURE) throw new TypeError("Invalid DLL portable executable header signature.");

        const { numberOfSections, sizeOfOptionalHeader } = header;

        // Skip optional header.
        view.shift(sizeOfOptionalHeader);

        let section: Section, directory: Directory, data: BufferView, decoder = new TextDecoder("utf-16");

        for (let i = 0; i < numberOfSections; i++) {
            (section = new Section()).readBuffer(view);

            const { name, pointerToRawData, sizeOfRawData } = section;
            if (name !== SECTION_RESOURCES) continue;

            // Get section data block view.
            data = new BufferView(view.buffer, pointerToRawData & OFFSET_MASK, sizeOfRawData);

            (directory = new Directory()).readBuffer(data);

            // List resource types.
            for (const { name, offsetToData } of directory.indexEntries) {
                data.offset = offsetToData & OFFSET_MASK;

                switch (name) {
                    case ResourceType.STRING:
                        for (const { name, offsetToString, size } of readStrings(data)) {
                            const stringView = new BufferView(view.buffer, offsetToString, size);
                
                            // Read block of 16 strings.
                            for (let i = 0, id = (name - 1) * STRING_BLOCK_SIZE; i < STRING_BLOCK_SIZE; i++) {
                                const length = stringView.readUint16() * 2;
                                if (length > 0) this.strings[id + i] = stringView.readString(length, decoder);
                            }
                        }

                        break;
                    case ResourceType.HTML:
                        for (const { name, offsetToString, size } of readStrings(data))
                            this.cards[name] = decoder.decode(new Uint8Array(view.buffer, offsetToString, size));

                        break;
                }
            }
        }

        return true;
    }

    async saveFile() {
        throw new Error("Method not implemented.");
    }
}

export class InfoDatabase {

    protected resources: InfoResource[] = [];

    clear() {
        this.resources.length = 0;
    }

    async load(file: File) {
        const resource = new InfoResource();
        this.resources.push(resource);

        await resource.loadFile(file); 
    }

    async saveFile() {
        throw new Error("Method not implemented.");
    }

    /**
     * Get info string (ids_name).
     * @param id
     */
    getString(id: number = 0) {
        const volume = id / MAX_RESOURCES >> 0, index = (id % MAX_RESOURCES) - volume;
        return this.resources[volume]?.strings[index];
    }

    /**
     * Get info card raw string (ids_info).
     * @param id
     */
    getCard(id: number = 0) {
        const volume = id / MAX_RESOURCES >> 0, index = (id % MAX_RESOURCES) - volume;
        return this.resources[volume]?.cards[index];
    }
    
    /**
     * Get info card XML document.
     * @param id 
     * @returns 
     */
    getCardXML(id: number) {
        const source = this.getCard(id);
        if (! source) return;

        const document = new DOMParser().parseFromString(source, "application/xml");
        
        if (document.documentElement.nodeName.toLowerCase() === "parsererror") {
            throw new SyntaxError(document.documentElement.textContent ?? "General error");
        }

        return document;
    }
}