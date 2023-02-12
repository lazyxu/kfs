import Icon from "components/Icon/Icon";
import {useClick} from "use";
import {open} from "api/api";
import {modeIsDir} from "api/utils/api";
import useResourceManager from 'hox/resourceManager';
import useSysConfig from 'hox/sysConfig';
import useContextMenu from "hox/contextMenu";
import styles from "./index.module.scss";

export default ({file}) => {
    // const [resourceManager, setResourceManager] = useResourceManager();
    // const {sysConfig} = useSysConfig();
    // const [contextMenu, setContextMenu] = useContextMenu();
    // const onClick = e => {
    //     console.log('onClick')
    // }
    // const onDoubleClick = e => {
    //     console.log('onDoubleClick')
    // }
    // let {filePath, branchName} = resourceManager;
    // const {Name, Mode} = dirItem;
    // filePath = filePath.concat(Name);
    console.log("FileViewer", file);
    return (
        <div className={styles.fileViewer}>
            {(new TextDecoder("utf-8")).decode((file.Content))}
        </div>
    )
};
