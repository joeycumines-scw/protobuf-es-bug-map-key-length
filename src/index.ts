import {fromBinary, toJsonString} from '@bufbuild/protobuf';
import {MessageSchema} from './gen/schema_pb.js';

const data = new Uint8Array(
  Buffer.from(process.argv[process.argv.length - 1], 'base64'),
);
const res = fromBinary(MessageSchema, data);
console.log(toJsonString(MessageSchema, res));
