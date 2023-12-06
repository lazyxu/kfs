import { AllInbox, Close, Download, Info, PrivacyTip } from "@mui/icons-material";
import { Box, Dialog, DialogContent, IconButton, Stack } from "@mui/material";
import { getMetadata } from "api/exif";
import { download, getDriverFile } from "api/fs";
import { getSysConfig } from "hox/sysConfig";
import useWindows, { closeWindow } from "hox/windows";
import humanize from "humanize";
import moment from "moment";
import FileAttribute from "pages/Files/DriverFiles/FileAttribute";
import { useEffect, useState } from "react";
import Metadata from "../../components/FileViewer/Metadata";
import SameFiles from "../../components/FileViewer/SameFiles";

export default function ({ id, props }) {
    let { driver, filePath } = props;
    // let { metadata, hash, attribute } = props;
    console.log("ImageViewer", id, props);
    const [windows, setWindows] = useWindows();
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
        <>
            <Dialog open={true} fullScreen={true} onClose={() => closeWindow(setWindows, id)}>
                <Stack direction="row" justifyContent="space-between" sx={{
                    color: theme => theme.context.secondary,
                    backgroundColor: theme => theme.background.secondary,
                }}
                >
                    <Box sx={{ height: "28px", lineHeight: "28px", paddingLeft: "1em" }}>
                        {filePath[filePath.length - 1]} - 图片查看器
                    </Box>
                    <Stack direction="row" justifyContent="flex-end" >
                        {driverFile && <>
                            <IconButton title="下载" onClick={() => { download(driver.id, filePath) }}
                                sx={{ height: "24px", width: "24px", color: theme => theme.context.secondary }}
                            >
                                <Download fontSize="small" sx={{ width: "20px", height: "20px" }} />
                            </IconButton>
                            <IconButton title="相同文件" onClick={() => setOpenSameFiles(true)}
                                sx={{ height: "24px", width: "24px", color: theme => theme.context.secondary }}
                            >
                                <AllInbox fontSize="small" sx={{ width: "20px", height: "20px" }} />
                            </IconButton>
                            <IconButton title="文件属性" onClick={() => setOpenAttribute(true)}
                                sx={{ height: "24px", width: "24px", color: theme => theme.context.secondary }}
                            >
                                <Info fontSize="small" sx={{ width: "20px", height: "20px" }} />
                            </IconButton>
                            <IconButton title="文件元数据" onClick={() => setOpenMetadata(true)}
                                sx={{ height: "24px", width: "24px", color: theme => theme.context.secondary }}
                            >
                                <PrivacyTip fontSize="small" sx={{ width: "20px", height: "20px" }} />
                            </IconButton>
                        </>
                        }
                        <IconButton aria-label="close" onClick={() => closeWindow(setWindows, id)}
                            sx={{
                                padding: "4px 12px", borderRadius: '0',
                                color: theme => theme.context.secondary,
                                ":hover": {
                                    backgroundColor: "red",
                                }
                            }}
                        >
                            <Close sx={{ width: "20px", height: "20px" }} />
                        </IconButton>
                    </Stack>
                </Stack>
                <DialogContent sx={{
                    padding: "0",
                    width: "100%",
                    height: "100%",
                    display: "flex",
                    justifyContent: "center",
                    alignItems: "center",
                    color: theme => theme.context.primary,
                    backgroundColor: theme => theme.background.primary,
                }}>
                    {driverFile && <img style={{ maxWidth: "100%", maxHeight: "100%" }} src={`${sysConfig.webServer}/api/v1/image?hash=${driverFile.hash}`} />}
                </DialogContent>
                <Box sx={{
                    flex: "0 0 auto", padding: "8px",
                    color: theme => theme.context.secondary,
                    backgroundColor: theme => theme.background.secondary,
                }}>
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
                </Box>
            </Dialog>
            {driverFile && <>
                <SameFiles open={openSameFiles} setOpen={setOpenSameFiles} hash={driverFile.hash} sameFiles={sameFiles} setSameFiles={setSameFiles} />
                <Metadata open={openMetadata} setOpen={setOpenMetadata} metadata={metadata} disabled={!metadata} />
                {openAttribute && <FileAttribute fileAttribute={{ driver, filePath, driverFile }} onClose={setOpenAttribute} />}
            </>}
        </>
    );
}
