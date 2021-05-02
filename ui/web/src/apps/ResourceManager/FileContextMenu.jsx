import React from 'react';

import ContextMenu from 'components/ContextMenu';

import {
  ctxInState, setState, StoreContext,
} from 'bus/bus';

@ctxInState(StoreContext, 'contextMenuForFile', 'chosen')
class component extends React.Component {
  static contextType = StoreContext

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
        下载: () => this.context.download(pathList),
        剪切: () => setState({
          clipboard: {
            cut: true,
            file: {
              branch: this.context.state.branch,
              pathList,
            },
          },
        }),
        复制: () => setState({
          clipboard: {
            copy: true,
            file: {
              branch: this.context.state.branch,
              pathList,
            },
          },
        }),
        删除: () => this.context.remove(pathList),
        重命名: () => this.context.setState({ chosen: (_chosen) => _chosen[pathList[0]] = 2 }),
        属性: { enabled: false, fn: () => { } },
        历史版本: { enabled: false, fn: () => { } },
      };
    } else {
      contextMenuForFileOptions = {
        打开: { enabled: false, fn: () => { } },
        [`下载${cnt}个文件`]: { enabled: false, fn: () => { } },
        剪切: () => setState({
          clipboard: {
            cut: true,
            file: {
              branch: this.context.state.branch,
              pathList,
            },
          },
        }),
        复制: () => setState({
          clipboard: {
            copy: true,
            file: {
              branch: this.context.state.branch,
              pathList,
            },
          },
        }),
        删除: () => this.context.remove(pathList),
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
        onFinish={() => this.context.setState({ contextMenuForFile: null })}
      />
    );
  }
}

export default component;
