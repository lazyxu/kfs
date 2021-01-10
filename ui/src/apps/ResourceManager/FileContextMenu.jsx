import React from 'react';

import ContextMenu from 'components/ContextMenu';

import {
  inState, setState,
} from 'bus/bus';
import {
  remove, download,
} from 'bus/fs';

@inState('contextMenuForFile', 'chosen')
class component extends React.Component {
  render() {
    const {
      contextMenuForFile, chosen,
    } = this.state;
    if (!contextMenuForFile) {
      return (
        <span />
      );
    }
    let contextMenuForFileOptions = {};
    const cnt = Object.values(chosen).filter((cnt) => cnt > 0).reduce((a, b) => a + b, 0);
    const pathList = Object.keys(chosen).filter((path) => chosen[path] > 0);
    if (cnt === 1) {
      contextMenuForFileOptions = {
        打开: { enabled: false, fn: () => { } },
        下载: () => download(pathList),
        剪切: () => setState({ cutFiles: pathList, copyFiles: [] }),
        复制: () => setState({ cutFiles: [], copyFiles: pathList }),
        删除: () => remove(pathList),
        重命名: () => setState({ chosen: (_chosen) => _chosen[pathList[0]] = 2 }),
        属性: { enabled: false, fn: () => { } },
        历史版本: { enabled: false, fn: () => { } },
      };
    } else {
      contextMenuForFileOptions = {
        打开: { enabled: false, fn: () => { } },
        [`下载${cnt}个文件`]: { enabled: false, fn: () => { } },
        剪切: () => setState({ cutFiles: pathList, copyFiles: [] }),
        复制: () => setState({ cutFiles: [], copyFiles: pathList }),
        删除: () => remove(pathList),
        属性: { enabled: false, fn: () => { } },
        历史版本: { enabled: false, fn: () => { } },
      };
    }
    return (
      <ContextMenu
        data-tag="contextMenu"
        x={contextMenuForFile.x}
        y={contextMenuForFile.y}
        options={contextMenuForFileOptions}
        onFinish={() => setState({ contextMenuForFile: null })}
      />
    );
  }
}

export default component;
