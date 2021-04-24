import React from 'react';

import ContextMenu from 'components/ContextMenu';

import {
  inState, ctxInState, StoreContext,
} from 'bus/bus';

@inState('clipboard')
@ctxInState(StoreContext, 'contextMenu')
class component extends React.Component {
  static contextType = StoreContext

  render() {
    const {
      contextMenu, clipboard,
    } = this.state;
    const {
      branch, pwd,
    } = this.context.state;
    if (!contextMenu) {
      return (
        <span />
      );
    }
    return (
      <ContextMenu
        data-tag="contextMenu"
        x={contextMenu.x}
        y={contextMenu.y}
        options={{
          上传文件: console.log,
          新建文件: () => this.context.newFile(),
          新建文件夹: () => this.context.newDir(),
          刷新: () => this.context.cd(branch, pwd),
          粘贴: {
            enabled: clipboard && clipboard.file,
            fn: () => {
              const { branch: srcBranch, pathList } = clipboard.file;
              if (clipboard.cut) {
                this.context.mv(srcBranch, pathList, branch, pwd);
                this.context.setState({ clipboard: undefined });
                return;
              }
              this.context.cp(srcBranch, pathList, branch, pwd);
              this.context.setState({ clipboard: undefined });
            },
          },
          历史版本: { enabled: false, fn: () => { } },
        }}
        onFinish={() => this.context.setState({ contextMenu: null })}
      />
    );
  }
}

export default component;
