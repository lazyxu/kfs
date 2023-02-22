import TextFileViewer from "./TextFileViewer";
import styles from './index.module.scss';
import moment from 'moment';
import humanize from 'humanize';
import ContentCopyIcon from '@mui/icons-material/ContentCopy';
import IconButton from "@mui/material/IconButton";
import {Box, Stack, Tooltip, useColorScheme} from "@mui/material";
import ModeEditIcon from '@mui/icons-material/ModeEdit';
import DeleteIcon from '@mui/icons-material/Delete';
import {download} from "../../../api/fs";
import useResourceManager from "../../../hox/resourceManager";
import FileDownloadIcon from '@mui/icons-material/FileDownload';

export default ({file}) => {
    const [resourceManager, setResourceManager] = useResourceManager();
    console.log("FileViewer", typeof file.content, file);
    let time = moment(file.modifyTime / 1000 / 1000).format("YYYY年MM月DD日 HH:mm:ss");
    const {mode} = useColorScheme();
    console.log("mode", mode)
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
                        download(resourceManager.branchName, resourceManager.filePath)
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
            <Box className={styles.fileViewer}>
                <TextFileViewer file={file}/>
            </Box>
        </>
    )
};
