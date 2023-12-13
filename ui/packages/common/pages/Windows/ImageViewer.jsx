import { getMetadata } from "@kfs/common/api/exif";
import { downloadByHash, listDriverFileByHash } from "@kfs/common/api/fs";
import { getSysConfig } from "@kfs/common/hox/sysConfig";
import FileAttribute from "@kfs/common/pages/Files/DriverFiles/FileAttribute";
import { AllInbox, Download, Info, PrivacyTip } from "@mui/icons-material";
import { Badge, Box, ButtonGroup, IconButton } from "@mui/material";
import { useEffect, useState } from "react";
import Metadata from "./Metadata";
import SameFiles from "./SameFiles";
import { TitleBar, Window, WorkingArea } from "./Window";

export default function ({ id, props }) {
    let {
        driver, filePath, driverFile,
        hash,
    } = props;
    console.log("ImageViewer", id, props);
    if (driverFile) {
        hash = driverFile.hash;
    }
    const [downloadName, setDownloadName] = useState();
    const [metadata, setMetadata] = useState();
    const [openMetadata, setOpenMetadata] = useState(false);
    const [sameFiles, setSameFiles] = useState([]);
    const [openSameFiles, setOpenSameFiles] = useState(false);
    const [openAttribute, setOpenAttribute] = useState(false);
    const sysConfig = getSysConfig();
    useEffect(() => {
        getMetadata(hash).then(setMetadata);
        if (driverFile) {
            listDriverFileByHash(hash).then(setSameFiles);
            setDownloadName(filePath[filePath.length - 1]);
        } else {
            listDriverFileByHash(hash).then(sf => {
                setSameFiles(sf);
                // TODO: select a best name.
                setDownloadName(sf[0].name);
            });
        }
    }, []);
    return (
        <Window id={id}>
            <TitleBar id={id} title={downloadName + " - 图片查看器"} buttons={<ButtonGroup variant="contained">
                <IconButton title="下载" disabled={!downloadName} onClick={() => { downloadByHash(hash, downloadName) }}
                    sx={{ color: theme => theme.context.secondary }}
                >
                    <Download fontSize="small" />
                </IconButton>
                <IconButton title="相同文件" onClick={() => setOpenSameFiles(true)}
                    sx={{ color: theme => theme.context.secondary }}
                >
                    <Badge badgeContent={sameFiles.length} color="secondary">
                        <AllInbox fontSize="small" />
                    </Badge>
                </IconButton>
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
                    {hash && <img style={{ maxWidth: "100%", maxHeight: "100%" }} src={`${sysConfig.webServer}/api/v1/image?hash=${hash}`} />}
                </Box>
            </WorkingArea>
            {openSameFiles && <SameFiles hash={hash} sameFiles={sameFiles} onClose={setOpenSameFiles} />}
            {openMetadata && <Metadata hash={hash} metadata={metadata} onClose={setOpenMetadata}/>}
            {driverFile && openAttribute && <FileAttribute fileAttribute={{ driver, filePath, driverFile }} onClose={setOpenAttribute} />}
        </Window>
    );
}
