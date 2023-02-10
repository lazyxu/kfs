import {useEffect} from 'react';

import {list} from "../../api/api";
import File from "../../components/File";
import styles from './index.module.scss';
import AbsolutePath from "../../components/AbsolutePath";
import useResourceManager from 'hox/resourceManager';
import useSysConfig from 'hox/sysConfig';
import DefaultContextMenu from "../../components/ContextMenu/DefaultContextMenu";
import useContextMenu from "../../hox/contextMenu";

function App() {
    const [resourceManager, setResourceManager] = useResourceManager();
    const {sysConfig} = useSysConfig();
    const [contextMenu, setContextMenu] = useContextMenu();
    useEffect(() => {
        console.log("mount");
        (async () => {
            let {filePath, branchName} = resourceManager;
            await list(sysConfig, setResourceManager, branchName, filePath);
        })()
    }, []);
    console.log(resourceManager)

    return (
        <div className={styles.right}>
            <AbsolutePath/>
            {resourceManager.content === null ?
                <div className={styles.filesGridview} onContextMenu={(e) => {
                    e.preventDefault();
                    console.log(e.target, e.currentTarget, e.target === e.currentTarget);
                    if (e.target === e.currentTarget) {
                        console.log("default");
                        const {clientX, clientY} = e;
                        let {x, y, width, height} = e.currentTarget.getBoundingClientRect();
                        setContextMenu({
                            type: 'default',
                            clientX, clientY,
                            x, y, width, height,
                        })
                    }
                    //     const { fileListView } = this.context.state;
                    //     if (e.target === fileListView || e.target.getAttribute('data-tag') !== 'choose-able') {
                    //         const { clientX, clientY } = e;
                    //         const { x, y } = fileListView.getBoundingClientRect();
                    //         this.context.setState({
                    //             contextMenuForFile: null,
                    //             contextMenu: {
                    //                 x: Math.min(clientX, x + fileListView.clientWidth - 200),
                    //                 y: Math.min(clientY, y + fileListView.clientHeight - 120),
                    //             },
                    //         });
                    //     }
                }}>
                    {resourceManager.dirItems.map((dirItem, i) => (
                        <File type={dirItem.Mode > 2147483648 ? 'dir' : 'file'} name={dirItem.Name}
                              key={dirItem.Name}/>
                    ))}
                </div> :
                <div className={styles.filesGridview}>
                    {(new TextDecoder("utf-8")).decode((resourceManager.content))}
                </div>
            }
            <DefaultContextMenu/>
        </div>
    );
}

export default App;
