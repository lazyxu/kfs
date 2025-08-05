import { getMetadata } from "@kfs/common/api/webServer/exif";
import { getSysConfig } from "@kfs/common/hox/sysConfig";
import { downloadByHash } from "@kfs/mui/api/fs";
import { TitleBar, Window, WorkingArea } from "@kfs/mui/components/Window/Window";
import { Download, Info, PrivacyTip } from "@mui/icons-material";
import { Box, ButtonGroup, IconButton } from "@mui/material";
import { useEffect, useState } from "react";
import Metadata from "./Metadata";
import SameFiles from "./SameFiles";
import FileAttributeWithPath from "../Files/DriverFiles/FileAttributeWithPath";

export default function ({ wid, props }) {
    let {
        driver, filePath, driverFile,
        metadata,
    } = props;
    console.log("VideoViewer", wid, props);
    const [downloadName, setDownloadName] = useState();
    const [metadataL, setMetadataL] = useState(metadata);
    const [openMetadata, setOpenMetadata] = useState(false);
    const [sameFiles, setSameFiles] = useState([]);
    const [openSameFiles, setOpenSameFiles] = useState(false);
    const [openAttribute, setOpenAttribute] = useState(false);
    const sysConfig = getSysConfig();
    useEffect(() => {
        if (!metadata) {
            getMetadata(wid).then(setMetadataL);
        }
        // if (driverFile) {
        //     listDriverFileByHash(hash).then(setSameFiles);
        //     setDownloadName(filePath[filePath.length - 1]);
        // } else {
        //     listDriverFileByHash(hash).then(sf => {
        //         setSameFiles(sf);
        //         // TODO: select a best name.
        //         setDownloadName(sf[0].name);
        //     });
        // }
    }, []);
    return (
        <Window id={wid}>
            <TitleBar wid={wid} title={downloadName + " - 视频查看器"} buttons={<ButtonGroup variant="contained">
                <IconButton title="下载" disabled={!downloadName} onClick={() => { downloadByHash(wid, downloadName) }}
                    sx={{ color: theme => theme.context.secondary }}
                >
                    <Download fontSize="small" />
                </IconButton>
                {/* <IconButton title="相同文件" onClick={() => setOpenSameFiles(true)}
                    sx={{ color: theme => theme.context.secondary }}
                >
                    <Badge badgeContent={sameFiles.length} color="secondary">
                        <AllInbox fontSize="small" />
                    </Badge>
                </IconButton> */}
                {driverFile && <IconButton title="文件属性" onClick={() => setOpenAttribute(true)}
                    sx={{ color: theme => theme.context.secondary }}
                >
                    <Info fontSize="small" />
                </IconButton>}
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
                    {metadataL && <video controls style={{ maxWidth: "100%", maxHeight: "100%" }} data-setup='{}'>
                        <source src={`${sysConfig.webServer}/api/v1/image?id=${metadataL.id}&size=${metadataL.size}&md5=${metadataL.md5}&sha256=${metadataL.sha256}`} />
                        您的浏览器不支持 HTML5 video 标签。
                    </video>}
                </Box>
            </WorkingArea>
            {openSameFiles && <SameFiles sameFiles={sameFiles} onClose={setOpenSameFiles} />}
            {openMetadata && <Metadata metadata={metadataL} onClose={setOpenMetadata}/>}
            {driverFile && openAttribute && <FileAttributeWithPath fileAttribute={{ driver, filePath, driverFile }} onClose={setOpenAttribute} />}
        </Window>
    );
}
