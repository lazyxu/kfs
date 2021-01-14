import {
  Path, MoveRequest, PathList, UploadRequest,
} from 'pb/fs_pb';
import { KoalaFS } from 'pb/fs_pb_service';
import { error } from 'bus/notification';
import { dirname, basename, join } from 'utils/filepath';

import { grpc } from '@improbable-eng/grpc-web';
import { getConfig } from 'adaptor/config';
import Store from './store';

Store.prototype.invoke = async function (method, request, metadata) {
  return new Promise((resolve) => {
    grpc.invoke(method, {
      request,
      host: getConfig().host,
      metadata: Object.assign(metadata || {}, { 'kfs-mount': 'default' }),
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
};

Store.prototype.refresh = async function () {
  const message = await this.invoke(KoalaFS.ls,
    new Path().setPath(this.state.pwd));
  if (message.getFilesList) {
    const files = message.getFilesList().map((f) => f.toObject());
    if (files) {
      this.setState({ files });
    }
  }
};

Store.prototype.cd = async function (path) {
  try {
    path = path || this.state.pwd;
    console.log('---grpc cd---', path);
    const message = await this.invoke(KoalaFS.ls,
      new Path().setPath(path));
    console.log('---grpc cd cb---', message);
    if (message.getFilesList) {
      const files = message.getFilesList().map((f) => f.toObject());
      if (files) {
        this.setState({ files });
      }
    }
    this.setState({
      pwd: path,
      chosen: {},
    });
  } catch (e) {
    console.error('---grpc cd error---', e);
    error('读取文件夹', e.message);
  }
};

Store.prototype.mv = async function (srcList, dst) {
  try {
    if (srcList.length === 1 && srcList[0] === dst) {
      return;
    }
    console.log('---grpc mv---', srcList, dst);
    const message = await this.invoke(KoalaFS.mv,
      new MoveRequest().setSrcList(srcList).setDst(dst));
    console.log('---grpc mv cb---', message);
    await this.refresh();
    this.setState({
      chosen: (_chosen) => {
        Object.keys(_chosen).forEach((item) => {
          delete _chosen[item];
        });
        console.log(dst, dirname(dst), this.state.pwd);
        if (dirname(dst) === this.state.pwd) {
          const f = Object.values(this.state.files).find((f) => f.name === basename(dst));
          if (f && f.type === 'dir') {
            _chosen[dst] = 1;
            return;
          }
        }
        srcList.forEach((src) => {
          _chosen[`${this.state.pwd}/${basename(src)}`] = 1;
        });
      },
    });
  } catch (e) {
    console.error('---grpc mv error---', e);
    error('移动文件', e.message);
  }
};

Store.prototype.cp = async function (srcList, dst) {
  try {
    if (srcList.length === 1 && srcList[0] === dst) {
      return;
    }
    console.log('---grpc cp---', srcList, dst);
    const message = await this.invoke(KoalaFS.cp,
      new MoveRequest().setSrcList(srcList).setDst(dst));
    console.log('---grpc cp cb---', message.toObject());
    const { pathList } = message.toObject();
    this.refresh();
    this.setState({
      chosen: (_chosen) => {
        Object.keys(_chosen).forEach((item) => {
          delete _chosen[item];
        });
        pathList.forEach((p) => {
          _chosen[p] = 1;
        });
      },
    });
  } catch (e) {
    console.error('---grpc cp error---', e);
    error('复制文件', e.message);
  }
};

Store.prototype.newFile = async function (p) {
  try {
    console.log('---grpc newFile---', p);
    const message = await this.invoke(KoalaFS.newFile,
      new Path().setPath(p));
    console.log('---grpc newFile cb---', message);
    const { path } = message.toObject();
    this.refresh();
    this.setState({
      chosen: (_chosen) => {
        Object.keys(_chosen).forEach((item) => {
          delete _chosen[item];
        });
        _chosen[path] = 2;
      },
    });
  } catch (e) {
    console.error('---grpc newFile error---', e);
    error('新建文件', e.message);
  }
};

Store.prototype.newDir = async function (p) {
  try {
    console.log('---grpc newDir---');
    const message = await this.invoke(KoalaFS.newDir,
      new Path().setPath(p));
    console.log('---grpc newDir cb---', message);
    const { path } = message.toObject();
    this.refresh();
    this.setState({
      chosen: (_chosen) => {
        Object.keys(_chosen).forEach((item) => {
          delete _chosen[item];
        });
        _chosen[path] = 2;
      },
    });
  } catch (e) {
    console.error('---grpc newDir error---', e);
    error('新建文件夹', e.message);
  }
};

Store.prototype.remove = async function (path) {
  try {
    const pathList = typeof path === 'string' ? [path] : path;
    console.log('---grpc remove---', path);
    const message = await this.invoke(KoalaFS.remove,
      new PathList().setPathList(pathList));
    console.log('---grpc remove cb---', message);
    await this.cd(this.state.pwd);
    this.setState({
      chosen: (_chosen) => {
        pathList.forEach((path) => delete _chosen[path]);
      },
    });
  } catch (e) {
    console.error('---grpc mv error---', e);
    error('删除文件或文件夹', e.message);
  }
};

Store.prototype.download = async function (pathList) {
  try {
    console.log('---grpc download---', pathList);
    const message = await this.invoke(KoalaFS.download,
      new PathList().setPathList(pathList));
    console.log('---grpc download cb---', message);
    for (const hash of message.getHashList()) {
      const response = await fetch(`${getConfig().host}/api/download/${hash}`);
      if (!response.ok) {
        throw Error(await response.text());
      }
      const blob = await response.blob();
      const aTag = document.createElement('a');
      aTag.download = basename(pathList[0]);
      aTag.href = URL.createObjectURL(blob);
      aTag.click();
      URL.revokeObjectURL(blob);
    }
    // const hashList = message.getHashList();
    // const blockCount = hashList.length;
    // if (hashList && blockCount > 0) {
    //   const blobs = await map(hashList, async (hash, i) => {
    //     const hashHex = hash.map((x) => (`00${x.toString(16)}`).slice(-2)).join('');
    //     console.log(`download block ${i}/${blockCount}---`, hashHex);
    //     const message = await this.invoke(KoalaFS.download,
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
};

Store.prototype.upload = async function (path, data, hashList = []) {
  try {
    const hash = await fetch(`${getConfig().host}/api/upload`, {
      method: 'POST',
      body: data,
    }).then(resp => resp.text());
    console.log('---grpc upload---', path, hash, data.size);
    const message = await this.invoke(KoalaFS.upload,
      new UploadRequest().setPath(path).setHash(hash).setSize(data.size));
    console.log('---grpc upload cb---', message);
    this.refresh();
    return hash;
  } catch (e) {
    console.error('---grpc upload error---', e);
    error('上传文件', e.message);
  }
  return null;
};
