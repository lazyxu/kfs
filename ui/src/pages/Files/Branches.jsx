import {useEffect, useRef} from 'react';
import styles from './index.module.scss';
import AbsolutePath from "components/AbsolutePath";
import useResourceManager from 'hox/resourceManager';
import useContextMenu from "hox/contextMenu";
import useDialog from "hox/dialog";
import {listBranch} from 'api/branch';
import Branch from "../../components/File/Branch";

export default function () {
    const [resourceManager, setResourceManager] = useResourceManager();
    const [contextMenu, setContextMenu] = useContextMenu();
    const [dialog, setDialog] = useDialog();
    const filesElm = useRef(null);
    useEffect(() => {
        listBranch(setResourceManager);
    }, []);
    console.log("branches.resourceManager", resourceManager, "contextMenu", contextMenu, "dialog", dialog)

    return (
        <>
            <AbsolutePath/>
            <div ref={filesElm} className={styles.filesGridview}>
                {resourceManager.branches.map((branch, i) => (
                    <Branch key={branch.name} branch={branch}>{branch.name}</Branch>
                ))}
            </div>
        </>
    );
}
