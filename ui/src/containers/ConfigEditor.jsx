import React from 'react';

import styled from 'styled-components';
import Modal from 'components/Modal';

const Textarea = styled.textarea`
  white-space: pre;
  word-break: break-all;
  position: relative;
  margin: 0;
  height: calc(100% - 2em);
  width: 100%;
  padding: 0;
  background: var(--modal-body-color);
  overflow: scroll;
  color: white;
  :focus{
    outline: none;
  }
`;
const Button = styled.button`
  float: right;
  :hover{
    cursor: pointer;
  }
  margin: 3px;
  line-height: 100%;
`;

export default React.memo(({
  ...props
}) => {
  const bodyRef = React.useRef();
  const [text, setText] = React.useState(false);
  return (
    <Modal
      title="配置"
      save={() => {
        console.log(text);
        try {
          JSON.parse(text);
        } catch (e) {
          console.warn(e);
          return false;
        }
        return true;
      }}
      disableSave={!text}
      {...props}
    >
      <Textarea
        spellCheck="false"
        ref={bodyRef}
        onChange={e => {
          const text = e.target.value;
          try {
            JSON.parse(text);
            setText(text);
          } catch (e) {
            setText(false);
            console.log(text, e);
          }
        }}
        defaultValue="frehgeh"
      />
      <Button
        disabled={!text}
        type="button"
        onClick={() => {
          console.log(text);
          try {
            JSON.parse(text);
          } catch (e) {
            console.warn(e);
            return;
          }
          props.close();
        }}
      >
        保存
      </Button>
    </Modal>
  );
});
