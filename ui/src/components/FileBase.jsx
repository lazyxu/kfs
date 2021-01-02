import React from 'react';

import styled from 'styled-components';
import { busState, setState, busValue } from 'bus/bus';
import { mv } from 'bus/fs';
import { warn } from 'bus/notification';
import { join } from 'utils/filepath';
import FileIconClickable from 'components/FileIconClickable';
import FileNameClickable from './FileNameClickable';

const File = styled.div`
  width: 5em;
  height: 8em;
  margin: 0.5em;
  background-color: transparent;
`;
const TextWrapper = styled.div`
  margin-top: 0.2em;
  height: 4em;
  width: 100%;
  text-align: center;
  background-color: transparent;
`;
const TextInEdit = styled.textarea`
  font-size: 1em;
  width: 100%;
  margin: 0.3em 0;
  color: white;
  background-color: transparent;
`;
class component extends React.Component {
  constructor(props) {
    super(props);
    this.clicked = false;
    this.state = {
      dragOver: false,
    };
  }

  // eslint-disable-next-line camelcase
  UNSAFE_componentWillReceiveProps(nextProps) {
    if (this.props.chosen === 2
      && nextProps.chosen !== 2) {
      const { name } = this.props;
      const fileName = this.fileNameElm.value;
      const pwd = this.props.dir;
      console.log('---rename---', pwd, join(pwd, name), join(pwd, fileName));
      const src = join(pwd, name);
      const dst = join(pwd, fileName);
      mv([src], dst);
    }
  }

  componentDidUpdate() {
    if (this.props.chosen && this.fileNameElm) {
      this.fileNameElm.focus();
      this.fileNameElm.select();
    }
  }

  componentWillUnmount() {
    console.log('componentWillUnmount');
  }

  onMouseDown() {
    console.log('---1---', this.clicked, this.props);
    const {
      atimems, mtimems, ctimems, birthtimems,
    } = this.props;
    setState({
      atimems, mtimems, ctimems, birthtimems,
    });
    if (this.props.onDoubleClick) {
      console.log('---2---', this.clicked);
      if (this.clicked) {
        console.log('---3---');
        this.props.onDoubleClick();
        this.clicked = false;
        return false;
      }
      this.clicked = true;
      setTimeout(() => {
        this.clicked = false;
      }, 200);
    }
    return true;
  }

  render() {
    const {
      type, name, chosen, path,
    } = this.props;
    console.log('---FileNameClickable._chosen.render---', chosen);
    return (
      <File
        name="file"
        data-path={path}
        draggable="true"
        onDragStart={(e) => {
          setState({
            chosen: (_chosen) => {
              _chosen[path] = 1;
            },
          });
          const files = Object.keys(busState.chosen).filter((k) => busState.chosen[k] > 0);
          e.dataTransfer.setData('text/plain', JSON.stringify(files));
          console.log('onDragStart', name, e.dataTransfer.getData('text/plain'));
        }}
        onDragEnter={(e) => {
          e.preventDefault();
          if (type === 'dir' && !this.state.dragOver && !busState.chosen[path]) {
            console.log('onDragEnter', e.dataTransfer.getData('text/plain'));
            this.setState({ dragOver: true });
          }
        }}
        onDragOver={(e) => {
          // console.log('onDragOver', e.dataTransfer.getData('text/plain'));
          if (type === 'dir' && !this.state.dragOver && !busState.chosen[path]) {
            this.setState({ dragOver: true });
          }
          e.preventDefault();
        }}
        onDragLeave={(e) => {
          e.preventDefault();
          if (type === 'dir' && this.state.dragOver) {
            console.log('onDragLeave', e.dataTransfer.getData('text/plain'));
            this.setState({ dragOver: false });
          }
        }}
        onDrop={(e) => {
          console.log('onDrop', e.dataTransfer.getData('text/plain'));
          if (type === 'dir' && !busState.chosen[path]) {
            let files = e.dataTransfer.getData('text/plain');
            files = JSON.parse(files);
            console.log('onDrop', files, path);
            if (files.includes(path)) {
              warn('移动文件至文件夹', '移动文件夹至本身');
            } else {
              mv(files, path);
            }
            if (type === 'dir' && this.state.dragOver) {
              this.setState({ dragOver: false });
            }
          }
          e.preventDefault();
        }}
        onContextMenu={(e) => {
          e.preventDefault();
          if (e.target.getAttribute('data-tag') === 'choose-able') {
            let { clientX, clientY } = e;
            (clientX > busValue.fileListView.clientWidth - 200)
              && (clientX = busValue.fileListView.clientWidth - 200);
            (clientY > busValue.fileListView.clientHeight - 200)
              && (clientY = busValue.fileListView.clientHeight - 200);
            setState({
              contextMenu: null,
              contextMenuForFile: { x: clientX, y: clientY },
              chosen: (_chosen) => {
                if (!chosen) {
                  Object.keys(_chosen).forEach((item) => {
                    delete _chosen[item];
                  });
                }
                _chosen[path] = 1;
              },
            });
          }
        }}
      >
        <FileIconClickable
          path={path}
          xlinkHref={type === 'file' ? '#icon-file3' : '#icon-floderblue'}
          style={{
            backgroundColor: chosen ? '#343537' : 'transparent',
            border: `1px dashed ${this.state.dragOver ? 'white' : 'transparent'}`,
          }}
          onClick={e => {
            console.log('---onClick---', this.clicked, this.props);
            const {
              atimems, mtimems, ctimems, birthtimems,
            } = this.props;
            setState({
              atimems, mtimems, ctimems, birthtimems,
            });
            setState({
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
                Object.keys(_chosen)
                  .forEach((path) => _chosen[path] === 2 && (_chosen[path] = 1));
                if (v === 2) {
                  _chosen[path] = 1;
                } else {
                  _chosen[path] = v ? 0 : 1;
                }
                return {};
              },
            });
          }}
          onDoubleClick={e => {
            console.log('---onDoubleClick---', this.clicked, this.props);
            this.props.onDoubleClick && this.props.onDoubleClick();
          }}
        />
        <TextWrapper>
          {chosen === 2
            ? (
              <TextInEdit
                data-tag="choose-able"
                ref={(fileNameElm) => this.fileNameElm = fileNameElm}
                defaultValue={name}
                rows="3"
                onKeyPress={(e) => {
                  if (e.which === 13) {
                    const fileName = e.target.value;
                    const pwd = this.props.dir;
                    const src = join(pwd, name);
                    const dst = join(pwd, fileName);
                    setState({
                      chosen: (_chosen) => {
                        delete _chosen[src];
                        _chosen[dst] = 1;
                      },
                    });
                  }
                  return true;
                }}
              />
            )
            : (
              <FileNameClickable
                name={name}
                style={{ backgroundColor: chosen === 1 ? '#0e5ccd' : 'transparent' }}
                onClick={e => {
                  console.log('---FileNameClickable.onClick---', this.clicked, this.props);
                  setState({
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
                onDoubleClick={e => {
                  console.log('---FileNameClickable.onDoubleClick---', this.clicked, this.props);
                  this.props.onDoubleClick && this.props.onDoubleClick();
                }}
              />
            )}
        </TextWrapper>
      </File>
    );
  }
}

export default component;
