import { AllInbox, Close, Download, Info, PrivacyTip } from "@mui/icons-material";
import { Box, Dialog, DialogContent, IconButton, Stack } from "@mui/material";
import { getSysConfig } from "hox/sysConfig";
import useWindows, { closeWindow } from "hox/windows";
import humanize from "humanize";
import moment from "moment";
import { useState } from "react";
import Attribute from "../../components/FileViewer/Attribute";
import Metadata from "../../components/FileViewer/Metadata";
import SameFiles from "../../components/FileViewer/SameFiles";

export default function ({ id, props }) {
    let { metadata, hash, attribute } = props;
    console.log("ImageViewer", id, props);
    const [windows, setWindows] = useWindows();
    const [openMetadata, setOpenMetadata] = useState(false);
    const [openSameFiles, setOpenSameFiles] = useState(false);
    const [openAttribute, setOpenAttribute] = useState(false);
    const [sameFiles, setSameFiles] = useState([]);
    const sysConfig = getSysConfig().sysConfig;
    let { hash: hash2, exif, fileType } = metadata;
    let time = moment(attribute?.modifyTime / 1000 / 1000).format("YYYY年MM月DD日 HH:mm:ss");
    return (
        <>
            <Dialog open={true} fullScreen={true} onClose={() => closeWindow(setWindows, id)}>
                <Stack direction="row" justifyContent="space-between" sx={{
                    color: theme => theme.context.secondary,
                    backgroundColor: theme => theme.background.secondary,
                }}
                >
                    <Box sx={{ height: "28px", lineHeight: "28px", paddingLeft: "1em" }}>
                        {hash} - 图片查看器
                    </Box>
                    <Stack direction="row" justifyContent="flex-end" >
                        <IconButton title="下载"
                            href={`${sysConfig.webServer}/api/v1/download?hash=${hash}`}
                            download={attribute ? attribute.name : sameFiles.length > 0 ? sameFiles[0].name : undefined}
                            sx={{ height: "24px", width: "24px", color: theme => theme.context.secondary }}
                        >
                            <Download fontSize="small" sx={{ width: "20px", height: "20px" }} />
                        </IconButton>
                        <IconButton title="相同文件" onClick={() => setOpenSameFiles(true)}
                            sx={{ height: "24px", width: "24px", color: theme => theme.context.secondary }}
                        >
                            <AllInbox fontSize="small" sx={{ width: "20px", height: "20px" }} />
                        </IconButton>
                        {attribute && <IconButton title="文件属性" onClick={() => setOpenAttribute(true)}
                            sx={{ height: "24px", width: "24px", color: theme => theme.context.secondary }}
                        >
                            <Info fontSize="small" sx={{ width: "20px", height: "20px" }} />
                        </IconButton>}
                        <IconButton title="文件元数据" onClick={() => setOpenMetadata(true)}
                            sx={{ height: "24px", width: "24px", color: theme => theme.context.secondary }}
                        >
                            <PrivacyTip fontSize="small" sx={{ width: "20px", height: "20px" }} />
                        </IconButton>
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
                    <img style={{ maxWidth: "100%", maxHeight: "100%" }} src={`${sysConfig.webServer}/api/v1/image?hash=${hash}`} />
                </DialogContent>
                <Box sx={{
                    flex: "0 0 auto", padding: "8px",
                    color: theme => theme.context.secondary,
                    backgroundColor: theme => theme.background.secondary,
                }}>
                    <Stack direction="row" justifyContent="space-between">
                        <Box >
                            {humanize.filesize(attribute?.size)}
                        </Box>
                        <Box >
                            {time}
                        </Box>
                    </Stack>
                </Box>
            </Dialog>
            <SameFiles open={openSameFiles} setOpen={setOpenSameFiles} hash={hash} sameFiles={sameFiles} setSameFiles={setSameFiles} />
            <Metadata open={openMetadata} setOpen={setOpenMetadata} metadata={metadata} />
            {attribute && <Attribute open={openAttribute} setOpen={setOpenAttribute} attribute={attribute} />}
        </>
    );
}
