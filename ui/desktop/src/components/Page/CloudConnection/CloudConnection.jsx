import useSysConfig from 'hox/sysConfig';
import styles from 'common/components/Page/Setting/index.module.scss';
import Status from 'common/components/Status/Status';

export default () => {
  const { sysConfig, setSysConfig, resetSysConfig } = useSysConfig();
  return (
    <div className={styles.page}>
      <header>远程连接</header>
      {!(sysConfig?.remotes) ? <span>加载中...</span>
        : sysConfig?.remotes.map((remote, i) => (
          <div className={styles.box} key={i}>
            <ul className={styles.configs}>
              <li className={styles.configs_item}>
                <Status style={{ backgroundColor: remote.status }} />
              </li>
              <li className={styles.configs_item}>
                <div className={styles.configs_item_key}>
                  <span>名称</span>
                </div>
                <div className={styles.configs_item_values}>
                  <input
                    type="text"
                    value={remote.name}
                    onChange={e => setSysConfig(c => {
                      c.remotes[i].name = e.target.value;
                      return { ...c };
                    })}
                  />
                </div>
              </li>
              <li className={styles.configs_item}>
                <div className={styles.configs_item_key}>
                  <span>类型</span>
                </div>
                <div className={styles.configs_item_values}>
                  <span>{remote.type}</span>
                </div>
              </li>
              <li className={styles.configs_item}>
                <div className={styles.configs_item_key}>
                  <span>登录方式</span>
                </div>
                <div className={styles.configs_item_values}>
                  <span>{remote.loginType}</span>
                </div>
              </li>
              <li className={styles.configs_item}>
                <div className={styles.configs_item_key}>
                  <span>refreshToken</span>
                </div>
                <div className={styles.configs_item_values}>
                  <input
                    type="text"
                    value={remote.refreshToken}
                    onChange={e => setSysConfig(c => {
                      c.remotes[i] = { ...c.remotes[i], refreshToken: e.target.value };
                      return { ...c };
                    })}
                  />
                </div>
              </li>
            </ul>
          </div>
        ))}
    </div>
  );
};
