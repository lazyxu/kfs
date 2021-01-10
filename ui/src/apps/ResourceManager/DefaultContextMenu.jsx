import React from 'react';

import ContextMenu from 'components/ContextMenu';

import { inState, busState, setState } from 'bus/bus';
import {
  cd, createFile, mkdir, mv, cp,
} from 'bus/fs';

@inState('contextMenu', 'cutFiles', 'copyFiles')
class component extends React.Component {
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
          新建文件: () => createFile(),
          新建文件夹: () => mkdir(),
          刷新: () => cd(),
          粘贴: {
            enabled: cutFiles.length > 0 || copyFiles.length > 0,
            fn: () => {
              if (cutFiles.length > 0) {
                mv(cutFiles, busState.pwd);
                setState({ cutFiles: [] });
              } else {
                cp(copyFiles, busState.pwd);
              }
            },
          },
          历史版本: { enabled: false, fn: () => { } },
        }}
        onFinish={() => setState({ contextMenu: null })}
      />
    );
  }
}

export default component;
