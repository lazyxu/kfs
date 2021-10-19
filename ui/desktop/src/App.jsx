import { useEffect } from 'react';

import Menu from 'common/components/Menu/Menu';
import { Body, Layout, Sider, Content } from 'common/components/Web/Web';
import DragableArea from 'components/DragableArea/DragableArea';
import SystemConfig from 'components/Page/Setting/SystemConfig';

import useMenu from 'common/hox/menu';
import useSysConfig from 'hox/sysConfig';

function App() {
  const { sysConfig } = useSysConfig();
  const { menu } = useMenu();
  useEffect(() => {
    document.body.setAttribute('data-theme', sysConfig.theme);
  }, [sysConfig.theme]);
  return (
    <div className="App">
      <DragableArea />
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
