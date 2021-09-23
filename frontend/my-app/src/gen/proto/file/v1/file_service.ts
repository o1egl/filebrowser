/* eslint-disable */
import Long from "long";
import _m0 from "protobufjs/minimal";
import { Timestamp } from "../../google/protobuf/timestamp";

export const protobufPackage = "file.v1";

export enum FileType {
  FILE_TYPE_UNSPECIFIED = 0,
  FILE_TYPE_BLOB = 1,
  FILE_TYPE_VIDEO = 2,
  FILE_TYPE_AUDIO = 3,
  FILE_TYPE_IMAGE = 4,
  FILE_TYPE_TEXT = 5,
  FILE_TYPE_DIR = 6,
  /** FILE_TYPE_SPECIAL - named pipe, device file, or socket */
  FILE_TYPE_SPECIAL = 7,
  UNRECOGNIZED = -1,
}

export function fileTypeFromJSON(object: any): FileType {
  switch (object) {
    case 0:
    case "FILE_TYPE_UNSPECIFIED":
      return FileType.FILE_TYPE_UNSPECIFIED;
    case 1:
    case "FILE_TYPE_BLOB":
      return FileType.FILE_TYPE_BLOB;
    case 2:
    case "FILE_TYPE_VIDEO":
      return FileType.FILE_TYPE_VIDEO;
    case 3:
    case "FILE_TYPE_AUDIO":
      return FileType.FILE_TYPE_AUDIO;
    case 4:
    case "FILE_TYPE_IMAGE":
      return FileType.FILE_TYPE_IMAGE;
    case 5:
    case "FILE_TYPE_TEXT":
      return FileType.FILE_TYPE_TEXT;
    case 6:
    case "FILE_TYPE_DIR":
      return FileType.FILE_TYPE_DIR;
    case 7:
    case "FILE_TYPE_SPECIAL":
      return FileType.FILE_TYPE_SPECIAL;
    case -1:
    case "UNRECOGNIZED":
    default:
      return FileType.UNRECOGNIZED;
  }
}

export function fileTypeToJSON(object: FileType): string {
  switch (object) {
    case FileType.FILE_TYPE_UNSPECIFIED:
      return "FILE_TYPE_UNSPECIFIED";
    case FileType.FILE_TYPE_BLOB:
      return "FILE_TYPE_BLOB";
    case FileType.FILE_TYPE_VIDEO:
      return "FILE_TYPE_VIDEO";
    case FileType.FILE_TYPE_AUDIO:
      return "FILE_TYPE_AUDIO";
    case FileType.FILE_TYPE_IMAGE:
      return "FILE_TYPE_IMAGE";
    case FileType.FILE_TYPE_TEXT:
      return "FILE_TYPE_TEXT";
    case FileType.FILE_TYPE_DIR:
      return "FILE_TYPE_DIR";
    case FileType.FILE_TYPE_SPECIAL:
      return "FILE_TYPE_SPECIAL";
    default:
      return "UNKNOWN";
  }
}

export enum FileSortBy {
  /** FILE_SORT_BY_NAME - buf:lint:ignore ENUM_ZERO_VALUE_SUFFIX */
  FILE_SORT_BY_NAME = 0,
  FILE_SORT_BY_SIZE = 1,
  FILE_SORT_BY_MOD_TIME = 2,
  UNRECOGNIZED = -1,
}

export function fileSortByFromJSON(object: any): FileSortBy {
  switch (object) {
    case 0:
    case "FILE_SORT_BY_NAME":
      return FileSortBy.FILE_SORT_BY_NAME;
    case 1:
    case "FILE_SORT_BY_SIZE":
      return FileSortBy.FILE_SORT_BY_SIZE;
    case 2:
    case "FILE_SORT_BY_MOD_TIME":
      return FileSortBy.FILE_SORT_BY_MOD_TIME;
    case -1:
    case "UNRECOGNIZED":
    default:
      return FileSortBy.UNRECOGNIZED;
  }
}

export function fileSortByToJSON(object: FileSortBy): string {
  switch (object) {
    case FileSortBy.FILE_SORT_BY_NAME:
      return "FILE_SORT_BY_NAME";
    case FileSortBy.FILE_SORT_BY_SIZE:
      return "FILE_SORT_BY_SIZE";
    case FileSortBy.FILE_SORT_BY_MOD_TIME:
      return "FILE_SORT_BY_MOD_TIME";
    default:
      return "UNKNOWN";
  }
}

export enum SortOrder {
  /** SORT_ORDER_DESC - buf:lint:ignore ENUM_ZERO_VALUE_SUFFIX */
  SORT_ORDER_DESC = 0,
  SORT_ORDER_ASC = 1,
  UNRECOGNIZED = -1,
}

export function sortOrderFromJSON(object: any): SortOrder {
  switch (object) {
    case 0:
    case "SORT_ORDER_DESC":
      return SortOrder.SORT_ORDER_DESC;
    case 1:
    case "SORT_ORDER_ASC":
      return SortOrder.SORT_ORDER_ASC;
    case -1:
    case "UNRECOGNIZED":
    default:
      return SortOrder.UNRECOGNIZED;
  }
}

export function sortOrderToJSON(object: SortOrder): string {
  switch (object) {
    case SortOrder.SORT_ORDER_DESC:
      return "SORT_ORDER_DESC";
    case SortOrder.SORT_ORDER_ASC:
      return "SORT_ORDER_ASC";
    default:
      return "UNKNOWN";
  }
}

export interface FileInfoRequest {
  path: string;
}

export interface FileInfoResponse {
  info: FileInfo | undefined;
}

export interface FileInfo {
  path: string;
  name: string;
  size: number;
  type: FileType;
  isSymlink: boolean;
  modTime: Date | undefined;
  mode: number;
}

export interface ListRequest {
  path: string;
  sortBy: FileSortBy;
  sortOrder: SortOrder;
}

export interface ListResponse {
  info: FileInfo | undefined;
  children: FileInfo[];
}

export interface CreateFileRequest {
  path: string;
  override: boolean;
  content: Uint8Array;
}

export interface CreateFileResponse {}

export interface CreateDirRequest {
  path: string;
}

export interface CreateDirResponse {}

export interface RemoveRequest {
  path: string;
  force: boolean;
}

export interface RemoveResponse {}

const baseFileInfoRequest: object = { path: "" };

export const FileInfoRequest = {
  encode(
    message: FileInfoRequest,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    if (message.path !== "") {
      writer.uint32(10).string(message.path);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): FileInfoRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseFileInfoRequest } as FileInfoRequest;
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.path = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): FileInfoRequest {
    const message = { ...baseFileInfoRequest } as FileInfoRequest;
    if (object.path !== undefined && object.path !== null) {
      message.path = String(object.path);
    } else {
      message.path = "";
    }
    return message;
  },

  toJSON(message: FileInfoRequest): unknown {
    const obj: any = {};
    message.path !== undefined && (obj.path = message.path);
    return obj;
  },

  fromPartial(object: DeepPartial<FileInfoRequest>): FileInfoRequest {
    const message = { ...baseFileInfoRequest } as FileInfoRequest;
    if (object.path !== undefined && object.path !== null) {
      message.path = object.path;
    } else {
      message.path = "";
    }
    return message;
  },
};

const baseFileInfoResponse: object = {};

export const FileInfoResponse = {
  encode(
    message: FileInfoResponse,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    if (message.info !== undefined) {
      FileInfo.encode(message.info, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): FileInfoResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseFileInfoResponse } as FileInfoResponse;
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.info = FileInfo.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): FileInfoResponse {
    const message = { ...baseFileInfoResponse } as FileInfoResponse;
    if (object.info !== undefined && object.info !== null) {
      message.info = FileInfo.fromJSON(object.info);
    } else {
      message.info = undefined;
    }
    return message;
  },

  toJSON(message: FileInfoResponse): unknown {
    const obj: any = {};
    message.info !== undefined &&
      (obj.info = message.info ? FileInfo.toJSON(message.info) : undefined);
    return obj;
  },

  fromPartial(object: DeepPartial<FileInfoResponse>): FileInfoResponse {
    const message = { ...baseFileInfoResponse } as FileInfoResponse;
    if (object.info !== undefined && object.info !== null) {
      message.info = FileInfo.fromPartial(object.info);
    } else {
      message.info = undefined;
    }
    return message;
  },
};

const baseFileInfo: object = {
  path: "",
  name: "",
  size: 0,
  type: 0,
  isSymlink: false,
  mode: 0,
};

export const FileInfo = {
  encode(
    message: FileInfo,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    if (message.path !== "") {
      writer.uint32(10).string(message.path);
    }
    if (message.name !== "") {
      writer.uint32(18).string(message.name);
    }
    if (message.size !== 0) {
      writer.uint32(24).int64(message.size);
    }
    if (message.type !== 0) {
      writer.uint32(32).int32(message.type);
    }
    if (message.isSymlink === true) {
      writer.uint32(40).bool(message.isSymlink);
    }
    if (message.modTime !== undefined) {
      Timestamp.encode(
        toTimestamp(message.modTime),
        writer.uint32(50).fork()
      ).ldelim();
    }
    if (message.mode !== 0) {
      writer.uint32(56).uint32(message.mode);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): FileInfo {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseFileInfo } as FileInfo;
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.path = reader.string();
          break;
        case 2:
          message.name = reader.string();
          break;
        case 3:
          message.size = longToNumber(reader.int64() as Long);
          break;
        case 4:
          message.type = reader.int32() as any;
          break;
        case 5:
          message.isSymlink = reader.bool();
          break;
        case 6:
          message.modTime = fromTimestamp(
            Timestamp.decode(reader, reader.uint32())
          );
          break;
        case 7:
          message.mode = reader.uint32();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): FileInfo {
    const message = { ...baseFileInfo } as FileInfo;
    if (object.path !== undefined && object.path !== null) {
      message.path = String(object.path);
    } else {
      message.path = "";
    }
    if (object.name !== undefined && object.name !== null) {
      message.name = String(object.name);
    } else {
      message.name = "";
    }
    if (object.size !== undefined && object.size !== null) {
      message.size = Number(object.size);
    } else {
      message.size = 0;
    }
    if (object.type !== undefined && object.type !== null) {
      message.type = fileTypeFromJSON(object.type);
    } else {
      message.type = 0;
    }
    if (object.isSymlink !== undefined && object.isSymlink !== null) {
      message.isSymlink = Boolean(object.isSymlink);
    } else {
      message.isSymlink = false;
    }
    if (object.modTime !== undefined && object.modTime !== null) {
      message.modTime = fromJsonTimestamp(object.modTime);
    } else {
      message.modTime = undefined;
    }
    if (object.mode !== undefined && object.mode !== null) {
      message.mode = Number(object.mode);
    } else {
      message.mode = 0;
    }
    return message;
  },

  toJSON(message: FileInfo): unknown {
    const obj: any = {};
    message.path !== undefined && (obj.path = message.path);
    message.name !== undefined && (obj.name = message.name);
    message.size !== undefined && (obj.size = message.size);
    message.type !== undefined && (obj.type = fileTypeToJSON(message.type));
    message.isSymlink !== undefined && (obj.isSymlink = message.isSymlink);
    message.modTime !== undefined &&
      (obj.modTime = message.modTime.toISOString());
    message.mode !== undefined && (obj.mode = message.mode);
    return obj;
  },

  fromPartial(object: DeepPartial<FileInfo>): FileInfo {
    const message = { ...baseFileInfo } as FileInfo;
    if (object.path !== undefined && object.path !== null) {
      message.path = object.path;
    } else {
      message.path = "";
    }
    if (object.name !== undefined && object.name !== null) {
      message.name = object.name;
    } else {
      message.name = "";
    }
    if (object.size !== undefined && object.size !== null) {
      message.size = object.size;
    } else {
      message.size = 0;
    }
    if (object.type !== undefined && object.type !== null) {
      message.type = object.type;
    } else {
      message.type = 0;
    }
    if (object.isSymlink !== undefined && object.isSymlink !== null) {
      message.isSymlink = object.isSymlink;
    } else {
      message.isSymlink = false;
    }
    if (object.modTime !== undefined && object.modTime !== null) {
      message.modTime = object.modTime;
    } else {
      message.modTime = undefined;
    }
    if (object.mode !== undefined && object.mode !== null) {
      message.mode = object.mode;
    } else {
      message.mode = 0;
    }
    return message;
  },
};

const baseListRequest: object = { path: "", sortBy: 0, sortOrder: 0 };

export const ListRequest = {
  encode(
    message: ListRequest,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    if (message.path !== "") {
      writer.uint32(10).string(message.path);
    }
    if (message.sortBy !== 0) {
      writer.uint32(16).int32(message.sortBy);
    }
    if (message.sortOrder !== 0) {
      writer.uint32(24).int32(message.sortOrder);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ListRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseListRequest } as ListRequest;
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.path = reader.string();
          break;
        case 2:
          message.sortBy = reader.int32() as any;
          break;
        case 3:
          message.sortOrder = reader.int32() as any;
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): ListRequest {
    const message = { ...baseListRequest } as ListRequest;
    if (object.path !== undefined && object.path !== null) {
      message.path = String(object.path);
    } else {
      message.path = "";
    }
    if (object.sortBy !== undefined && object.sortBy !== null) {
      message.sortBy = fileSortByFromJSON(object.sortBy);
    } else {
      message.sortBy = 0;
    }
    if (object.sortOrder !== undefined && object.sortOrder !== null) {
      message.sortOrder = sortOrderFromJSON(object.sortOrder);
    } else {
      message.sortOrder = 0;
    }
    return message;
  },

  toJSON(message: ListRequest): unknown {
    const obj: any = {};
    message.path !== undefined && (obj.path = message.path);
    message.sortBy !== undefined &&
      (obj.sortBy = fileSortByToJSON(message.sortBy));
    message.sortOrder !== undefined &&
      (obj.sortOrder = sortOrderToJSON(message.sortOrder));
    return obj;
  },

  fromPartial(object: DeepPartial<ListRequest>): ListRequest {
    const message = { ...baseListRequest } as ListRequest;
    if (object.path !== undefined && object.path !== null) {
      message.path = object.path;
    } else {
      message.path = "";
    }
    if (object.sortBy !== undefined && object.sortBy !== null) {
      message.sortBy = object.sortBy;
    } else {
      message.sortBy = 0;
    }
    if (object.sortOrder !== undefined && object.sortOrder !== null) {
      message.sortOrder = object.sortOrder;
    } else {
      message.sortOrder = 0;
    }
    return message;
  },
};

const baseListResponse: object = {};

export const ListResponse = {
  encode(
    message: ListResponse,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    if (message.info !== undefined) {
      FileInfo.encode(message.info, writer.uint32(10).fork()).ldelim();
    }
    for (const v of message.children) {
      FileInfo.encode(v!, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ListResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseListResponse } as ListResponse;
    message.children = [];
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.info = FileInfo.decode(reader, reader.uint32());
          break;
        case 2:
          message.children.push(FileInfo.decode(reader, reader.uint32()));
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): ListResponse {
    const message = { ...baseListResponse } as ListResponse;
    message.children = [];
    if (object.info !== undefined && object.info !== null) {
      message.info = FileInfo.fromJSON(object.info);
    } else {
      message.info = undefined;
    }
    if (object.children !== undefined && object.children !== null) {
      for (const e of object.children) {
        message.children.push(FileInfo.fromJSON(e));
      }
    }
    return message;
  },

  toJSON(message: ListResponse): unknown {
    const obj: any = {};
    message.info !== undefined &&
      (obj.info = message.info ? FileInfo.toJSON(message.info) : undefined);
    if (message.children) {
      obj.children = message.children.map((e) =>
        e ? FileInfo.toJSON(e) : undefined
      );
    } else {
      obj.children = [];
    }
    return obj;
  },

  fromPartial(object: DeepPartial<ListResponse>): ListResponse {
    const message = { ...baseListResponse } as ListResponse;
    message.children = [];
    if (object.info !== undefined && object.info !== null) {
      message.info = FileInfo.fromPartial(object.info);
    } else {
      message.info = undefined;
    }
    if (object.children !== undefined && object.children !== null) {
      for (const e of object.children) {
        message.children.push(FileInfo.fromPartial(e));
      }
    }
    return message;
  },
};

const baseCreateFileRequest: object = { path: "", override: false };

export const CreateFileRequest = {
  encode(
    message: CreateFileRequest,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    if (message.path !== "") {
      writer.uint32(10).string(message.path);
    }
    if (message.override === true) {
      writer.uint32(16).bool(message.override);
    }
    if (message.content.length !== 0) {
      writer.uint32(26).bytes(message.content);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): CreateFileRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseCreateFileRequest } as CreateFileRequest;
    message.content = new Uint8Array();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.path = reader.string();
          break;
        case 2:
          message.override = reader.bool();
          break;
        case 3:
          message.content = reader.bytes();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): CreateFileRequest {
    const message = { ...baseCreateFileRequest } as CreateFileRequest;
    message.content = new Uint8Array();
    if (object.path !== undefined && object.path !== null) {
      message.path = String(object.path);
    } else {
      message.path = "";
    }
    if (object.override !== undefined && object.override !== null) {
      message.override = Boolean(object.override);
    } else {
      message.override = false;
    }
    if (object.content !== undefined && object.content !== null) {
      message.content = bytesFromBase64(object.content);
    }
    return message;
  },

  toJSON(message: CreateFileRequest): unknown {
    const obj: any = {};
    message.path !== undefined && (obj.path = message.path);
    message.override !== undefined && (obj.override = message.override);
    message.content !== undefined &&
      (obj.content = base64FromBytes(
        message.content !== undefined ? message.content : new Uint8Array()
      ));
    return obj;
  },

  fromPartial(object: DeepPartial<CreateFileRequest>): CreateFileRequest {
    const message = { ...baseCreateFileRequest } as CreateFileRequest;
    if (object.path !== undefined && object.path !== null) {
      message.path = object.path;
    } else {
      message.path = "";
    }
    if (object.override !== undefined && object.override !== null) {
      message.override = object.override;
    } else {
      message.override = false;
    }
    if (object.content !== undefined && object.content !== null) {
      message.content = object.content;
    } else {
      message.content = new Uint8Array();
    }
    return message;
  },
};

const baseCreateFileResponse: object = {};

export const CreateFileResponse = {
  encode(
    _: CreateFileResponse,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): CreateFileResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseCreateFileResponse } as CreateFileResponse;
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(_: any): CreateFileResponse {
    const message = { ...baseCreateFileResponse } as CreateFileResponse;
    return message;
  },

  toJSON(_: CreateFileResponse): unknown {
    const obj: any = {};
    return obj;
  },

  fromPartial(_: DeepPartial<CreateFileResponse>): CreateFileResponse {
    const message = { ...baseCreateFileResponse } as CreateFileResponse;
    return message;
  },
};

const baseCreateDirRequest: object = { path: "" };

export const CreateDirRequest = {
  encode(
    message: CreateDirRequest,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    if (message.path !== "") {
      writer.uint32(10).string(message.path);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): CreateDirRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseCreateDirRequest } as CreateDirRequest;
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.path = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): CreateDirRequest {
    const message = { ...baseCreateDirRequest } as CreateDirRequest;
    if (object.path !== undefined && object.path !== null) {
      message.path = String(object.path);
    } else {
      message.path = "";
    }
    return message;
  },

  toJSON(message: CreateDirRequest): unknown {
    const obj: any = {};
    message.path !== undefined && (obj.path = message.path);
    return obj;
  },

  fromPartial(object: DeepPartial<CreateDirRequest>): CreateDirRequest {
    const message = { ...baseCreateDirRequest } as CreateDirRequest;
    if (object.path !== undefined && object.path !== null) {
      message.path = object.path;
    } else {
      message.path = "";
    }
    return message;
  },
};

const baseCreateDirResponse: object = {};

export const CreateDirResponse = {
  encode(
    _: CreateDirResponse,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): CreateDirResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseCreateDirResponse } as CreateDirResponse;
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(_: any): CreateDirResponse {
    const message = { ...baseCreateDirResponse } as CreateDirResponse;
    return message;
  },

  toJSON(_: CreateDirResponse): unknown {
    const obj: any = {};
    return obj;
  },

  fromPartial(_: DeepPartial<CreateDirResponse>): CreateDirResponse {
    const message = { ...baseCreateDirResponse } as CreateDirResponse;
    return message;
  },
};

const baseRemoveRequest: object = { path: "", force: false };

export const RemoveRequest = {
  encode(
    message: RemoveRequest,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    if (message.path !== "") {
      writer.uint32(10).string(message.path);
    }
    if (message.force === true) {
      writer.uint32(16).bool(message.force);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): RemoveRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseRemoveRequest } as RemoveRequest;
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.path = reader.string();
          break;
        case 2:
          message.force = reader.bool();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): RemoveRequest {
    const message = { ...baseRemoveRequest } as RemoveRequest;
    if (object.path !== undefined && object.path !== null) {
      message.path = String(object.path);
    } else {
      message.path = "";
    }
    if (object.force !== undefined && object.force !== null) {
      message.force = Boolean(object.force);
    } else {
      message.force = false;
    }
    return message;
  },

  toJSON(message: RemoveRequest): unknown {
    const obj: any = {};
    message.path !== undefined && (obj.path = message.path);
    message.force !== undefined && (obj.force = message.force);
    return obj;
  },

  fromPartial(object: DeepPartial<RemoveRequest>): RemoveRequest {
    const message = { ...baseRemoveRequest } as RemoveRequest;
    if (object.path !== undefined && object.path !== null) {
      message.path = object.path;
    } else {
      message.path = "";
    }
    if (object.force !== undefined && object.force !== null) {
      message.force = object.force;
    } else {
      message.force = false;
    }
    return message;
  },
};

const baseRemoveResponse: object = {};

export const RemoveResponse = {
  encode(
    _: RemoveResponse,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): RemoveResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseRemoveResponse } as RemoveResponse;
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(_: any): RemoveResponse {
    const message = { ...baseRemoveResponse } as RemoveResponse;
    return message;
  },

  toJSON(_: RemoveResponse): unknown {
    const obj: any = {};
    return obj;
  },

  fromPartial(_: DeepPartial<RemoveResponse>): RemoveResponse {
    const message = { ...baseRemoveResponse } as RemoveResponse;
    return message;
  },
};

export interface FileService {
  FileInfo(request: FileInfoRequest): Promise<FileInfoResponse>;
  List(request: ListRequest): Promise<ListResponse>;
  CreateFile(request: CreateFileRequest): Promise<CreateFileResponse>;
  CreateDir(request: CreateDirRequest): Promise<CreateDirResponse>;
  Remove(request: RemoveRequest): Promise<RemoveResponse>;
}

declare var self: any | undefined;
declare var window: any | undefined;
declare var global: any | undefined;
var globalThis: any = (() => {
  if (typeof globalThis !== "undefined") return globalThis;
  if (typeof self !== "undefined") return self;
  if (typeof window !== "undefined") return window;
  if (typeof global !== "undefined") return global;
  throw "Unable to locate global object";
})();

const atob: (b64: string) => string =
  globalThis.atob ||
  ((b64) => globalThis.Buffer.from(b64, "base64").toString("binary"));
function bytesFromBase64(b64: string): Uint8Array {
  const bin = atob(b64);
  const arr = new Uint8Array(bin.length);
  for (let i = 0; i < bin.length; ++i) {
    arr[i] = bin.charCodeAt(i);
  }
  return arr;
}

const btoa: (bin: string) => string =
  globalThis.btoa ||
  ((bin) => globalThis.Buffer.from(bin, "binary").toString("base64"));
function base64FromBytes(arr: Uint8Array): string {
  const bin: string[] = [];
  for (const byte of arr) {
    bin.push(String.fromCharCode(byte));
  }
  return btoa(bin.join(""));
}

type Builtin =
  | Date
  | Function
  | Uint8Array
  | string
  | number
  | boolean
  | undefined;
export type DeepPartial<T> = T extends Builtin
  ? T
  : T extends Array<infer U>
  ? Array<DeepPartial<U>>
  : T extends ReadonlyArray<infer U>
  ? ReadonlyArray<DeepPartial<U>>
  : T extends {}
  ? { [K in keyof T]?: DeepPartial<T[K]> }
  : Partial<T>;

function toTimestamp(date: Date): Timestamp {
  const seconds = date.getTime() / 1_000;
  const nanos = (date.getTime() % 1_000) * 1_000_000;
  return { seconds, nanos };
}

function fromTimestamp(t: Timestamp): Date {
  let millis = t.seconds * 1_000;
  millis += t.nanos / 1_000_000;
  return new Date(millis);
}

function fromJsonTimestamp(o: any): Date {
  if (o instanceof Date) {
    return o;
  } else if (typeof o === "string") {
    return new Date(o);
  } else {
    return fromTimestamp(Timestamp.fromJSON(o));
  }
}

function longToNumber(long: Long): number {
  if (long.gt(Number.MAX_SAFE_INTEGER)) {
    throw new globalThis.Error("Value is larger than Number.MAX_SAFE_INTEGER");
  }
  return long.toNumber();
}

if (_m0.util.Long !== Long) {
  _m0.util.Long = Long as any;
  _m0.configure();
}
