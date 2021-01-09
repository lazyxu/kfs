import React from 'react';

import styled from 'styled-components';
import Modal from 'components/Modal';

import { getConfig, setConfig } from 'adaptor/config';

const Textarea = styled.textarea`
  white-space: pre;
  word-break: break-all;
  position: relative;
  margin: 0;
  height: calc(100% - 2.5em);
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
  const [textarea, setTextarea] = React.useState({
    text: '',
    valid: false,
  });
  React.useEffect(() => {
    if (props.isOpen) {
      console.log('load config', JSON.stringify(getConfig(), undefined, 2));
      setTextarea({
        text: JSON.stringify(getConfig(), undefined, 2),
        valid: true,
      });
    }
  }, [props.isOpen]);
  return (
    <Modal
      title="系统配置"
      {...props}
    >
      <Textarea
        spellCheck="false"
        onChange={e => {
          const text = e.target.value;
          try {
            JSON.parse(text);
            setTextarea({
              text,
              valid: true,
            });
          } catch (e) {
            setTextarea({
              text,
              valid: false,
            });
            console.log(text, e);
          }
        }}
        value={textarea.text}
      />
      <Button
        disabled={!textarea.valid}
        type="button"
        onClick={() => {
          console.log(textarea);
          try {
            JSON.parse(textarea.text);
            setConfig(textarea.text);
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
