import {useRef} from 'react';
import File from "components/File";
import styles from './index.module.scss';
import AbsolutePath from "components/AbsolutePath";
import useResourceManager from 'hox/resourceManager';
import DefaultContextMenu from "components/ContextMenu/DefaultContextMenu";
import useContextMenu from "hox/contextMenu";
import FileContextMenu from "components/ContextMenu/FileContextMenu";
import FileViewer from "./FileViewer/FileViewer";
import Dialog from "components/Dialog";

export default function () {
    const [resourceManager, setResourceManager] = useResourceManager();
    const [contextMenu, setContextMenu] = useContextMenu();
    const filesElm = useRef(null);

    return (
        <>
            <AbsolutePath/>
            {resourceManager.file ?
                <FileViewer file={resourceManager.file}/> :
                <div ref={filesElm} className={styles.filesGridview} onContextMenu={(e) => {
                    e.preventDefault();
                    // console.log(e.target, e.currentTarget, e.target === e.currentTarget);
                    // if (e.target === e.currentTarget) {
                    const {clientX, clientY} = e;
                    let {x, y, width, height} = e.currentTarget.getBoundingClientRect();
                    setContextMenu({
                        type: 'default',
                        clientX, clientY,
                        x, y, width, height,
                    })
                    // }
                }}>
                    {resourceManager.dirItems.map((dirItem, i) => (
                        <File filesElm={filesElm} dirItem={dirItem} key={dirItem.name}/>
                    ))}
                </div>
            }
            <DefaultContextMenu/>
            <FileContextMenu/>
            <Dialog/>
        </>
    );
}
