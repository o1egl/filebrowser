import React from 'react';
import logo from './logo.svg';
import axios from 'axios';
import './App.css';
import './gen/proto/file/v1/file_service';
import {FileServiceClientJSON, Rpc} from "./gen/proto/file/v1/file_service.twirp-client";
import {FileSortBy, SortOrder} from "./gen/proto/file/v1/file_service";

 function App() {
  const client = axios.create({
    baseURL: "http://localhost:3000/twirp",
  })

  const implementation: Rpc = {
    request(service, method, contentType, data) {
      return client.post(`${service}/${method}`, data, {
        responseType: contentType === "application/protobuf" ? 'arraybuffer' : "json",
        headers: {
          "content-type": contentType,
        }
      }).then(response => {
        return response.data
      });
    }
  }

  const jsonClient = new FileServiceClientJSON(implementation);

 jsonClient.List({path: "/file/path.txt", sortBy: FileSortBy.FILE_SORT_BY_NAME, sortOrder: SortOrder.SORT_ORDER_ASC}).then((value => console.log(value))).catch(reason => console.log(reason))
  //console.log(resp)

  return (
    <div className="App">
      <header className="App-header">
        <img src={logo} className="App-logo" alt="logo" />
        <p>
          Edit <code>src/App.tsx</code> and save to reload.
        </p>
        <a
          className="App-link"
          href="https://reactjs.org"
          target="_blank"
          rel="noopener noreferrer"
        >
          Learn React
        </a>
      </header>
    </div>
  );
}

export default App;
