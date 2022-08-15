import { useEffect } from 'react';

import Menu from 'components/Menu/Menu';
import { Body, Layout, Sider, Content } from 'components/Web/Web';
import SystemConfig from 'components/Page/Setting/SystemConfig';
import Version from 'components/Version';

import useMenu from 'hox/menu';
import useSysConfig from 'hox/sysConfig';

function App() {
  const { sysConfig } = useSysConfig();
  const { menu } = useMenu();
  useEffect(() => {
    document.body.setAttribute('data-theme', sysConfig.theme);
  }, [sysConfig.theme]);
  return (
    <div className="App">
      <Body>
        <Layout>
          <Sider>
            <Menu items={[
              { icon: 'wangpan', name: '文件' },
              { icon: 'tongbu', name: '传输列表' },
              { icon: 'peizhi', name: '设置' },
              { icon: 'system', name: '资源监控' },
            ]}
            />
            <Version />
          </Sider>
          <Content>
            {menu === '文件' && <span>{menu}</span>}
            {menu === '传输列表' && <span>{menu}</span>}
            {menu === '设置' && <SystemConfig />}
            {menu === '资源监控' && <span>{menu}</span>}
          </Content>
        </Layout>
      </Body>
    </div>
  );
}

export default App;
