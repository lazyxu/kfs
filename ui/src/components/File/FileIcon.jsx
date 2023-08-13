import SvgIcon from "../Icon/SvgIcon";
import { modeIsDir, isDCIM } from "api/utils/api";
import { Box, Stack } from "@mui/material";

export default function ({ dirItem }) {
    return (
        <Box className="file-icon-box">
            {modeIsDir(dirItem.mode) ?
                <SvgIcon icon="folder1" className='file-icon file-icon-folder' fontSize="inherit" /> :
                isDCIM(dirItem.name) ?
                    <img src={"http://127.0.0.1:1123/thumbnail?hash=" + dirItem.hash} loading="lazy" />
                    :
                    <SvgIcon icon="file12" className='file-icon file-icon-file' fontSize="inherit" />
            }
        </Box>
    );
}
