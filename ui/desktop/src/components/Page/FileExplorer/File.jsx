import { useState, useEffect } from 'react';
import Button from '@mui/material/Button';
import ButtonGroup from '@mui/material/ButtonGroup';
import TextField from '@mui/material/TextField';
import Dialog from '@mui/material/Dialog';
import DialogActions from '@mui/material/DialogActions';
import DialogContent from '@mui/material/DialogContent';
import DialogTitle from '@mui/material/DialogTitle';
import Icon from 'common/components/Icon/Icon';
import useNotification from 'common/components/Notification/Notification';
import useSysConfig from 'hox/sysConfig';
import useFilePath from 'hox/filepath';
import { Menu, MenuItem } from 'remote/menu';
import { deleteBranch, getBranchHash, renameBranch } from 'remote/branch';
import { readDirectory, createDirectory } from 'remote/object';

import styles from './index.module.scss';

export default () => {
  const [_, __, sendError] = useNotification();
  const [branches, setBranches] = useState([]);
  const [showCreateDirectory, setShowCreateDirectory] = useState(false);
  const [createDirectoryName, setCreateDirectoryName] = useState('');
  const [showRenameBranch, setShowRenameBranch] = useState(undefined);
  const [renameBranchName, setRenameBranchName] = useState('');
  const { sysConfig } = useSysConfig();
  const [filepath, setFilepath] = useFilePath();
  const refresh = () => {
    getBranchHash(sysConfig.clientID, filepath).then(readDirectory).then(data => {
      console.log(typeof data, data);
    }).catch(sendError);
  };
  useEffect(refresh, []);
  return (
    <div
      onContextMenu={e => {
        const menu = new Menu();
        menu.append(new MenuItem({
          label: '新建文件夹',
          click: () => setShowCreateDirectory(true),
        }));
        menu.append(new MenuItem({
          label: '上传',
          click: () => setShowCreateDirectory(true),
          enabled: false,
        }));
        menu.append(new MenuItem({
          label: '刷新',
          click: refresh,
        }));
        menu.popup();
      }}
    >
      <div>{filepath}</div>
      <ul>
        {branches.map(branch => (
          <li
            key={branch.branchName}
            onContextMenu={e => {
              const menu = new Menu();
              menu.append(new MenuItem({
                label: '打开',
                click: () => setFilepath(branch.branchName),
              }));
              menu.append(new MenuItem({
                label: '重命名',
                click: () => {
                  setRenameBranchName(branch.branchName);
                  setShowRenameBranch(branch.branchName);
                },
              }));
              menu.append(new MenuItem({
                label: '拷贝分支',
                enabled: false,
              }));
              menu.append(new MenuItem({
                label: '删除',
                click: () => {
                  deleteBranch(sysConfig.clientID, branch.branchName).then(data => {
                    setShowCreateDirectory(false);
                    refresh();
                  }).catch(sendError);
                },
              }));
              menu.append(new MenuItem({
                label: '属性',
                click: () => {
                  console.log('item 1 clicked');
                },
                enabled: false,
              }));
              menu.popup();
            }}
          >
            <div><Icon icon="wangpan" className={styles.icon} /></div>
            <span className={styles.text}>{branch.branchName}</span>
          </li>
        ))}
      </ul>
      <Dialog open={showCreateDirectory} onClose={() => setShowCreateDirectory(false)}>
        <DialogTitle>新建文件夹</DialogTitle>
        <DialogContent>
          <TextField
            autoFocus
            margin="dense"
            id="name"
            label="文件夹名称"
            type="email"
            fullWidth
            variant="standard"
            value={createDirectoryName}
            onChange={e => setCreateDirectoryName(e.target.value)}
          />
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setShowCreateDirectory(false)}>取消</Button>
          <Button onClick={() => {
            createDirectory(sysConfig.clientID, filepath, createDirectoryName).then(data => {
              setShowCreateDirectory(false);
              refresh();
            }).catch(sendError);
          }}
          >
            确定
          </Button>
        </DialogActions>
      </Dialog>
      <Dialog open={!!showRenameBranch} onClose={() => setRenameBranchName(undefined)}>
        <DialogTitle>重命名分支</DialogTitle>
        <DialogContent>
          <TextField
            autoFocus
            margin="dense"
            id="name"
            label="分支名称"
            type="email"
            fullWidth
            variant="standard"
            value={renameBranchName}
            onChange={e => setRenameBranchName(e.target.value)}
          />
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setShowRenameBranch(undefined)}>取消</Button>
          <Button onClick={() => {
            renameBranch(sysConfig.clientID, showRenameBranch, renameBranchName).then(data => {
              setShowRenameBranch(undefined);
              refresh();
            }).catch(sendError);
          }}
          >
            确定
          </Button>
        </DialogActions>
      </Dialog>
    </div>
  );
};
