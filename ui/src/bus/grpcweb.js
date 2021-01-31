import { grpc } from '@improbable-eng/grpc-web';

import { getConfig } from 'adaptor/config';
import { KoalaFS } from 'pb/fs_pb_service';
import { Void } from 'pb/fs_pb';
import { error } from 'bus/notification';

export function invoke(method, request, metadata) {
  return new Promise((resolve) => {
    grpc.invoke(method, {
      request,
      host: getConfig().host,
      metadata,
      onHeaders: (headers) => {
        // console.log(headers);
      },
      onMessage: (message) => {
        return resolve(message);
      },
      onEnd: (code, msg, trailers) => {
        if (code === grpc.Code.OK) {
          console.log('all ok');
        } else {
          console.error('hit an error', code, msg, trailers);
          error(msg);
        }
      },
    });
  });
}

export async function getBranchList() {
  const message = await invoke(KoalaFS.branches, new Void());
  return message.getBranchList();
}
