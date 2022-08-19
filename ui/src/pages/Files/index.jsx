import { useEffect } from 'react';

import { list } from "../../rpc/ws";
import File from "../../components/File";
import styles from './index.module.scss';
import FilePath from "../../components/FilePath";
import useResourceManager from 'hox/resourceManager';
import useSysConfig from 'hox/sysConfig';

function App() {
    const [resourceManager, setResourceManager] = useResourceManager();
    const { sysConfig } = useSysConfig();
    useEffect(() => {
        (async () => {
            let dirItems;
            let { filePath, branchName } = resourceManager;
            await list(sysConfig, branchName, filePath, (total) => {
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
