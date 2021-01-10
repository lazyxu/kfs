import React from 'react';
import styled from 'styled-components';

import FileIconClickable from 'components/FileIconClickable';
import FileNameClickable from 'components/FileNameClickable';

import {
  busState, setState, busValue, StoreContext,
} from 'bus/bus';
import { mv } from 'bus/fs';
import { warn } from 'bus/notification';
import { join } from 'utils/filepath';

const File = styled.div`
  width: 5em;
  height: 8em;
  margin: 0.5em;
  background-color: transparent;
  z-index: 1;
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

export default React.memo(({
  type, name, chosen, path, dir,
  atimems, mtimems, ctimems, birthtimems,
  onDoubleClick,
}) => {
  const [dragOver, setDragOver] = React.useState(false);
  const fileNameElm = React.useRef();
  const context = React.useContext(StoreContext);
  React.useMemo(() => {
    // componentWillReceiveProps
    if (fileNameElm.current) {
      const fileName = fileNameElm.current.value;
      const src = join(dir, name);
      const dst = join(dir, fileName);
      console.log('---rename---', dir, src, dst);
      mv([src], dst);
    }
  }, [chosen]);
  React.useEffect(() => {
    // componentDidUpdate
    if (chosen && fileNameElm.current) {
      fileNameElm.current.focus();
      fileNameElm.current.select();
    }
  });
  const isDargTargetValid = type === 'dir' && !busState.chosen[path];
  return (
    <File
      name="file"
      data-path={path}
      draggable="true"
      onDragStart={e => {
        const files = Object.keys(busState.chosen).filter((k) => busState.chosen[k] > 0);
        e.dataTransfer.setData('text/plain', JSON.stringify(files));
        console.log('onDragStart', name, e.dataTransfer.getData('text/plain'));
      }}
      onDragOver={e => {
        e.preventDefault();
        if (!dragOver && isDargTargetValid) {
          console.log('onDragEnter', e.dataTransfer.getData('text/plain'));
          setDragOver(true);
        }
      }}
      onDragLeave={e => {
        e.preventDefault();
        if (dragOver && isDargTargetValid) {
          console.log('onDragLeave', e.dataTransfer.getData('text/plain'));
          setDragOver(false);
        }
      }}
      onDrop={(e) => {
        e.preventDefault();
        console.log('onDrop', e.dataTransfer.getData('text/plain'));
        if (dragOver && isDargTargetValid) {
          let files = e.dataTransfer.getData('text/plain');
          files = JSON.parse(files);
          console.log('onDrop', files, path);
          if (files.includes(path)) {
            warn('移动文件至文件夹', '移动文件夹至本身');
          } else {
            mv(files, path);
          }
          setDragOver(false);
        }
      }}
      onContextMenu={(e) => {
        e.preventDefault();
        console.log('onContextMenu', e.target, e.target.getAttribute('data-tag'));
        if (e.target.getAttribute('data-tag') === 'choose-able') {
          let { clientX, clientY } = e;
          const { clientWidth, clientHeight } = context.state.fileListView;
          (clientX > clientWidth - 200) && (clientX = clientWidth - 200);
          (clientY > clientHeight - 200) && (clientY = clientHeight - 200);
          setState({
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
          context.setState({
            contextMenu: null,
          });
        }
      }}
    >
      <FileIconClickable
        path={path}
        xlinkHref={type === 'file' ? '#icon-file3' : '#icon-floderblue'}
        style={{
          backgroundColor: chosen ? '#343537' : 'transparent',
          border: `1px dashed ${dragOver ? 'white' : 'transparent'}`,
        }}
        onClick={e => {
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
        onDoubleClick={onDoubleClick}
      />
      <TextWrapper>
        {chosen === 2
          ? (
            <TextInEdit
              data-tag="choose-able"
              ref={fileNameElm}
              defaultValue={name}
              rows="3"
              onKeyPress={(e) => {
                if (e.which === 13) {
                  const fileName = e.target.value;
                  const src = join(dir, name);
                  const dst = join(dir, fileName);
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
              onDoubleClick={onDoubleClick}
            />
          )}
      </TextWrapper>
    </File>
  );
});
