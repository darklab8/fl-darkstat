
export interface BufferObject {

    /* Object will tell us array length to allocate or reserve ahead. */
    readonly byteLength: number;

    /** Object reads its state from buffer. */
    readBuffer(view: BufferView): void;

    /** Objects writes its state into buffer. */
    writeBuffer(view: BufferView): void;
}

export enum DataType { UINT8, INT8, UINT16, INT16, UINT32, INT32, BIGINT64, BIGUINT64, FLOAT32, FLOAT64 };

/**
 * A much more usable version of original DataView with internal offset needle and endianness setting.
 */
export default class BufferView extends DataView {
    constructor(
        buffer: ArrayBufferLike,
        byteOffset?: number,
        byteLength?: number,

        /** Multibyte endianness. */
        public littleEndian = true,

        /** Internal offset. */
        public offset = 0,
    ) {
        super(buffer, byteOffset, byteLength);
    }

    /**
     * Create BufferView from BufferObject(s).
     * @param objects
     * @returns 
     */
    static fromObjects(...objects: BufferObject[]) {
        const view = this.allocate(...objects);
        view.writeObjects(...objects);
        return view;
    }

    /**
     * Create BufferView from another ArrayBufferView.
     * @param value 
     * @returns 
     */
    static fromArray(value: ArrayBufferView) {
        return new this(value.buffer, value.byteOffset, value.byteLength);
    }

    /**
     * Allocate new buffer for BufferObject(s).
     * @param objects 
     * @returns 
     */
    static allocate(...objects: BufferObject[]) {
        return new this(new ArrayBuffer(objects.reduce((length, target) => length + target.byteLength, 0)));
    }

    /**
     * Create BufferView asynchronously from Blob.
     * @param value
     * @returns 
     */
    static async fromBlob(value: Blob) {
        return new this(await value.arrayBuffer());
    }

    /** Bytes left to end of buffer. */
    get byteRemain() {
        return this.byteLength - this.offset;
    }

    /**
     * Find sequence of bytes from offset.
     * @param value
     * @param fromOffset Starting offset
     * @returns Offset to matching sequence of bytes.
     */
    indexOf(value: ArrayBufferView, fromOffset = this.offset) {

        // Ideallly get optimal length view based on byteOffset (must be aligned) and haystack byteLength.
        const haystack = new Uint8Array(this.buffer, this.byteOffset, this.byteLength);
        const needle = new Uint8Array(value.buffer, value.byteOffset, value.byteLength);

        let offset = fromOffset, index = 0;

        // Compare bytes from current offset.
        while (haystack.byteLength - offset >= needle.byteLength) {
            for (index = 0; index < needle.length; index++)
                if (haystack[offset + index] !== needle[index])
                    break; // Values don't match.

            // Match found. 
            if (index >= needle.length) return offset;

            offset++;
        }

        return -1;
    }

    /**
     * Slice view from current offset.
     * @param length 
     * @returns 
     */
    slice(length = this.byteRemain) {
        return new BufferView(this.buffer, this.byteOffset + this.shift(length), length);
    }

    /**
     * Shifts internal offset by amount of bytes.
     * @param value Amount of bytes to shift by.
     * @returns Previous offset value.
     */
    shift(value: number) {
        return (this.offset += value) - value;
    }

    /** Read unsigned 8-bit integer number. */
    readUint8() {
        return this.getUint8(this.shift(Uint8Array.BYTES_PER_ELEMENT));
    }

    /** Write unsigned 8-bit integer number. */
    writeUint8(value: number) {
        return this.setUint8(this.shift(Uint8Array.BYTES_PER_ELEMENT), value);
    }

    /** Read signed 8-bit integer number. */
    readInt8() {
        return this.getInt8(this.shift(Int8Array.BYTES_PER_ELEMENT));
    }

    /** Write signed 8-bit integer number. */
    writeInt8(value: number) {
        return this.setInt8(this.shift(Int8Array.BYTES_PER_ELEMENT), value);
    }

    /** Read unsigned 16-bit integer number. */
    readUint16(littleEndian = this.littleEndian) {
        return this.getUint16(this.shift(Uint16Array.BYTES_PER_ELEMENT), littleEndian);
    }

    /** Write unsigned 16-bit integer number. */
    writeUint16(value: number, littleEndian = this.littleEndian) {
        return this.setUint16(this.shift(Uint16Array.BYTES_PER_ELEMENT), value, littleEndian);
    }

    /** Read signed 16-bit integer number. */
    readInt16(littleEndian = this.littleEndian) {
        return this.getInt16(this.shift(Int16Array.BYTES_PER_ELEMENT), littleEndian);
    }

    /** Write signed 16-bit integer number. */
    writeInt16(value: number, littleEndian = this.littleEndian) {
        return this.setInt16(this.shift(Int16Array.BYTES_PER_ELEMENT), value, littleEndian);
    }

    /** Read unsigned 32-bit integer number. */
    readUint32(littleEndian = this.littleEndian) {
        return this.getUint32(this.shift(Uint32Array.BYTES_PER_ELEMENT), littleEndian);
    }

    /** Write unsigned 32-bit integer number. */
    writeUint32(value: number, littleEndian = this.littleEndian) {
        return this.setUint32(this.shift(Uint32Array.BYTES_PER_ELEMENT), value, littleEndian);
    }

    /** Read unsigned 32-bit integer number. */
    readInt32(littleEndian = this.littleEndian) {
        return this.getInt32(this.shift(Int32Array.BYTES_PER_ELEMENT), littleEndian);
    }

    /** Write signed 32-bit integer number. */
    writeInt32(value: number, littleEndian = this.littleEndian) {
        return this.setInt32(this.shift(Int32Array.BYTES_PER_ELEMENT), value, littleEndian);
    }

    /** Read unsigned 64-bit integer number. */
    readBigUint64(littleEndian = this.littleEndian) {
        return this.getBigUint64(this.shift(BigUint64Array.BYTES_PER_ELEMENT), littleEndian);
    }

    /** Write unsigned 64-bit integer number. */
    writeBigUint64(value: bigint, littleEndian = this.littleEndian) {
        return this.setBigUint64(this.shift(BigUint64Array.BYTES_PER_ELEMENT), value, littleEndian);
    }

    /** Read signed 64-bit integer number. */
    readBigInt64(littleEndian = this.littleEndian) {
        return this.getBigInt64(this.shift(BigInt64Array.BYTES_PER_ELEMENT), littleEndian);
    }

    /** Write signed 64-bit integer number. */
    writeBigInt64(value: bigint, littleEndian = this.littleEndian) {
        return this.setBigInt64(this.shift(BigInt64Array.BYTES_PER_ELEMENT), value, littleEndian);
    }

    /** Read 32-bit float point number. */
    readFloat32(littleEndian = this.littleEndian) {
        return this.getFloat32(this.shift(Float32Array.BYTES_PER_ELEMENT), littleEndian);
    }

    /** Write 32-bit float point number. */
    writeFloat32(value: number, littleEndian = this.littleEndian) {
        return this.setFloat32(this.shift(Float32Array.BYTES_PER_ELEMENT), value, littleEndian);
    }

    /** Read 64-bit float point number. */
    readFloat64(littleEndian = this.littleEndian) {
        return this.getFloat64(this.shift(Float64Array.BYTES_PER_ELEMENT), littleEndian);
    }

    /** Write 64-bit float point number. */
    writeFloat64(value: number, littleEndian = this.littleEndian) {
        return this.setFloat64(this.shift(Float64Array.BYTES_PER_ELEMENT), value, littleEndian);
    }

    /** Read bytes into TypedArray value. */
    readBytes(value: ArrayBufferView) {
        return new Uint8Array(value.buffer, value.byteOffset, value.byteLength).set(new Uint8Array(this.buffer, this.byteOffset + this.shift(value.byteLength), value.byteLength));
    }

    /** Write bytes from TypedArray value. */
    writeBytes(value: ArrayBufferView) {
        return new Uint8Array(this.buffer, this.byteOffset + this.shift(value.byteLength), value.byteLength).set(new Uint8Array(value.buffer, value.byteOffset, value.byteLength));
    }

    /**
     * Read fixed length string.
     * @param length Byte length.
     * @param decoder Text decoder.
     * @returns Output string
     */
    readString(length: number, decoder = new TextDecoder()) {
        return decoder.decode(new Uint8Array(this.buffer, this.byteOffset + this.shift(length), length));
    }

    /**
     * Write string.
     * @param value Input string.
     * @param encoder Text encoder.
     * @returns Bytes written
     */
    writeString(value: string, encoder = new TextEncoder()) {
        const { written = 0 } = encoder.encodeInto(value, new Uint8Array(this.buffer, this.byteOffset + this.offset));
        this.shift(written);
        return written;
    }

    /**
     * Read object(s).
     * @param values
     * @returns
     */
    readObjects(...values: BufferObject[]) {
        let offset = this.offset;

        for (const value of values) {
            this.offset = offset;
            value.readBuffer(this);
            offset += value.byteLength;
        }
    }

    /**
     * Write object(s).
     * @param value
     * @returns 
     */
    writeObjects(...values: BufferObject[]) {
        let offset = this.offset;

        for (const value of values) {
            this.offset = offset;
            value.writeBuffer(this);
            offset += value.byteLength;
        }
    }

    /**
     * Read numbers by specified types.
     * @param format 
     * @returns 
     */
    read(...format: DataType[]) {
        const values = [];

        for (const type of format) {
            switch (type) {
                case DataType.UINT8:
                    values.push(this.readUint8());
                    break;
                case DataType.INT8:
                    values.push(this.readInt8());
                    break;
                case DataType.UINT16:
                    values.push(this.readUint16());
                    break;
                case DataType.INT16:
                    values.push(this.readInt16());
                    break;
                case DataType.UINT32:
                    values.push(this.readUint32());
                    break;
                case DataType.INT32:
                    values.push(this.readInt32());
                    break;
                case DataType.FLOAT32:
                    values.push(this.readFloat32());
                    break;
                case DataType.FLOAT64:
                    values.push(this.readFloat64());
                    break;
            }
        }

        return values;
    }

    /**
     * Write numbers by specified types.
     * @param values 
     * @param format 
     */
    write(values: number[], ...format: DataType[]) {
        let value: number;

        for (const [index, type] of format.entries()) {
            value = values[index] ?? 0;

            switch (type) {
                case DataType.UINT8:
                    this.writeUint8(value);
                    break;
                case DataType.INT8:
                    this.writeInt8(value);
                    break;
                case DataType.UINT16:
                    this.writeUint16(value);
                    break;
                case DataType.INT16:
                    this.writeInt16(value);
                    break;
                case DataType.UINT32:
                    this.writeUint32(value);
                    break;
                case DataType.INT32:
                    this.writeInt32(value);
                    break;
                case DataType.FLOAT32:
                    this.writeFloat32(value);
                    break;
                case DataType.FLOAT64:
                    this.writeFloat64(value);
                    break;
            }
        }
    }
}
