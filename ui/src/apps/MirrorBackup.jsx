import React from 'react';

import styled from 'styled-components';
import Modal from 'components/Modal';

import { getConfig, setConfig } from 'adaptor/config';
import { invoke } from 'adaptor/ws';

const Button = styled.button`
  list-style-type: none;
  padding-left: 1em;
`;

export default React.memo(({
  ...props
}) => {
  const [status, setStatus] = React.useState();
  return (
    <Modal
      title="镜像备份"
      {...props}
    >
      <Button onClick={() => {
        invoke('backup', { path: '/Users/xuliang/repos/kfs-network/kfscore' }, setStatus);
      }}
      >
        开始备份
      </Button>
      <div>
        {status && (
          <div>
            <div>文件总大小：{status.FileSize}</div>
            <div>文件数量：{status.FileCount}</div>
            <div>目录数量：{status.DirCount}</div>
            {/* <div>大文件大于100M）：{JSON.stringify(status.LargeFiles)}</div> */}
            <div>已排除文件：{JSON.stringify(status.IgnoredFiles)}</div>
            <div>是否取消：{status.Canceled ? '是' : '否'}</div>
            <div>扫描进度：{status.ScanProcess}</div>
            <div>待上传文件：{status.UploadingCount}</div>
            <div>备份完成：{status.Done ? '是' : '否'}</div>
          </div>
        )}
      </div>
    </Modal>
  );
});
