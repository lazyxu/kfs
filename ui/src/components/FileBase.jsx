import React from 'react';

import styled from 'styled-components';
import { busState, setState, busValue } from 'bus/bus';
import { mv } from 'bus/fs';
import { warn } from 'bus/notification';
import { join } from 'utils/filepath';

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
const Icon = styled.svg`
  height: 100%;
  width: 100%;
  vertical-align: -0.15em;
  fill: #dddddd;
  background-color: transparent;
`;
const TextInEdit = styled.textarea`
  font-size: 1em;
  width: 100%;
  margin: 0.3em 0;
  color: white;
  background-color: transparent;
`;

const fileImg = document.createElement('img');
fileImg.src = '/file.png';
fileImg.height = '3.7em';
fileImg.width = '3.7em';
const dirImg = document.createElement('img');
dirImg.src = '/dir.png';
dirImg.height = '3.7em';
dirImg.width = '3.7em';

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
    if (this.props.type === nextProps.type
      && this.props.name === nextProps.name
      && this.props.chosen === 2
      && nextProps.chosen !== 2) {
      const { name } = this.props;
      const fileName = this.fileNameElm.value;
      const { pwd } = busState;
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
    const { type, name, chosen } = this.props;
    const IconWrapper = styled.div`
      padding: 0.3em;
      margin: 0 0.3em;
      height: 4em;
      width: cal(100% - 0.5em);
      background-color: transparent;
      background-color: ${chosen && '#343537'};
      border: 1px dashed transparent;
      border-color: ${this.state.dragOver && 'white'};
      border-radius: 0.3em;
    `;
    const Text = styled.p`
      font-size: 1em;
      padding: 0;
      overflow : hidden;
      text-overflow: ellipsis;
      background-color: transparent;
      background-color: ${chosen === 1 && '#0e5ccd'};
      border-radius: 0.3em;
      display: -webkit-box;
      -webkit-line-clamp: 2;
      -webkit-box-orient: vertical;
      user-select: none;
      overflow-wrap: break-word;
      margin: 0;
    `;
    const { pwd } = busState;
    const path = join(pwd, name);
    return (
      <File
        name="file"
        data-path={path}
        draggable="true"
        onDragStart={(e) => {
          const path = join(busState.pwd, name);
          setState({
            chosen: (_chosen) => {
              _chosen[path] = 1;
            },
          });
          const files = Object.keys(busState.chosen).filter((k) => busState.chosen[k] > 0);
          e.dataTransfer.setData('text/plain', JSON.stringify(files));
          e.dataTransfer.setDragImage(type === 'file' ? fileImg : dirImg, 0, 0);
          console.log('onDragStart', name, e.dataTransfer.getData('text/plain'));
        }}
        onDragEnter={(e) => {
          e.preventDefault();
          if (type === 'dir' && !this.state.dragOver) {
            // console.log('onDragEnter', e.dataTransfer.getData('text/plain'));
            this.setState({ dragOver: true });
          }
        }}
        onDragOver={(e) => {
          e.preventDefault();
        }}
        onDragLeave={(e) => {
          e.preventDefault();
          if (type === 'dir' && this.state.dragOver) {
            // console.log('onDragLeave', e.dataTransfer.getData('text/plain'));
            this.setState({ dragOver: false });
          }
        }}
        onDrop={(e) => {
          if (type === 'dir') {
            let files = e.dataTransfer.getData('text/plain');
            files = JSON.parse(files);
            const { pwd } = busState;
            const dst = join(pwd, name);
            console.log('onDrop', files, dst);
            if (files.includes(dst)) {
              warn('移动文件至文件夹', '移动文件夹至本身');
            } else {
              mv(files, dst);
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
            (clientX > busValue.fileListView.clientWidth - 200) && (clientX = busValue.fileListView.clientWidth - 200);
            (clientY > busValue.fileListView.clientHeight - 200) && (clientY = busValue.fileListView.clientHeight - 200);
            const path = join(busState.pwd, name);
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
        <IconWrapper
          data-tag="choose-able"
          onMouseDown={(e) => {
            if (e.button === 2) {
              return;
            }
            if (this.onMouseDown()) {
              const path = join(busState.pwd, name);
              setState({
                chosen: (_chosen) => {
                  const v = _chosen[path];
                  if (v === 1) {
                    const cnt = Object.values(_chosen).filter((v) => v > 0).reduce((a, b) => a + b, 0);
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
                  if (v === 2) {
                    _chosen[path] = 1;
                  } else {
                    _chosen[path] = v ? 0 : 1;
                  }
                },
              });
            }
          }}
        >
          <Icon data-tag="choose-able" aria-hidden="true">
            <use data-tag="choose-able" xlinkHref={type === 'file' ? '#icon-file3' : '#icon-floderblue'} />
          </Icon>
        </IconWrapper>
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
                    const { pwd } = busState;
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
              <Text
                data-tag="choose-able"
                onMouseDown={(e) => {
                  if (e.button === 2) {
                    return;
                  }
                  if (this.onMouseDown()) {
                    setState({
                      chosen: (_chosen) => {
                        const path = join(busState.pwd, name);
                        const v = _chosen[path];
                        if (chosen === 2 || !e.metaKey) {
                          Object.keys(_chosen).forEach((item) => {
                            delete _chosen[item];
                          });
                        }
                        Object.keys(_chosen).forEach((path) => _chosen[path] === 2 && (_chosen[path] = 1));
                        if (e.metaKey) {
                          _chosen[path] = v ? 0 : 1;
                        } else {
                          _chosen[path] = v ? 2 : 1;
                        }
                      },
                    });
                  }
                }}
              >
                {name}
              </Text>
            )}
        </TextWrapper>
      </File>
    );
  }
}

export default component;
