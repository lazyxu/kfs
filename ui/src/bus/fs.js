import {
  PathRequest, MoveRequest, PathListRequest, DownloadRequest, UploadRequest,
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
      metadata: Object.assign(metadata || {}, { 'kfs-pwd': this.state.pwd, 'kfs-mount': 'default' }),
      onHeaders: (headers) => {
        // console.log(headers);
      },
      onMessage: (message) => {
        if (message.getFilesList) {
          const files = message.getFilesList().map((f) => f.toObject());
          if (files) {
            this.setState({ files });
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
};

Store.prototype.cd = async function (path) {
  try {
    path = path || this.state.pwd;
    console.log('---grpc cd---', path);
    const message = await this.invoke(KoalaFS.ls,
      new PathRequest().setPath(path));
    console.log('---grpc cd cb---', message);
    const { path: pwd } = message.toObject();
    if (pwd !== this.state.pwd) {
      this.setState({ chosenFiles: {} });
    }
    this.setState({
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
    if (srcList.length !== 1 || dirname(srcList[0]) === dirname(dst)) {
      await this.cd(this.state.pwd);
    }
    this.setState({
      chosen: (_chosen) => {
        Object.keys(_chosen).forEach((item) => {
          delete _chosen[item];
        });
        console.log(dst, dirname(dst), this.state.pwd);
        if (dirname(dst) === this.state.pwd
          && Object.values(this.state.files).find((f) => f.name === basename(dst)).type === 'dir') {
          _chosen[dst] = 1;
          return;
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
    console.log('---grpc cp cb---', message);
    if (srcList.length !== 1 || dirname(srcList[0]) === dirname(dst)) {
      await this.cd(this.state.pwd);
    }
    this.setState({
      chosen: (_chosen) => {
        Object.keys(_chosen).forEach((item) => {
          delete _chosen[item];
        });
        srcList.forEach((src) => {
          _chosen[`${this.state.pwd}/${basename(src)}`] = 1;
        });
      },
    });
  } catch (e) {
    console.error('---grpc cp error---', e);
    error('复制文件', e.message);
  }
};

Store.prototype.createFile = async function (path) {
  try {
    console.log('---grpc createFile---', path);
    const message = await this.invoke(KoalaFS.createFile,
      new PathRequest().setPath(path));
    console.log('---grpc createFile cb---', message);
    const { name } = message.toObject();
    this.setState({
      chosen: (_chosen) => {
        Object.keys(_chosen).forEach((item) => {
          delete _chosen[item];
        });
        _chosen[join(this.state.pwd, name)] = 2;
      },
    });
  } catch (e) {
    console.error('---grpc createFile error---', e);
    error('新建文件', e.message);
  }
};

Store.prototype.mkdir = async function (path) {
  try {
    console.log('---grpc mkdir---');
    const { pwd } = this.state;
    const message = await this.invoke(KoalaFS.mkdir,
      new PathRequest().setPath(path));
    console.log('---grpc mkdir cb---', message);
    const { name } = message.toObject();
    this.setState({
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
};

Store.prototype.remove = async function (path) {
  try {
    const pathList = typeof path === 'string' ? [path] : path;
    console.log('---grpc remove---', path);
    const message = await this.invoke(KoalaFS.remove,
      new PathListRequest().setPathList(pathList));
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
      new DownloadRequest().setPathList(pathList));
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
    return hash;
  } catch (e) {
    console.error('---grpc upload error---', e);
    error('上传文件', e.message);
  }
  return null;
};
