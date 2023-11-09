import ContentCopyIcon from '@mui/icons-material/ContentCopy';
import DeleteIcon from '@mui/icons-material/Delete';
import FileDownloadIcon from '@mui/icons-material/FileDownload';
import ModeEditIcon from '@mui/icons-material/ModeEdit';
import { Box, Stack, Tooltip } from "@mui/material";
import IconButton from "@mui/material/IconButton";
import humanize from 'humanize';
import moment from 'moment';
import { download } from "../../../api/fs";
import useResourceManager from "../../../hox/resourceManager";
import TextFileViewer from "./TextFileViewer";
import styles from './index.module.scss';

export default ({file}) => {
    const [resourceManager, setResourceManager] = useResourceManager();
    console.log("FileViewer", typeof file.content, file);
    let time = moment(file.modifyTime / 1000 / 1000).format("YYYY年MM月DD日 HH:mm:ss");
    return (
        <>
            <Stack
                className={styles.fileHeaderViewer}
                direction="row"
                justifyContent="space-between"
                alignItems="center"
                spacing={0.5}
            >
                {humanize.filesize(file.size)} | {time}
                <Stack
                    direction="row"
                    justifyContent="flex-end"
                    alignItems="flex-end"
                    spacing={0.5}
                >
                    <Tooltip title="编辑">
                    <span><IconButton color="inherit" disabled={true}>
                    <ModeEditIcon fontSize="small"/>
                    </IconButton></span>
                    </Tooltip>
                    <Tooltip title="复制文本内容">
                    <span><IconButton onClick={() => {
                        navigator.clipboard.writeText(file.content);
                    }} disabled={file.tooLarge}>
                    <ContentCopyIcon fontSize="small"/>
                    </IconButton></span>
                    </Tooltip>
                    <Tooltip title="下载">
                    <span><IconButton onClick={() => {
                        download(resourceManager.driverId, resourceManager.filePath)
                    }}>
                    <FileDownloadIcon fontSize="small"/>
                    </IconButton></span>
                    </Tooltip>
                    <Tooltip title="删除">
                    <span><IconButton disabled={true}>
                    <DeleteIcon fontSize="small"/>
                    </IconButton></span>
                    </Tooltip>
                </Stack>
            </Stack>
            <Box style={{flex: "auto"}} className={styles.fileViewer}>
                <TextFileViewer file={file}/>
            </Box>
        </>
    )
};
