import {useEffect} from 'react';

import Menu from 'components/Menu/Menu';
import {Body, Content, Layout, Sider} from 'components/Web/Web';
import SystemConfig from 'pages/Setting/SystemConfig';
import Version from 'components/Version';

import useMenu from 'hox/menu';
import useSysConfig from 'hox/sysConfig';
import Files from "./pages/Files";
import {useColorScheme} from "@mui/material";

function App() {
    const {sysConfig} = useSysConfig();
    const {menu} = useMenu();
    const {mode, setMode} = useColorScheme();
    useEffect(() => {
        document.body.setAttribute('data-theme', sysConfig.theme);
        console.log("mode:", mode, "=>", sysConfig.theme);
        setMode(sysConfig.theme);
    }, [sysConfig.theme]);
    return (
        <div className="App">
            <Body>
                <Layout>
                    <Sider>
                        <Menu items={[
                            {icon: 'wangpan', name: '文件'},
                            {icon: 'tongbu', name: '传输列表'},
                            {icon: 'peizhi', name: '设置'},
                            {icon: 'system', name: '资源监控'},
                        ]}
                        />
                        <Version/>
                    </Sider>
                    <Content>
                        {menu === '文件' && <Files/>}
                        {menu === '传输列表' && <span>{menu}</span>}
                        {menu === '设置' && <SystemConfig/>}
                        {menu === '资源监控' && <span>{menu}</span>}
                    </Content>
                </Layout>
            </Body>
        </div>
    );
}

export default App;
