import React from 'react';
import styled from 'styled-components';

import { FileIconClickable, FileNameClickable } from 'kfs-components';

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
  type, name, chosen, path,
  onDoubleClick, onRename, onClickName, onEditNameComplete, onIconClick, onDrag, onDrop, onContextMenu,
}) => {
  const [dragOver, setDragOver] = React.useState(false);
  const fileNameElm = React.useRef();
  React.useMemo(() => {
    // componentWillReceiveProps
    if (fileNameElm.current) {
      const fileName = fileNameElm.current.value;
      onRename(fileName);
    }
  }, [chosen]);
  // React.useEffect(() => {
  //   // componentDidUpdate
  //   if (chosen && fileNameElm.current) {
  //     console.log('componentDidUpdate', chosen, fileNameElm.current);
  //     fileNameElm.current.focus();
  //     fileNameElm.current.select();
  //   }
  // });
  const isDargTargetValid = type === 'dir' && !chosen;
  return (
    <File
      name="file"
      data-path={path}
      draggable="true"
      onDragStart={e => {
        e.dataTransfer.setData('text/plain', onDrag());
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
          onDrop(e.dataTransfer.getData('text/plain'));
          setDragOver(false);
        }
      }}
      onContextMenu={(e) => {
        e.preventDefault();
        onContextMenu(e);
      }}
    >
      <FileIconClickable
        path={path}
        icon={type === 'file' ? 'file3' : 'floderblue'}
        style={{
          backgroundColor: chosen ? '#343537' : 'transparent',
          border: `1px dashed ${dragOver ? 'white' : 'transparent'}`,
        }}
        onClick={onIconClick}
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
                  onEditNameComplete(e.target.value);
                }
                return true;
              }}
            />
          )
          : (
            <FileNameClickable
              name={name}
              style={{ backgroundColor: chosen === 1 ? '#0e5ccd' : 'transparent' }}
              onClick={onClickName}
              onDoubleClick={onDoubleClick}
            />
          )}
      </TextWrapper>
    </File>
  );
});
