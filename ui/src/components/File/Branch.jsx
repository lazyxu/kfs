import './index.scss';
import {useClick} from "use";
import {list} from "api/fs";
import useResourceManager from 'hox/resourceManager';
import SvgIcon from "../Icon/SvgIcon";
import {Stack} from "@mui/material";

export default ({branch}) => {
    const [resourceManager, setResourceManager] = useResourceManager();
    return (
        <Stack component="span" sx={{":hover": {backgroundColor: (theme) => theme.palette.action.hover}}}
               className='file-normal'
               justifyContent="flex-start"
               alignItems="center"
               spacing={1}>
            <div onMouseDown={useClick(null, () => {
                list(setResourceManager, branch.name, []);
            })}>
                <SvgIcon icon="wangpan" className='file-icon file-icon-file' fontSize="inherit"/>
            </div>
            <div className='file-name-wrapper'>
                <p className='file-name'>{branch.name}</p>
            </div>
        </Stack>
    )
};
