import SvgIcon from "../Icon/SvgIcon";
import styles from './index.module.scss';

export default () => (
    <div className={styles.logo}>
        <SvgIcon icon="wangpan1" className={styles.icon} />
        <span className={styles.name}>考拉云盘</span>
    </div>
);
