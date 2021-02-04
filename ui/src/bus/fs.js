import {
  Path, MoveRequest, PathList, UploadRequest, PathReq,
} from 'pb/fs_pb';
import { KoalaFS } from 'pb/fs_pb_service';
import { error } from 'bus/notification';
import { dirname, basename, join } from 'utils/filepath';

import { grpc } from '@improbable-eng/grpc-web';
import { getConfig } from 'adaptor/config';
import Store from './store';

Store.prototype.invoke = async function (method, request) {
  return new Promise((resolve) => {
    grpc.invoke(method, {
      request,
      host: getConfig().host,
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
    new PathReq().setBranch(this.state.branch).setPath(this.state.pwd));
  if (message.getFilesList) {
    const files = message.getFilesList().map((f) => f.toObject());
    if (files) {
      this.setState({ files });
    }
  }
};

Store.prototype.cd = async function (branch, path) {
  try {
    const { pwd } = this.state;
    path = path || pwd;
    console.log('---grpc cd---', path);
    const message = await this.invoke(KoalaFS.ls,
      new PathReq().setBranch(branch).setPath(path));
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

Store.prototype.mv = async function (srcBranch, srcPathList, dstBranch, dstPath) {
  try {
    if (srcBranch === dstBranch && srcPathList.length === 1 && srcPathList[0] === dstPath) {
      return;
    }
    console.log('---grpc mv---', srcBranch, srcPathList, dstBranch, dstPath);
    const message = await this.invoke(KoalaFS.mv,
      new MoveRequest().setSrcbranch(srcBranch).setSrcpathList(srcPathList)
        .setDstbranch(dstBranch)
        .setDstpath(dstPath));
    console.log('---grpc mv cb---', message.toObject());
    await this.refresh();
    this.setState({
      chosen: (_chosen) => {
        Object.keys(_chosen).forEach((item) => {
          delete _chosen[item];
        });
        if (srcBranch === dstBranch) {
          console.log(dstPath, dirname(dstPath), this.state.pwd);
          if (dirname(dstPath) === this.state.pwd) {
            const f = Object.values(this.state.files).find((f) => f.name === basename(dstPath));
            if (f && f.type === 'dir') {
              _chosen[dstPath] = 1;
              return;
            }
          }
          srcPathList.forEach((src) => {
            _chosen[`${this.state.pwd}/${basename(src)}`] = 1;
          });
        }
      },
    });
  } catch (e) {
    console.error('---grpc mv error---', e);
    error('移动文件', e.message);
  }
};

Store.prototype.cp = async function (srcBranch, srcPathList, dstBranch, dstPath) {
  try {
    if (srcBranch === dstBranch && srcPathList.length === 1 && srcPathList[0] === dstPath) {
      return;
    }
    console.log('---grpc cp---', srcBranch, srcPathList, dstBranch, dstPath);
    const message = await this.invoke(KoalaFS.cp,
      new MoveRequest().setSrcbranch(srcBranch).setSrcpathList(srcPathList)
        .setDstbranch(dstBranch)
        .setDstpath(dstPath));
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

Store.prototype.newFile = async function () {
  try {
    const { branch, pwd } = this.state;
    console.log('---grpc newFile---', branch, pwd);
    const message = await this.invoke(KoalaFS.newFile,
      new PathReq().setBranch(branch).setPath(pwd));
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

Store.prototype.newDir = async function () {
  try {
    const { branch, pwd } = this.state;
    console.log('---grpc newDir---', branch, pwd);
    const message = await this.invoke(KoalaFS.newDir,
      new PathReq().setBranch(branch).setPath(pwd));
    console.log('---grpc newDir cb---', message.toObject());
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
    const { branch, pwd } = this.state;
    const pathList = typeof path === 'string' ? [path] : path;
    console.log('---grpc remove---', path);
    const message = await this.invoke(KoalaFS.remove,
      new PathList().setBranch(branch).setPathList(pathList));
    console.log('---grpc remove cb---', message);
    await this.cd(branch, pwd);
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
    const { branch } = this.state;
    console.log('---grpc download---', pathList);
    const message = await this.invoke(KoalaFS.download,
      new PathList().setBranch(branch).setPathList(pathList));
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
    const { branch } = this.state;
    const hash = await fetch(`${getConfig().host}/api/upload`, {
      method: 'POST',
      body: data,
    }).then(resp => resp.text());
    console.log('---grpc upload---', path, hash, data.size);
    const message = await this.invoke(KoalaFS.upload,
      new UploadRequest().setBranch(branch).setPath(path).setHash(hash)
        .setSize(data.size));
    console.log('---grpc upload cb---', message);
    this.refresh();
    return hash;
  } catch (e) {
    console.error('---grpc upload error---', e);
    error('上传文件', e.message);
  }
  return null;
};
