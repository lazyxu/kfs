import { useState, useEffect } from 'react';
import useFilePath from 'hox/filepath';

import styles from './index.module.scss';
import Branch from './Branch';
import File from './File';

export default () => {
  const [filepath, setFilepath] = useFilePath();
  useEffect(() => setFilepath(''), []);
  return (
    <div className={styles.pan}>
      {filepath === '' ? <Branch /> : <File />}
    </div>
  );
};
