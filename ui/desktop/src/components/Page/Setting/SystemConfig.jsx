import useSysConfig from 'hox/sysConfig';
import { chooseDir } from 'remote/ui';
import styles from 'common/components/Page/Setting/index.module.scss';
import Status from 'common/components/Status/Status';

export default () => {
  const { sysConfig, setSysConfig, resetSysConfig } = useSysConfig();
  return (
    <div className={styles.page}>
      <header>设置</header>
      {!sysConfig ? <span>加载中...</span>
        : (
          <ul className={styles.configs}>
            <li className={styles.configs_item}>
              <div className={styles.configs_item_key} />
              <div className={styles.configs_item_values}>
                <button type="button" onClick={() => resetSysConfig()}>重置</button>
              </div>
            </li>
            <li className={styles.configs_item}>
              <div className={styles.configs_item_key}>
                <span>主题</span>
              </div>
              <div className={styles.configs_item_values}>
                <input type="radio" name="theme" checked={sysConfig.theme === 'light'} onChange={() => setSysConfig(c => ({ ...c, theme: 'light' }))} />
                <span> 浅色 </span>
                <input type="radio" name="theme" checked={sysConfig.theme === 'dark'} onChange={() => setSysConfig(c => ({ ...c, theme: 'dark' }))} />
                <span> 深色 </span>
                <input type="radio" name="theme" checked={sysConfig.theme === 'system'} onChange={() => setSysConfig(c => ({ ...c, theme: 'system' }))} />
                <span> 跟随系统 </span>
              </div>
            </li>
            <li className={styles.configs_item}>
              <div className={styles.configs_item_key}>
                <span>本地https服务端口</span>
              </div>
              <div className={styles.configs_item_values}>
                <input
                  type="text"
                  value={sysConfig?.backendProcess?.port}
                  onChange={e => setSysConfig(c => ({
                    ...c, backendProcess: {
                      ...sysConfig?.backendProcess, port: e.target.value,
                    },
                  }))}
                />
                <Status style={{ backgroundColor: sysConfig?.backendProcess?.status }} />
              </div>
            </li>
            <li className={styles.configs_item}>
              <div className={styles.configs_item_key}>
                <span>文件下载位置</span>
              </div>
              <div className={styles.configs_item_values}>
                <input type="text" value={sysConfig.downloadPath} disabled />
                <input
                  type="button"
                  value="更改"
                  onClick={e => {
                    e.preventDefault();
                    chooseDir().then(downloadPath => setSysConfig(c => ({ ...c, downloadPath })));
                  }}
                />
              </div>
            </li>
          </ul>
        )}
    </div>
  );
};
