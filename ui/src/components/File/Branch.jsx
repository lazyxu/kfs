import './index.scss';
import {useClick} from "use";
import {open} from "api/api";
import {modeIsDir} from "api/utils/api";
import useResourceManager from 'hox/resourceManager';
import useContextMenu from "../../hox/contextMenu";
import SvgIcon from "../Icon/SvgIcon";
import {Box} from "@mui/material";

export default ({branch}) => {
    const [resourceManager, setResourceManager] = useResourceManager();
    return (
        <Box component="span" sx={{":hover": {backgroundColor: (theme) => theme.palette.action.hover}}}
             className='file-normal'>
            <div onMouseDown={useClick(null, () => {
                open(setResourceManager, branch.name, []);
            })}>
                <SvgIcon icon="wangpan" className='file-icon file-icon-file' fontSize="inherit"/>
            </div>
            <div className='file-name-wrapper'>
                <p className='file-name'>{branch.name}</p>
            </div>
        </Box>
    )
};
