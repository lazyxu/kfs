import React from 'react';

import styled from 'styled-components';
import Modal from 'components/Modal';

import { getConfig, setConfig } from 'adaptor/config';
import rws, { invoke } from 'adaptor/ws';
import { openDir } from 'adaptor/backup';
import { StoreContext } from 'bus/bus';
import { getBranchList } from 'bus/grpcweb';

const Button = styled.button`
  list-style-type: none;
  padding-left: 1em;
`;

const stateParams = 0;
const stateReady = 1;
const stateBackup = 2;
const stateDone = 3;
const stateCancel = 4;

export default React.memo(({
  ...props
}) => {
  const [backupState, setBackupState] = React.useState();
  const [status, setStatus] = React.useState();
  const [uploadDir, setUploadDir] = React.useState();
  const [branch, setBranch] = React.useState();
  const [options, setOptions] = React.useState([]);
  const context = React.useContext(StoreContext);
  const branches = React.useRef();
  if (!branches.current) {
    // componentWillMount
    branches.current = true;
    getBranchList().then(list => {
      setOptions(list);
      const branch = list[0];
      console.log(branch);
    });
  }
  return (
    <Modal
      title="镜像备份"
      {...props}
    >
      <div>
        <button
          type="button"
          onClick={async () => {
            const path = await openDir();
            setUploadDir(path);
          }}
        >
          选择需要备份的文件夹
        </button>
      </div>
      <div>
        {uploadDir}
      </div>
      <div>
        <input
          type="text"
          name="greeting"
          list="greetings"
          onChange={e => {
            setBranch(e.target.value);
          }}
        />
        <datalist
          id="greetings"
          style={{ display: 'none' }}
        >
          {options.map(o => <option key={o} value={o}>{o}</option>)}
        </datalist>
      </div>
      {backupState === stateBackup
        ? (
          <Button
            onClick={() => {
              rws.reconnect(); // TODO: 更好的实现方式
              setBackupState(stateCancel); // TODO: 暂停备份
            }}
          >
            取消备份
          </Button>
        )
        : (
          <Button
            onClick={() => {
              invoke('backup', { path: uploadDir, branch }, ({ id, result }) => {
                setStatus(result);
                if (result.Done) {
                  setBackupState(stateDone);
                  return;
                }
                setBackupState(stateBackup);
              });
            }}
            disabled={!branch || !uploadDir}
          >
            开始备份
          </Button>
        )}
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
