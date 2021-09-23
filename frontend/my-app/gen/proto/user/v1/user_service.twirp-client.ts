import { FindRequest, FindResponse } from "./user_service";

//==================================//
//          Client Code             //
//==================================//

interface Rpc {
  request(
    service: string,
    method: string,
    contentType: "application/json" | "application/protobuf",
    data: object | Uint8Array
  ): Promise<object | Uint8Array>;
}

export interface UserServiceClient {
  Find(request: FindRequest): Promise<FindResponse>;
}

export class UserServiceClientJSON implements UserServiceClient {
  private readonly rpc: Rpc;
  constructor(rpc: Rpc) {
    this.rpc = rpc;
    this.Find.bind(this);
  }
  Find(request: FindRequest): Promise<FindResponse> {
    const data = FindRequest.toJSON(request);
    const promise = this.rpc.request(
      "user.v1.UserService",
      "Find",
      "application/json",
      data as object
    );
    return promise.then((data) => FindResponse.fromJSON(data as any));
  }
}

export class UserServiceClientProtobuf implements UserServiceClient {
  private readonly rpc: Rpc;
  constructor(rpc: Rpc) {
    this.rpc = rpc;
    this.Find.bind(this);
  }
  Find(request: FindRequest): Promise<FindResponse> {
    const data = FindRequest.encode(request).finish();
    const promise = this.rpc.request(
      "user.v1.UserService",
      "Find",
      "application/protobuf",
      data
    );
    return promise.then((data) => FindResponse.decode(data as Uint8Array));
  }
}
