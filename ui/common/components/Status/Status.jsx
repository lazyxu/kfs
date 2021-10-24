import styles from './index.module.scss';

export default function ({ style }) {
  return (
    <div
      className={styles.status}
      style={style}
    />
  );
}
