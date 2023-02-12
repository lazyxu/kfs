import {useEffect, useMemo} from 'react';

import Menu from 'components/Menu/Menu';
import {Body, Content, Layout, Sider} from 'components/Web/Web';
import SystemConfig from 'pages/Setting/SystemConfig';
import Version from 'components/Version';

import useMenu from 'hox/menu';
import useSysConfig from 'hox/sysConfig';
import Files from "./pages/Files";
import {createTheme, ThemeProvider} from "@mui/material";

function App() {
    const {sysConfig} = useSysConfig();
    const {menu} = useMenu();
    useEffect(() => {
        document.body.setAttribute('data-theme', sysConfig.theme);
    }, [sysConfig.theme]);
    const theme = useMemo(
        () =>
            createTheme({
                palette: {
                    mode: sysConfig.theme,
                },
            }),
        [sysConfig.theme],
    );
    return (
        <ThemeProvider theme={theme}>
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
        </ThemeProvider>
    );
}

export default App;
