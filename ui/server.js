/*
 *
 * Copyright 2018 Google LLC
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

const PROTO_PATH = `${__dirname}/protos/fs.proto`;

const assert = require('assert');
const grpc = require('@grpc/grpc-js');
const protoLoader = require('@grpc/proto-loader');

const packageDefinition = protoLoader.loadSync(
  PROTO_PATH,
  {
    keepCase: true,
    longs: String,
    enums: String,
    defaults: true,
    oneofs: true,
  },
);
const protoDescriptor = grpc.loadPackageDefinition(packageDefinition);
const { fileSystem } = protoDescriptor.grpc.koala;
const implementation = require('./fs');
/**
 * Get a new server with the handler functions in this file bound to the
 * methods it serves.
 * @return {!Server} The new server object
 */
function getServer() {
  const server = new grpc.Server();
  server.addService(fileSystem.KoalaFS.service, implementation);
  return server;
}

const server = getServer();
server.bindAsync(
  '0.0.0.0:9090', grpc.ServerCredentials.createInsecure(), (err, port) => {
    console.log(err, port);
    assert.ifError(err);
    server.start();
  },
);
