import {useEffect} from 'react';

import {list} from "../../api/api";
import File from "../../components/File";
import styles from './index.module.scss';
import AbsolutePath from "../../components/AbsolutePath";
import useResourceManager from 'hox/resourceManager';
import useSysConfig from 'hox/sysConfig';

function App() {
    const [resourceManager, setResourceManager] = useResourceManager();
    const {sysConfig} = useSysConfig();
    useEffect(() => {
        console.log("mount");
        (async () => {
            let {filePath, branchName} = resourceManager;
            await list(sysConfig, setResourceManager, branchName, filePath);
        })()
    }, []);
    console.log(resourceManager)

    return (
        <>
            <AbsolutePath/>
            {resourceManager.content === null ?
                <div className={styles.filesGridview}>
                    {resourceManager.dirItems.map((dirItem, i) => (
                        <File type={dirItem.Mode > 2147483648 ? 'dir' : 'file'} name={dirItem.Name} key={dirItem.Name}/>
                    ))}
                </div> :
                <div className={styles.filesGridview}>
                    { (new TextDecoder("utf-8")).decode((resourceManager.content))}
                </div>
            }
        </>
    );
}

export default App;
