import styles from './index.module.scss';
import SvgIcon from "../Icon/SvgIcon";

export default () => (
    <div className={styles.logo}>
        <SvgIcon icon="wangpan1" className={styles.icon}/>
        <span className={styles.name}>考拉云盘</span>
    </div>
);
