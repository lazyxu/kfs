import { listDriverFileByHash } from "@kfs/common/api/webServer/driverfile";
import { getMetadata } from "@kfs/common/api/webServer/exif";
import { getSysConfig } from "@kfs/common/hox/sysConfig";
import { downloadByHash } from "@kfs/mui/api/fs";
import { TitleBar, Window, WorkingArea } from "@kfs/mui/components/Window/Window";
import FileAttribute from "@kfs/mui/pages/Files/DriverFiles/FileAttribute";
import { AllInbox, Download, Info, PrivacyTip } from "@mui/icons-material";
import { Badge, Box, ButtonGroup, IconButton } from "@mui/material";
import { useEffect, useState } from "react";
import Metadata from "./Metadata";
import SameFiles from "./SameFiles";

export default function ({ wid, props }) {
    let {
        filePath, driverFile,
        metadata,
    } = props;
    console.log("ImageViewer", wid, props);
    const [downloadName, setDownloadName] = useState();
    const [metadataL, setMetadataL] = useState(metadata);
    const [openMetadata, setOpenMetadata] = useState(false);
    const [sameFiles, setSameFiles] = useState([]);
    const [openSameFiles, setOpenSameFiles] = useState(false);
    const [openAttribute, setOpenAttribute] = useState(false);
    const sysConfig = getSysConfig();
    useEffect(() => {
        if (!metadata) {
            getMetadata(driverFile.file_id).then(setMetadataL);
        }
        // if (driverFile) {
        //     listDriverFileByHash(driverFile.id).then(setSameFiles);
        //     setDownloadName(filePath[filePath.length - 1]);
        // } else {
        //     listDriverFileByHash(driverFile.id).then(sf => {
        //         setSameFiles(sf);
        //         // TODO: select a best name.
        //         setDownloadName(sf[0].name);
        //     });
        // }
    }, []);
    return (
        <Window id={wid}>
            <TitleBar wid={wid} title={downloadName + " - 图片查看器"} buttons={<ButtonGroup variant="contained">
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
                }}>
                    {metadataL && <img style={{ maxWidth: "100%", maxHeight: "100%" }} src={`${sysConfig.webServer}/api/v1/image?id=${metadataL.id}&size=${metadataL.size}&md5=${metadataL.md5}&sha256=${metadataL.sha256}`} />}
                </Box>
            </WorkingArea>
            {openSameFiles && <SameFiles sameFiles={sameFiles} onClose={setOpenSameFiles} />}
            {openMetadata && <Metadata metadata={metadataL} onClose={setOpenMetadata}/>}
            {driverFile && openAttribute && <FileAttribute fileAttribute={{ filePath, driverFile }} onClose={setOpenAttribute} />}
        </Window>
    );
}
