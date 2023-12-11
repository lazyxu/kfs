import { AllInbox, ContentCopy, Info, Save } from '@mui/icons-material';
import { default as FileDownload } from '@mui/icons-material/FileDownload';
import { Badge, Box, Link, Stack } from "@mui/material";
import IconButton from "@mui/material/IconButton";
import { download, listDriverFileByHash, openFile } from "api/fs";
import { getSysConfig } from "hox/sysConfig";
import humanize from 'humanize';
import moment from 'moment';
import FileAttribute from 'pages/Files/DriverFiles/FileAttribute';
import { useEffect, useState } from 'react';
import SameFiles from "./SameFiles";
import { StatusBar, TitleBar, Window, WorkingArea } from './Window';

export default ({ id, props }) => {
    const { driver, filePath, driverFile } = props;
    console.log("TextViewer", id, props);
    const { hash } = driverFile;
    const [loaded, setLoaded] = useState();
    const [openAttribute, setOpenAttribute] = useState(false);
    const [sameFiles, setSameFiles] = useState([]);
    const [openSameFiles, setOpenSameFiles] = useState(false);
    useEffect(() => {
        listDriverFileByHash(hash).then(setSameFiles);
        openFile(driver.id, filePath).then(setLoaded);
    }, []);
    return (
        <Window id={id}>
            <TitleBar id={id} title={filePath[filePath.length - 1] + " - 文本编辑器"} buttons={<>
                <IconButton title="保存" disabled
                    sx={{ color: theme => theme.context.secondary }}
                >
                    <Save fontSize="small" />
                </IconButton>
                <IconButton title="复制文本内容" onClick={() => {
                    navigator.clipboard.writeText(loaded.content);
                }} disabled={!loaded || loaded.tooLarge}
                    sx={{ color: theme => theme.context.secondary }}
                >
                    <ContentCopy fontSize="small" />
                </IconButton>
                <IconButton title="下载" onClick={() => { download(driver.id, filePath) }}
                    sx={{ color: theme => theme.context.secondary }}
                >
                    <FileDownload fontSize="small" />
                </IconButton>
                <IconButton title="相同文件" onClick={() => setOpenSameFiles(true)}
                    sx={{ color: theme => theme.context.secondary }}
                >
                    <Badge badgeContent={sameFiles.length} color="secondary">
                        <AllInbox fontSize="small" />
                    </Badge>
                </IconButton>
                <IconButton title="文件属性" onClick={() => setOpenAttribute(true)}
                    sx={{ color: theme => theme.context.secondary }}
                >
                    <Info />
                </IconButton>
            </>} />
            <WorkingArea>
                {!loaded ?
                    <>加载中...</> :
                    loaded.tooLarge ?
                        <>
                            文件大于{humanize.filesize(getSysConfig().sysConfig.maxContentSize)}，不支持在线查看，你可以选择
                            <Link underline="hover" onClick={() => {
                                download(driver.id, filePath)
                            }}>下载该文件</Link>。
                        </> :
                        <p style={{ wordBreak: "break-all", whiteSpace: "break-spaces", outline: "none" }} contentEditable>
                            {loaded.content}
                        </p>
                }
            </WorkingArea>
            <StatusBar>
                <Stack direction="row" justifyContent="space-between">
                    <Box >
                        {humanize.filesize(driverFile.size)}
                    </Box>
                    <Box >
                        {moment(driverFile.modifyTime / 1000 / 1000).format("YYYY年MM月DD日 HH:mm:ss")}
                    </Box>
                </Stack>
            </StatusBar>
            {openSameFiles && <SameFiles hash={hash} sameFiles={sameFiles} onClose={setOpenSameFiles} />}
            {openAttribute && <FileAttribute fileAttribute={{ driver, filePath, driverFile }} onClose={setOpenAttribute} />}
        </Window>
    )
};
