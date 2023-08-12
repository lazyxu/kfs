import './index.scss';
import {useClick} from "use";
import {list, openFile} from "api/fs";
import {modeIsDir} from "api/utils/api";
import useResourceManager from 'hox/resourceManager';
import useContextMenu from "../../hox/contextMenu";
import SvgIcon from "../Icon/SvgIcon";
import {Box, Stack} from "@mui/material";

function isDCIM(name) {
    if (name.endsWith(".jpg")) {
        return true;
    }
    return false;
}

export default function ({dirItem, filesElm}) {
    const [resourceManager, setResourceManager] = useResourceManager();
    const [contextMenu, setContextMenu] = useContextMenu();
    let {filePath, driverName} = resourceManager;
    const {name, mode} = dirItem;
    filePath = filePath.concat(name);
    return (
        <Stack component="span" sx={{":hover": {backgroundColor: (theme) => theme.palette.action.hover}}}
               className='file-normal'
               justifyContent="flex-start"
               alignItems="center"
               spacing={1}
               onContextMenu={(e) => {
                   e.preventDefault();
                   e.stopPropagation();
                   const {clientX, clientY} = e;
                   let {x, y, width, height} = filesElm.current.getBoundingClientRect();
                   setContextMenu({
                       type: 'file',
                       dirItem,
                       clientX, clientY,
                       x, y, width, height,
                   })
               }}>
            <Box onMouseDown={useClick(null, () => {
                if (modeIsDir(mode)) {
                    list(setResourceManager, driverName, filePath);
                } else {
                    openFile(setResourceManager, driverName, filePath, dirItem);
                }
            })}>
                {modeIsDir(mode) ?
                    <SvgIcon icon="folder1" className='file-icon file-icon-folder' fontSize="inherit"/> :
                    isDCIM(name) ? <img src={"http://127.0.0.1:1123/thumbnail?hash="+dirItem.hash} loading="lazy"/> :
                    <SvgIcon icon="file12" className='file-icon file-icon-file' fontSize="inherit"/>
                }
                
            </Box>
            <Box kfs-attr="file" style={{width: "100%", overflowWrap: "break-word", textAlign: "center"}}>{name}</Box>
        </Stack>
    )
};
