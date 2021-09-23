import {
  FileInfoRequest,
  FileInfoResponse,
  ListRequest,
  ListResponse,
  CreateFileRequest,
  CreateFileResponse,
  CreateDirRequest,
  CreateDirResponse,
  RemoveRequest,
  RemoveResponse,
} from "./file_service";

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

export interface FileServiceClient {
  FileInfo(request: FileInfoRequest): Promise<FileInfoResponse>;
  List(request: ListRequest): Promise<ListResponse>;
  CreateFile(request: CreateFileRequest): Promise<CreateFileResponse>;
  CreateDir(request: CreateDirRequest): Promise<CreateDirResponse>;
  Remove(request: RemoveRequest): Promise<RemoveResponse>;
}

export class FileServiceClientJSON implements FileServiceClient {
  private readonly rpc: Rpc;
  constructor(rpc: Rpc) {
    this.rpc = rpc;
    this.FileInfo.bind(this);
    this.List.bind(this);
    this.CreateFile.bind(this);
    this.CreateDir.bind(this);
    this.Remove.bind(this);
  }
  FileInfo(request: FileInfoRequest): Promise<FileInfoResponse> {
    const data = FileInfoRequest.toJSON(request);
    const promise = this.rpc.request(
      "file.v1.FileService",
      "FileInfo",
      "application/json",
      data as object
    );
    return promise.then((data) => FileInfoResponse.fromJSON(data as any));
  }

  List(request: ListRequest): Promise<ListResponse> {
    const data = ListRequest.toJSON(request);
    const promise = this.rpc.request(
      "file.v1.FileService",
      "List",
      "application/json",
      data as object
    );
    return promise.then((data) => ListResponse.fromJSON(data as any));
  }

  CreateFile(request: CreateFileRequest): Promise<CreateFileResponse> {
    const data = CreateFileRequest.toJSON(request);
    const promise = this.rpc.request(
      "file.v1.FileService",
      "CreateFile",
      "application/json",
      data as object
    );
    return promise.then((data) => CreateFileResponse.fromJSON(data as any));
  }

  CreateDir(request: CreateDirRequest): Promise<CreateDirResponse> {
    const data = CreateDirRequest.toJSON(request);
    const promise = this.rpc.request(
      "file.v1.FileService",
      "CreateDir",
      "application/json",
      data as object
    );
    return promise.then((data) => CreateDirResponse.fromJSON(data as any));
  }

  Remove(request: RemoveRequest): Promise<RemoveResponse> {
    const data = RemoveRequest.toJSON(request);
    const promise = this.rpc.request(
      "file.v1.FileService",
      "Remove",
      "application/json",
      data as object
    );
    return promise.then((data) => RemoveResponse.fromJSON(data as any));
  }
}

export class FileServiceClientProtobuf implements FileServiceClient {
  private readonly rpc: Rpc;
  constructor(rpc: Rpc) {
    this.rpc = rpc;
    this.FileInfo.bind(this);
    this.List.bind(this);
    this.CreateFile.bind(this);
    this.CreateDir.bind(this);
    this.Remove.bind(this);
  }
  FileInfo(request: FileInfoRequest): Promise<FileInfoResponse> {
    const data = FileInfoRequest.encode(request).finish();
    const promise = this.rpc.request(
      "file.v1.FileService",
      "FileInfo",
      "application/protobuf",
      data
    );
    return promise.then((data) => FileInfoResponse.decode(data as Uint8Array));
  }

  List(request: ListRequest): Promise<ListResponse> {
    const data = ListRequest.encode(request).finish();
    const promise = this.rpc.request(
      "file.v1.FileService",
      "List",
      "application/protobuf",
      data
    );
    return promise.then((data) => ListResponse.decode(data as Uint8Array));
  }

  CreateFile(request: CreateFileRequest): Promise<CreateFileResponse> {
    const data = CreateFileRequest.encode(request).finish();
    const promise = this.rpc.request(
      "file.v1.FileService",
      "CreateFile",
      "application/protobuf",
      data
    );
    return promise.then((data) =>
      CreateFileResponse.decode(data as Uint8Array)
    );
  }

  CreateDir(request: CreateDirRequest): Promise<CreateDirResponse> {
    const data = CreateDirRequest.encode(request).finish();
    const promise = this.rpc.request(
      "file.v1.FileService",
      "CreateDir",
      "application/protobuf",
      data
    );
    return promise.then((data) => CreateDirResponse.decode(data as Uint8Array));
  }

  Remove(request: RemoveRequest): Promise<RemoveResponse> {
    const data = RemoveRequest.encode(request).finish();
    const promise = this.rpc.request(
      "file.v1.FileService",
      "Remove",
      "application/protobuf",
      data
    );
    return promise.then((data) => RemoveResponse.decode(data as Uint8Array));
  }
}
