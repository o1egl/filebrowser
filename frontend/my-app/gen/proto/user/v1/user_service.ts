/* eslint-disable */
import Long from "long";
import _m0 from "protobufjs/minimal";

export const protobufPackage = "user.v1";

export interface FindRequest {}

export interface FindResponse {}

const baseFindRequest: object = {};

export const FindRequest = {
  encode(_: FindRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): FindRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseFindRequest } as FindRequest;
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

  fromJSON(_: any): FindRequest {
    const message = { ...baseFindRequest } as FindRequest;
    return message;
  },

  toJSON(_: FindRequest): unknown {
    const obj: any = {};
    return obj;
  },

  fromPartial(_: DeepPartial<FindRequest>): FindRequest {
    const message = { ...baseFindRequest } as FindRequest;
    return message;
  },
};

const baseFindResponse: object = {};

export const FindResponse = {
  encode(
    _: FindResponse,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): FindResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseFindResponse } as FindResponse;
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

  fromJSON(_: any): FindResponse {
    const message = { ...baseFindResponse } as FindResponse;
    return message;
  },

  toJSON(_: FindResponse): unknown {
    const obj: any = {};
    return obj;
  },

  fromPartial(_: DeepPartial<FindResponse>): FindResponse {
    const message = { ...baseFindResponse } as FindResponse;
    return message;
  },
};

export interface UserService {
  Find(request: FindRequest): Promise<FindResponse>;
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

if (_m0.util.Long !== Long) {
  _m0.util.Long = Long as any;
  _m0.configure();
}
