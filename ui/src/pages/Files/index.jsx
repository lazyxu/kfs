import {useEffect, useState} from 'react';

import {list} from "../../rpc/ws";
import File from "../../components/File";
import styles from './index.module.scss';

function App() {
    const [dirItems, setDirItems] = useState([]);
    useEffect(() => {
        (async () =>{
            let newDirItems;
            await list((total) => {
                newDirItems = new Array(total);
            }, (dirItem, i) => {
                newDirItems[i] = dirItem;
            });
            setDirItems(newDirItems);
        })()
    }, []);

    return (
        <div className={styles.filesGridview}>
            {dirItems.map((dirItem, i) => (
                <File type='dir' name={dirItem.Name} key={dirItem.Name}/>
            ))}
        </div>
    );
}

export default App;
