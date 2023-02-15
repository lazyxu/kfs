import useMenu from 'hox/menu';

import styles from './index.module.scss';
import SvgIcon from "../Icon/SvgIcon";

export default function ({items}) {
    const {menu, setMenu} = useMenu();
    return (
        <ul className={styles.nav_menu}>
            {items.map(item =>
                <li key={item.name} className={`${styles.nav_menu_item} ${menu === item.name ? styles.is_active : ''}`}
                    onClick={() => setMenu(item.name)}>
          <span className={styles.nav_menu_item_icon}>
            <SvgIcon icon={item.icon}/>
          </span>
                    <span>{item.name}</span>
                </li>)}
        </ul>
    );
}
