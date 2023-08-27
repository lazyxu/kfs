import SvgIcon from "../Icon/SvgIcon";
import { modeIsDir, isDCIM } from "api/utils/api";
import { Box, Stack } from "@mui/material";
import { getMetadata } from "api/web/exif";
import ImageViewer from "components/FileViewer/ImageViewer";
import { useState } from "react";
import VideoViewer from "components/FileViewer/VideoViewer";

export default function ({ dirItem }) {
    const [open, setOpen] = useState(false);
    const [metadata, setMetadata] = useState();
    return (
        <Box className="file-icon-box">
            {modeIsDir(dirItem.mode) ?
                <SvgIcon icon="folder1" className='file-icon file-icon-folder' fontSize="inherit" /> :
                isDCIM(dirItem.name) ?
                    <img src={`http://127.0.0.1:1123/thumbnail?size=64&hash=${dirItem.hash}`} loading="lazy" onClick={() => getMetadata(dirItem.hash).then(m => {
                        setMetadata(m);
                        setOpen(true);
                    })} />
                    :
                    <SvgIcon icon="file12" className='file-icon file-icon-file' fontSize="inherit" />
            }
            {metadata?.fileType?.type == "video" && <VideoViewer open={open} setOpen={setOpen} metadata={metadata} hash={metadata.hash} />}
            {metadata?.fileType?.type == "image" && <ImageViewer open={open} setOpen={setOpen} metadata={metadata} hash={metadata.hash} />}
        </Box>
    );
}
