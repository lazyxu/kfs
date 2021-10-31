import { listBranches } from 'remote/branch';
import { useState, useEffect } from 'react';
import Icon from 'common/components/Icon/Icon';
import { onContextMenu } from 'remote/menu';

import styles from './index.module.scss';

export default () => {
  const [branches, setBranches] = useState([]);
  useEffect(() => {
    listBranches().then(setBranches).catch(console.log);
  }, []);
  return (
    <div>
      <header>云盘</header>
      <ul className={styles.pan}>
        {branches.map(branch => (
          <li
            key={branch.branchName}
            onContextMenu={e => {
              onContextMenu(e);
            }}
          >
            <div><Icon icon="wangpan" className={styles.icon} /></div>
            <span className={styles.text}>{branch.branchName}</span>
          </li>
        ))}
      </ul>
    </div>
  );
};
