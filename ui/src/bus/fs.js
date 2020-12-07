import {
  PathRequest, MoveRequest, PathListRequest, DownloadRequest, UploadRequest,
} from 'pb/fs_pb';
import { KoalaFS } from 'pb/fs_pb_service';
import { setState, busState } from 'bus/bus';
import { error } from 'bus/notification';
import { dirname, basename, join } from 'utils/filepath';

import { grpc } from '@improbable-eng/grpc-web';
import map from 'promise.map';

function invoke(method, request, metadata) {
  return new Promise((resolve) => {
    grpc.invoke(method, {
      request,
      host: 'https://localhost:9091',
      metadata: Object.assign(metadata || {}, { 'kfs-pwd': busState.pwd, 'kfs-mount': 'default' }),
      onHeaders: (headers) => {
        // console.log(headers);
      },
      onMessage: (message) => {
        if (message.getFilesList) {
          const files = message.getFilesList().map((f) => f.toObject());
          window.message = message;
          if (files) {
            setState({ files });
          }
        }
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

export async function cd(path) {
  try {
    console.log('---grpc cd---', path);
    const message = await invoke(KoalaFS.ls,
      new PathRequest().setPath(path));
    console.log('---grpc cd cb---', message);
    const { path: pwd } = message.toObject();
    if (pwd !== busState.pwd) {
      setState({ chosenFiles: {} });
    }
    setState({
      pwd,
      chosen: (_chosen) => {
        Object.keys(_chosen).forEach((item) => {
          delete _chosen[item];
        });
      },
    });
  } catch (e) {
    console.error('---grpc cd error---', e);
    error('读取文件夹', e.message);
  }
}

export async function mv(srcList, dst) {
  try {
    if (srcList.length === 1 && srcList[0] === dst) {
      return;
    }
    console.log('---grpc mv---', srcList, dst);
    const message = await invoke(KoalaFS.mv,
      new MoveRequest().setSrcList(srcList).setDst(dst));
    console.log('---grpc mv cb---', message);
    if (srcList.length !== 1 || dirname(srcList[0]) === dirname(dst)) {
      await cd(busState.pwd);
    }
    setState({
      chosen: (_chosen) => {
        Object.keys(_chosen).forEach((item) => {
          delete _chosen[item];
        });
        console.log(dst, dirname(dst), busState.pwd);
        if (dirname(dst) === busState.pwd
          && Object.values(busState.files).find((f) => f.name === basename(dst)).type === 'dir') {
          _chosen[dst] = 1;
          return;
        }
        srcList.forEach((src) => {
          _chosen[`${busState.pwd}/${basename(src)}`] = 1;
        });
      },
    });
  } catch (e) {
    console.error('---grpc mv error---', e);
    error('移动文件', e.message);
  }
}

export async function cp(srcList, dst) {
  try {
    if (srcList.length === 1 && srcList[0] === dst) {
      return;
    }
    console.log('---grpc cp---', srcList, dst);
    const message = await invoke(KoalaFS.cp,
      new MoveRequest().setSrcList(srcList).setDst(dst));
    console.log('---grpc cp cb---', message);
    if (srcList.length !== 1 || dirname(srcList[0]) === dirname(dst)) {
      await cd(busState.pwd);
    }
    setState({
      chosen: (_chosen) => {
        Object.keys(_chosen).forEach((item) => {
          delete _chosen[item];
        });
        srcList.forEach((src) => {
          _chosen[`${busState.pwd}/${basename(src)}`] = 1;
        });
      },
    });
  } catch (e) {
    console.error('---grpc cp error---', e);
    error('复制文件', e.message);
  }
}

export async function createFile(path) {
  try {
    console.log('---grpc createFile---', path);
    const message = await invoke(KoalaFS.createFile,
      new PathRequest().setPath(path));
    console.log('---grpc createFile cb---', message);
    const { name } = message.toObject();
    setState({
      chosen: (_chosen) => {
        Object.keys(_chosen).forEach((item) => {
          delete _chosen[item];
        });
        _chosen[join(busState.pwd, name)] = 2;
      },
    });
  } catch (e) {
    console.error('---grpc createFile error---', e);
    error('新建文件', e.message);
  }
}

export async function mkdir(path) {
  try {
    console.log('---grpc mkdir---');
    const { pwd } = busState;
    const message = await invoke(KoalaFS.mkdir,
      new PathRequest().setPath(path));
    console.log('---grpc mkdir cb---', message);
    const { name } = message.toObject();
    setState({
      chosen: (_chosen) => {
        Object.keys(_chosen).forEach((item) => {
          delete _chosen[item];
        });
        _chosen[join(pwd, name)] = 2;
      },
    });
  } catch (e) {
    console.error('---grpc mkdir error---', e);
    error('新建文件夹', e.message);
  }
}

export async function remove(path) {
  try {
    const pathList = typeof path === 'string' ? [path] : path;
    console.log('---grpc remove---', path);
    const message = await invoke(KoalaFS.remove,
      new PathListRequest().setPathList(pathList));
    console.log('---grpc remove cb---', message);
    await cd(busState.pwd);
    setState({
      chosen: (_chosen) => {
        pathList.forEach((path) => delete _chosen[path]);
      },
    });
  } catch (e) {
    console.error('---grpc mv error---', e);
    error('删除文件或文件夹', e.message);
  }
}

function createAndDownloadFile(fileName, content) {
  const aTag = document.createElement('a');
  const blob = new Blob(content);
  aTag.download = fileName;
  aTag.href = URL.createObjectURL(blob);
  aTag.click();
  URL.revokeObjectURL(blob);
}

export async function download(pathList) {
  try {
    console.log('---grpc download---', pathList);
    const message = await invoke(KoalaFS.download,
      new DownloadRequest().setPathList(pathList));
    console.log('---grpc download cb---', message);
    for (const hash of message.getHashList()) {
      fetch("http://localhost:9999/api/download/" + hash)
        .then(response => response.blob())
        .then(blob => {
          const aTag = document.createElement('a');
          aTag.download = basename(pathList[0]);
          aTag.href = URL.createObjectURL(blob);
          aTag.click();
          URL.revokeObjectURL(blob);
        })
    }
    // const hashList = message.getHashList();
    // const blockCount = hashList.length;
    // if (hashList && blockCount > 0) {
    //   const blobs = await map(hashList, async (hash, i) => {
    //     const hashHex = hash.map((x) => (`00${x.toString(16)}`).slice(-2)).join('');
    //     console.log(`download block ${i}/${blockCount}---`, hashHex);
    //     const message = await invoke(KoalaFS.download,
    //       new DownloadRequest().setHash(hash));
    //     console.log(`download block ${i}/${blockCount} cb---`, hashHex, message);
    //     return message.getSinglefilecontent();
    //   }, 2);
    //   createAndDownloadFile(basename(pathList[0]), [...blobs]);
    // } else {
    //   createAndDownloadFile(basename(pathList[0]), [message.getSinglefilecontent()]);
    // }
  } catch (e) {
    console.error('---grpc download error---', e);
    error('下载文件', e.message);
  }
}

export async function upload(path, data, hashList = []) {
  try {
    let hash = await fetch("http://localhost:9999/api/upload", {
      method: 'POST',
      body: data,
    }).then(resp=>resp.text())
    console.log('---grpc upload---', path, hash);
    const message = await invoke(KoalaFS.upload,
      new UploadRequest().setPath(path).setHash(hash));
    console.log('---grpc upload cb---', message);
    return hash;
  } catch (e) {
    console.error('---grpc upload error---', e);
    error('上传文件', e.message);
  }
  return null;
}
