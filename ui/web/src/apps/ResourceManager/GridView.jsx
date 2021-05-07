import React from 'react';
import styled from 'styled-components';

import DefaultContextMenu from 'apps/ResourceManager/DefaultContextMenu';
import FileContextMenu from 'apps/ResourceManager/FileContextMenu';
import { BoxSelection, FileIconNameClickable } from 'kfs-components';

import {
  StoreContext, ctxInState,
} from 'bus/bus';
import { warn } from 'bus/notification';
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
  user-select: element;
  :focus {
    outline: none;
  }
`;

@ctxInState(StoreContext, 'files', 'chosen', 'boxChosen')
class component extends React.Component {
  static contextType = StoreContext;

  render() {
    const {
      files, chosen, boxChosen,
    } = this.state;
    return (
      <View
        onContextMenu={(e) => {
          e.preventDefault();
          const { fileListView } = this.context.state;
          if (e.target === fileListView || e.target.getAttribute('data-tag') !== 'choose-able') {
            const { clientX, clientY } = e;
            const { x, y } = fileListView.getBoundingClientRect();
            this.context.setState({
              contextMenuForFile: null,
              contextMenu: {
                x: Math.min(clientX, x + fileListView.clientWidth - 200),
                y: Math.min(clientY, y + fileListView.clientHeight - 120),
              },
            });
          }
        }}
        onMouseDown={(e) => {
          if (e.button === 2) {
            return;
          }
          if (e.target.getAttribute('data-tag') !== 'choose-able') {
            if (!e.metaKey) {
              this.context.setState({
                chosen: (_chosen) => {
                  Object.keys(_chosen).forEach((item) => {
                    delete _chosen[item];
                  });
                },
              });
            }
          }
        }}
        onMouseUp={(e) => {
          const keys = Object.keys(this.state.boxChosen);
          if (keys.length !== 0) {
            this.context.setState({
              chosen: (_chosen) => {
                keys.forEach((key) => _chosen[key] = 1);
              },
              boxChosen: {},
            });
          }
        }}
        tabIndex="-1"
        onKeyDown={(e) => {
          console.log(e.keyCode, e.metaKey);
          if (e.keyCode === 65 && e.metaKey === true) {
            this.context.setState({
              chosen: (_chosen) => {
                this.context.state.files.map((f) => join(this.context.state.pwd, f.name)).forEach((path) => {
                  _chosen[path] = 1;
                });
              },
            });
            e.preventDefault();
          }
        }}
        ref={(fileListView) => {
          this.context.setState({ fileListView });
        }}
      >
        <DefaultContextMenu />
        <FileContextMenu />
        <BoxSelection
          onPosChange={({
            x1, x2, y1, y2,
          }) => {
            const boxChosen = {};
            window.sss = this.context.state;
            const { childNodes, scrollTop } = this.context.state.fileListView;
            for (let i = 0; i < childNodes.length; i++) {
              const e = childNodes[i];
              const path = e.getAttribute('data-path');
              if (path) {
                if (!(e.offsetLeft + e.clientWidth < x1
                  || e.offsetTop + e.clientHeight - scrollTop < y1
                  || e.offsetLeft > x2
                  || e.offsetTop - scrollTop > y2)) {
                  boxChosen[path] = 1;
                }
              }
            }
            this.context.setState({ boxChosen });
          }}
        />
        {files.map((f) => {
          const {
            name,
            atimems, mtimems, ctimems, birthtimems,
          } = f;
          const { branch, pwd } = this.context.state;
          const path = join(pwd, name);
          return (
            <FileIconNameClickable
              key={`${f.type}-${path}`}
              type={f.type}
              name={f.name}
              chosen={chosen[path] || boxChosen[path]}
              path={path}
              onDoubleClick={() => f.type === 'dir' && this.context.cd(branch, path)}
              onRename={newName => {
                const src = join(pwd, name);
                const dst = join(pwd, newName);
                console.log('---rename---', pwd, src, dst);
                this.context.mv(branch, [src], branch, dst);
              }}
              onClickName={e => {
                this.context.setState({
                  contextMenuForFile: null,
                  contextMenu: null,
                  chosen: (_chosen) => {
                    const v = _chosen[path];
                    console.log('---FileNameClickable._chosen---', _chosen[path]);
                    if (chosen === 2 || !e.metaKey) {
                      Object.keys(_chosen).forEach((item) => {
                        delete _chosen[item];
                      });
                    }
                    Object.keys(_chosen)
                      .forEach((path) => _chosen[path] === 2 && (_chosen[path] = 1));
                    if (e.metaKey) {
                      _chosen[path] = v ? 0 : 1;
                    } else {
                      _chosen[path] = v ? 2 : 1;
                    }
                    console.log('---FileNameClickable._chosen.done---', _chosen[path]);
                  },
                });
              }}
              onEditNameComplete={newName => {
                const src = join(pwd, name);
                const dst = join(pwd, newName);
                this.context.setState({
                  chosen: (_chosen) => {
                    delete _chosen[src];
                    _chosen[dst] = 1;
                  },
                });
              }}
              onIconClick={e => {
                this.context.setState({
                  contextMenuForFile: null,
                  contextMenu: null,
                  atimems, mtimems, ctimems, birthtimems,
                  chosen: (_chosen) => {
                    const v = _chosen[path];
                    if (v === 1) {
                      const cnt = Object.values(_chosen)
                        .filter((v) => v > 0).reduce((a, b) => a + b, 0);
                      if (cnt !== 1) {
                        return {};
                      }
                    }
                    if (!e.metaKey) {
                      Object.keys(_chosen).forEach((item) => {
                        delete _chosen[item];
                      });
                    }
                    Object.keys(_chosen).forEach((path) => _chosen[path] === 2 && (_chosen[path] = 1));
                    _chosen[path] = 1;
                    return {};
                  },
                });
              }}
              onDrag={() => {
                const files = Object.keys(this.context.state.chosen).filter((k) => this.context.state.chosen[k] > 0);
                return JSON.stringify({ branch, files });
              }}
              onDrop={data => {
                const { files, branch } = JSON.parse(data);
                console.log('onDrop', files, path);
                if (files.includes(path)) {
                  warn('移动文件至文件夹', '移动文件夹至本身');
                } else {
                  this.context.mv(branch, files, branch, path);
                }
              }}
              onContextMenu={e => {
                const { fileListView } = this.context.state;
                console.log('onContextMenu', e.target, e.target.getAttribute('data-tag'));
                if (e.target.getAttribute('data-tag') === 'choose-able') {
                  const { clientX, clientY } = e;
                  const { x, y } = fileListView.getBoundingClientRect();
                  this.context.setState({
                    contextMenuForFile: {
                      x: Math.min(clientX, x + fileListView.clientWidth - 200),
                      y: Math.min(clientY, y + fileListView.clientHeight - 120),
                    },
                    contextMenu: null,
                  });
                }
              }}
            />
          );
        })}
      </View>
    );
  }
}

export default component;
