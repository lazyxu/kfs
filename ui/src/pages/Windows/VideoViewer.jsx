import { AllInbox, Download, Info, PrivacyTip } from "@mui/icons-material";
import { Box, ButtonGroup, IconButton, Stack } from "@mui/material";
import { getMetadata } from "api/exif";
import { download, getDriverFile } from "api/fs";
import { getSysConfig } from "hox/sysConfig";
import humanize from "humanize";
import moment from "moment";
import FileAttribute from "pages/Files/DriverFiles/FileAttribute";
import { useEffect, useState } from "react";
import Metadata from "../../components/FileViewer/Metadata";
import SameFiles from "../../components/FileViewer/SameFiles";
import { StatusBar, TitleBar, Window, WorkingArea } from "./Window";

export default function ({ id, props }) {
    let { driver, filePath } = props;
    console.log("ImageViewer", id, props);
    const [driverFile, setDriverFile] = useState();
    const [metadata, setMetadata] = useState(false);
    const [openMetadata, setOpenMetadata] = useState(false);
    const [openSameFiles, setOpenSameFiles] = useState(false);
    const [openAttribute, setOpenAttribute] = useState(false);
    const [sameFiles, setSameFiles] = useState([]);
    const sysConfig = getSysConfig().sysConfig;
    useEffect(() => {
        getDriverFile(driver.id, filePath).then(df => {
            setDriverFile(df);
            getMetadata(df.hash).then(setMetadata);
        });
    }, []);
    return (
        <Window id={id}>
            <TitleBar id={id} title={filePath[filePath.length - 1] + " - 视频查看器"} buttons={driverFile && <ButtonGroup variant="contained">
                <IconButton title="下载" onClick={() => { download(driver.id, filePath) }}
                    sx={{ color: theme => theme.context.secondary }}
                >
                    <Download fontSize="small" />
                </IconButton>
                <IconButton title="相同文件" onClick={() => setOpenSameFiles(true)}
                    sx={{ color: theme => theme.context.secondary }}
                >
                    <AllInbox fontSize="small" />
                </IconButton>
                <IconButton title="文件属性" onClick={() => setOpenAttribute(true)}
                    sx={{ color: theme => theme.context.secondary }}
                >
                    <Info fontSize="small" />
                </IconButton>
                <IconButton title="文件元数据" onClick={() => setOpenMetadata(true)}
                    sx={{ color: theme => theme.context.secondary }}
                >
                    <PrivacyTip fontSize="small" />
                </IconButton>
            </ButtonGroup>} />
            <WorkingArea>
                <Box sx={{
                    padding: "0",
                    width: "100%",
                    height: "100%",
                    display: "flex",
                    justifyContent: "center",
                    alignItems: "center",
                    color: theme => theme.context.primary,
                    backgroundColor: theme => theme.background.primary,
                }}>
                    {driverFile && <video controls style={{ maxWidth: "100%", maxHeight: "100%" }} data-setup='{}'>
                        <source src={`${sysConfig.webServer}/api/v1/image?hash=${driverFile.hash}`} />
                        您的浏览器不支持 HTML5 video 标签。
                    </video>}
                </Box>
            </WorkingArea>
            <StatusBar>
                <Stack direction="row" justifyContent="space-between">
                    {driverFile ? <>
                        <Box >
                            {humanize.filesize(driverFile.size)}
                        </Box>
                        <Box >
                            {moment(driverFile.modifyTime / 1000 / 1000).format("YYYY年MM月DD日 HH:mm:ss")}
                        </Box>
                    </> :
                        <Box >
                            ...
                        </Box>}
                </Stack>
            </StatusBar>
            {driverFile && <>
                <SameFiles open={openSameFiles} setOpen={setOpenSameFiles} hash={driverFile.hash} sameFiles={sameFiles} setSameFiles={setSameFiles} />
                <Metadata open={openMetadata} setOpen={setOpenMetadata} metadata={metadata} disabled={!metadata} />
                {openAttribute && <FileAttribute fileAttribute={{ driver, filePath, driverFile }} onClose={setOpenAttribute} />}
            </>}
        </Window>
    );
}
