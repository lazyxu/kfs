import styles from './index.module.scss';

export function Body({ children }) {
  return (
    <div className={styles.body}>
      {children}
    </div>
  );
}

export function Layout({ children }) {
  return (
    <div className={styles.layout}>
      {children}
    </div>
  );
}

export function Sider({ children }) {
  return (
    <div className={styles.sider}>
      {children}
    </div>
  );
}

export function Content({ children }) {
  return (
    <div className={styles.content}>
      {children}
    </div>
  );
}
