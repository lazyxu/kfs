import React from 'react';
import styled from 'styled-components';

import Dir from 'components/Dir';
import File from 'components/File';
import ContextMenu from 'components/ContextMenu';

import {
  inState, busState, setState, busValue,
} from 'bus/bus';
import {
  cd, createFile, mkdir, remove, download, mv, cp,
} from 'bus/fs';
import { join } from 'utils/filepath';

const View = styled.div`
  position: relative;
  height: 100%;
  width: 100%;
  overflow: scroll;
  display: flex;
  flex-flow:row wrap;
  align-content:flex-start;
  background-color: transparent;
`;

@inState('files', 'chosen', 'contextMenu', 'contextMenuForFile', 'boxChosen', 'cutFiles', 'copyFiles')
class component extends React.Component {
  render() {
    const {
      contextMenu, files, chosen, contextMenuForFile, boxSelection, boxChosen,
    } = this.state;
    let contextMenuForFileOptions = {};
    if (this.state.contextMenuForFile) {
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
          [`选中了${cnt}个文件`]: { enabled: false, fn: () => { } },
          打开: { enabled: false, fn: () => { } },
          [`下载${cnt}个文件`]: { enabled: false, fn: () => { } },
          剪切: () => setState({ cutFiles: pathList, copyFiles: [] }),
          复制: () => setState({ cutFiles: [], copyFiles: pathList }),
          删除: () => remove(pathList),
          属性: { enabled: false, fn: () => { } },
          历史版本: { enabled: false, fn: () => { } },
        };
      }
    }

    let left;
    let top;
    let height;
    let width;
    if (boxSelection) {
      if (boxSelection.x1 < boxSelection.x2) {
        left = boxSelection.x1;
        width = boxSelection.x2 - boxSelection.x1;
      } else {
        left = boxSelection.x2;
        width = boxSelection.x1 - boxSelection.x2;
      }
      if (boxSelection.y1 < boxSelection.y2) {
        top = boxSelection.y1;
        height = boxSelection.y2 - boxSelection.y1;
      } else {
        top = boxSelection.y2;
        height = boxSelection.y1 - boxSelection.y2;
      }
    }
    const BoxSelection = styled.div`
      position: fixed;
      left: ${typeof left !== 'undefined' ? `${left}px` : ''};
      top: ${typeof top !== 'undefined' ? `${top}px` : ''};
      width: ${typeof width !== 'undefined' ? `${width}px` : ''};
      height: ${typeof height !== 'undefined' ? `${height}px` : ''};
      border: 1px solid #5d5e60;
      background-color: #343537;
      opacity: .7;
      z-index: var(--z-body-select);
    `;
    return (
      <View
        onContextMenu={(e) => {
          e.preventDefault();
          const { fileListView } = busValue;
          if (e.target === fileListView || e.target.getAttribute('data-tag') !== 'choose-able') {
            let { clientX, clientY } = e;
            (clientX > fileListView.clientWidth - 200) && (clientX = fileListView.clientWidth - 200);
            (clientY > fileListView.clientHeight - 120) && (clientY = fileListView.clientHeight - 120);
            setState({
              contextMenu: { x: clientX, y: clientY },
              contextMenuForFile: null,
            });
          }
        }}
        onMouseDown={(e) => {
          if (e.button === 2) {
            return;
          }
          if (e.target.getAttribute('data-tag') !== 'choose-able') {
            if (!e.metaKey) {
              setState({
                chosen: (_chosen) => {
                  Object.keys(_chosen).forEach((item) => {
                    delete _chosen[item];
                  });
                },
              });
            }
            this.setState({ boxSelection: { x1: e.clientX, y1: e.clientY } });
          }
        }}
        onMouseMove={(e) => {
          if (this.state.boxSelection && !this.moving) {
            this.moving = true;
            setTimeout(() => {
              this.moving = false;
            }, 20);
            const { clientX: x2, clientY: y2 } = e;
            const { x1, y1 } = this.state.boxSelection;
            this.setState((prevState) => {
              prevState.boxSelection.x2 = x2;
              prevState.boxSelection.y2 = y2;
              return { boxSelection: prevState.boxSelection };
            });
            let left;
            let top;
            let right;
            let bottom;
            if (x1 < x2) {
              left = x1;
              right = x2;
            } else {
              left = x2;
              right = x1;
            }
            if (y1 < y2) {
              top = y1;
              bottom = y2;
            } else {
              top = y2;
              bottom = y1;
            }
            const boxChosen = {};
            const elements = document.getElementsByName('file');
            const { scrollTop } = busValue.fileListView;
            elements.forEach((e) => {
              if (!(e.offsetLeft + e.clientWidth < left
                || e.offsetTop + e.clientHeight - scrollTop < top
                || e.offsetLeft > right
                || e.offsetTop - scrollTop > bottom)) {
                const path = e.getAttribute('data-path');
                boxChosen[path] = 1;
              }
            });
            setState({ boxChosen });
          }
        }}
        onMouseUp={(e) => {
          const keys = Object.keys(this.state.boxChosen);
          if (keys.length !== 0) {
            setState({
              chosen: (_chosen) => {
                keys.forEach((key) => _chosen[key] = 1);
              },
              boxChosen: {},
            });
          }
          this.setState({ boxSelection: null });
        }}
        ref={(fileListView) => busValue.fileListView = fileListView}
      >
        {contextMenu && (
          <ContextMenu
            data-tag="contextMenu"
            x={contextMenu.x}
            y={contextMenu.y}
            options={{
              上传文件: console.log,
              新建文件: () => createFile(),
              新建文件夹: () => mkdir(),
              刷新: () => cd(busState.pwd),
              粘贴: {
                enabled: this.state.cutFiles.length > 0 || this.state.copyFiles.length > 0,
                fn: () => {
                  if (this.state.cutFiles.length > 0) {
                    mv(this.state.cutFiles, busState.pwd);
                    setState({ cutFiles: [] });
                  } else {
                    cp(this.state.copyFiles, busState.pwd);
                  }
                },
              },
              历史版本: { enabled: false, fn: () => { } },
            }}
            onFinish={() => setState({ contextMenu: null })}
          />
        )}
        <BoxSelection />
        {contextMenuForFile && (
          <ContextMenu
            x={contextMenuForFile.x}
            y={contextMenuForFile.y}
            options={contextMenuForFileOptions}
            onFinish={() => setState({ contextMenuForFile: null })}
          />
        )}
        {files.map((f) => {
          const path = join(busState.pwd, f.name);
          return f.type === 'file'
            ? <File key={`${f.type}-${path}`} {...f} chosen={chosen[path] || boxChosen[path]} dir={busState.pwd} path={path} />
            : <Dir key={`${f.type}-${path}`} {...f} chosen={chosen[path] || boxChosen[path]} dir={busState.pwd} path={path} />;
        })}
      </View>
    );
  }
}

export default component;
