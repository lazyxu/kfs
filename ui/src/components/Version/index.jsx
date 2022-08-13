import styles from './index.module.scss';

export default () => (
  <div className={styles.version}>
    {process.env.REACT_APP_PLATFORM}.{process.env.NODE_ENV}
  </div>
);
