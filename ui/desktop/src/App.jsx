import { useEffect } from 'react';

import Menu from 'common/components/Menu/Menu';
import { Body, Layout, Sider, Content } from 'common/components/Web/Web';
import DragableArea from 'components/DragableArea/DragableArea';
import SystemConfig from 'components/Page/Setting/SystemConfig';
import CloudConnection from 'components/Page/CloudConnection/CloudConnection';

import useMenu from 'common/hox/menu';
import useSysConfig from 'hox/sysConfig';
import { backendProcess } from 'remote/backendProcess';
import { getBackendInstance } from 'remote/axios';
import FileExplorer from 'components/Page/FileExplorer/FileExplorer';

function App() {
  const { sysConfig, setSysConfig } = useSysConfig();
  const { menu } = useMenu();
  useEffect(() => {
    document.body.setAttribute('data-theme', sysConfig.theme);
  }, [sysConfig.theme]);
  useEffect(() => {
    backendProcess(sysConfig?.backendProcess?.port).then(() => setSysConfig(c => ({
      ...c, backendProcess: {
        ...sysConfig?.backendProcess, status: 'green',
      },
    }))).catch(() => setSysConfig(c => ({
      ...c, backendProcess: {
        ...sysConfig?.backendProcess, status: 'red',
      },
    })));
  }, [sysConfig?.backendProcess?.port]);
  // TODO: listen to remote config changes
  useEffect(() => {
    setSysConfig(c => {
      c.remotes.forEach(remote => {
        remote.status = 'yellow';
      });
      return { ...c };
    });
    getBackendInstance().get('/api/connect').then(res => {
      res.data.map((status, i) => setSysConfig(c => {
        c.remotes[i].status = status;
        return { ...c };
      }));
    }).catch(() => {
      setSysConfig(c => {
        c.remotes.forEach(remote => {
          remote.status = 'red';
        });
        return { ...c };
      });
    });
  }, [sysConfig?.remotes]);
  return (
    <div className="App">
      <DragableArea />
      <Body>
        <Layout>
          <Sider>
            <Menu items={[
              { icon: 'wangpan', name: '云盘' },
              { icon: 'tongbu', name: '传输列表' },
              { icon: 'peizhi', name: '设置' },
              { icon: 'cloud-connection', name: '远程连接' },
              { icon: 'system', name: '资源监控' },
            ]}
            />
          </Sider>
          <Content>
            {menu === '云盘' && <FileExplorer />}
            {menu === '传输列表' && <span>{menu}</span>}
            {menu === '设置' && <SystemConfig />}
            {menu === '远程连接' && <CloudConnection />}
            {menu === '资源监控' && <span>{menu}</span>}
          </Content>
        </Layout>
      </Body>
    </div>
  );
}

export default App;
