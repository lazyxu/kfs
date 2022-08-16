import { useEffect } from 'react';

import { list } from "../../rpc/ws";
import File from "../../components/File";
import styles from './index.module.scss';
import FilePath from "../../components/FilePath";
import useResourceManager from 'hox/resourceManager';

function App() {
    const [resourceManager, setResourceManager] = useResourceManager();
    useEffect(() => {
        (async () => {
            let dirItems;
            let { filePath, branchName } = resourceManager;
            await list(branchName, filePath, (total) => {
                dirItems = new Array(total);
            }, (dirItem, i) => {
                dirItems[i] = dirItem;
            });
            setResourceManager(prev => {
                return { ...prev, branchName, filePath, dirItems };
            });
        })()
    }, []);
    console.log(resourceManager.dirItems)

    return (
        <>
            <FilePath />
            <div className={styles.filesGridview}>
                {resourceManager.dirItems.map((dirItem, i) => (
                    <File type={dirItem.Mode > 2147483648 ? 'dir' : 'file'} name={dirItem.Name} key={dirItem.Name} />
                ))}
            </div>
        </>
    );
}

export default App;
