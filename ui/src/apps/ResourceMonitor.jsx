import React from 'react';

import styled from 'styled-components';
import Modal from 'components/Modal';

import { getConfig, setConfig } from 'adaptor/config';
import { invoke } from 'adaptor/ws';
import { getStatus } from 'bus/grpcweb';

const Button = styled.button`
  list-style-type: none;
  padding-left: 1em;
`;

export default React.memo(({
  ...props
}) => {
  const [status, setStatus] = React.useState();
  const context = React.useRef();
  if (!context.current) {
    // componentWillMount
    context.current = true;
    getStatus().then(setStatus);
  }
  return (
    <Modal
      title="资源管理"
      {...props}
    >
      <div>
        {status ? (
          <div>
            <div>占用空间总大小：{status.totalsize}</div>
            <div>文件总大小：{status.filesize}</div>
            <div>文件数量：{status.filecount}</div>
            <div>目录数量：{status.dircount}</div>
          </div>
        ) : (
          <div>统计中...</div>
        )}
      </div>
    </Modal>
  );
});
