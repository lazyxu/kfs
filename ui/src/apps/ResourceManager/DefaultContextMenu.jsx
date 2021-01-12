import React from 'react';

import ContextMenu from 'components/ContextMenu';

import {
  inState, ctxInState, StoreContext,
} from 'bus/bus';

@ctxInState(StoreContext, 'contextMenu', 'cutFiles', 'copyFiles')
class component extends React.Component {
  static contextType = StoreContext

  render() {
    const {
      contextMenu, copyFiles, cutFiles,
    } = this.state;
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
          新建文件: () => this.context.createFile(),
          新建文件夹: () => this.context.mkdir(),
          刷新: () => this.context.cd(),
          粘贴: {
            enabled: cutFiles.length > 0 || copyFiles.length > 0,
            fn: () => {
              if (cutFiles.length > 0) {
                this.context.mv(cutFiles, this.context.state.pwd);
                this.context.setState({ cutFiles: [] });
              } else {
                this.context.cp(copyFiles, this.context.state.pwd);
              }
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
