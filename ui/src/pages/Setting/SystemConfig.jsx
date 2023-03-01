import useSysConfig from 'hox/sysConfig';
import styles from 'pages/Setting/index.module.scss';

export default () => {
    const {sysConfig, setSysConfig, resetSysConfig} = useSysConfig();
    return (
        <div className={styles.page}>
            <header>设置</header>
            {!sysConfig ? <span>加载中...</span>
                : (
                    <ul className={styles.configs}>
                        <li className={styles.configs_item}>
                            <div className={styles.configs_item_key}/>
                            <div className={styles.configs_item_values}>
                                <button type="button" onClick={() => resetSysConfig()}>重置</button>
                            </div>
                        </li>
                        <li className={styles.configs_item}>
                            <div className={styles.configs_item_key}>
                                <span>主题</span>
                            </div>
                            <div className={styles.configs_item_values}>
                                <input type="radio" name="theme" checked={sysConfig.theme === 'light'}
                                       onChange={() => setSysConfig(c => ({...c, theme: 'light'}))}/>
                                <span> 浅色 </span>
                                <input type="radio" name="theme" checked={sysConfig.theme === 'dark'}
                                       onChange={() => setSysConfig(c => ({...c, theme: 'dark'}))}/>
                                <span> 深色 </span>
                                {process.env.REACT_APP_PLATFORM === 'web' ? [] : <>
                                    <input type="radio" name="theme" checked={sysConfig.theme === 'system'}
                                           onChange={() => setSysConfig(c => ({...c, theme: 'system'}))} disabled/>
                                    <span> 跟随系统 </span>
                                </>
                                }
                            </div>
                        </li>
                        <li className={styles.configs_item}>
                            <div className={styles.configs_item_key}>
                                <span>API</span>
                            </div>
                            <div className={styles.configs_item_values}>
                                <input type="radio" name="api" checked={sysConfig.api === 'mock'}
                                       onChange={() => setSysConfig(c => ({...c, api: 'mock'}))}/>
                                <span> mock </span>
                                <input type="radio" name="api" checked={sysConfig.api === 'web'}
                                       onChange={() => setSysConfig(c => ({...c, api: 'web'}))}/>
                                <span> web </span>
                            </div>
                        </li>
                        {process.env.NODE_ENV === 'production' ? [] :
                            <li className={styles.configs_item}>
                                <div className={styles.configs_item_key}>
                                    <span>Web服务器</span>
                                </div>
                                <div className={styles.configs_item_values}>
                                    <input type="text" value={sysConfig.webServer}
                                           onChange={e => setSysConfig(c => ({...c, webServer: e.target.value}))}/>
                                </div>
                            </li>
                        }
                        <li className={styles.configs_item}>
                            <div className={styles.configs_item_key}>
                                <span>Socket服务器</span>
                            </div>
                            <div className={styles.configs_item_values}>
                                <input type="text" value={sysConfig.socketServer}
                                       onChange={e => setSysConfig(c => ({...c, socketServer: e.target.value}))}/>
                            </div>
                        </li>
                        <li className={styles.configs_item}>
                            <div className={styles.configs_item_key}>
                                <span>客户端端口</span>
                            </div>
                            <div className={styles.configs_item_values}>
                                <input type="text" value={sysConfig.port}/>
                            </div>
                        </li>
                    </ul>
                )}
        </div>
    );
};
