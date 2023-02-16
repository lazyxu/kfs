import Menu from 'components/Menu/Menu';
import {Body, Content, Layout, Sider} from 'components/Web/Web';
import SystemConfig from 'pages/Setting/SystemConfig';
import Version from 'components/Version';
import Files from "./pages/Files";
import {useColorScheme} from "@mui/material";
import React, {useEffect} from "react";
import useMenu from "./hox/menu";
import useSysConfig from "./hox/sysConfig";

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
                            {icon: 'wangpan', name: '我的文件'},
                            {icon: 'tongbu', name: '传输列表'},
                            {icon: 'peizhi', name: '设置'},
                            {icon: 'system', name: '资源监控'},
                            {icon: '', name: '我的书签'},
                            {icon: '', name: '分享历史'},
                            {icon: '', name: '操作历史'},
                            {icon: '', name: '搜索'},
                        ]}
                        />
                        <Version/>
                    </Sider>
                    <Content>
                        {menu === '我的文件' && <Files/>}
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
