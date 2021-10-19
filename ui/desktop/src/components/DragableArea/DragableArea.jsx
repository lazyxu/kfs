import styles from './index.module.scss';

export default function ({ children }) {
  return (
    <div className={styles.dragable_area}>
      {children}
    </div>
  );
}
