import TextFileViewer from "./TextFileViewer";
import styles from './index.module.scss';
import moment from 'moment';
import humanize from 'humanize';

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
    console.log("FileViewer", typeof file.Content, file);
    let time = moment(file.ModifyTime / 1000 / 1000).format("YYYY-MM-DD HH:mm:ss");
    return (
        <>
            <div className={styles.fileHeaderViewer}>
                {humanize.filesize(file.Size)} | {time}
            </div>
            <div className={styles.fileViewer}>
                <TextFileViewer file={file}/>
            </div>
        </>
    )
};
