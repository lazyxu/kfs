import {useClick} from "use";
import {list} from "api/fs";
import useResourceManager from 'hox/resourceManager';
import SvgIcon from "components/Icon/SvgIcon";
import {Box} from "@mui/material";
import useContextMenu from "hox/contextMenu";
import humanize from 'humanize';

export default ({branchesElm, branch}) => {
    const [resourceManager, setResourceManager] = useResourceManager();
    const [contextMenu, setContextMenu] = useContextMenu();
    return (
        <Box sx={{
            ":hover": {backgroundColor: (theme) => theme.palette.action.hover},
            display: 'grid',
            gridAutoFlow: 'row',
            gridTemplateColumns: '4em 1fr',
            gridTemplateRows: '1em 1em 1em 1em',
            columnGap: 1,
            width: "10em",
        }}
             onContextMenu={(e) => {
                 e.preventDefault();
                 e.stopPropagation();
                 const {clientX, clientY} = e;
                 let {x, y, width, height} = branchesElm.current.getBoundingClientRect();
                 setContextMenu({
                     type: 'branch',
                     branch,
                     clientX, clientY,
                     x, y, width, height,
                 })
             }}>
            <Box sx={{gridColumn: '1', gridRow: '1 / 4'}} onMouseDown={useClick(null, () => {
                list(setResourceManager, branch.name, []);
            })}>
                <SvgIcon icon="wangpan" style={{height: "4em", width: "4em"}} fontSize="inherit"/>
            </Box>
            <Box sx={{gridColumn: '2', gridRow: '1', height: "1em"}}>{branch.name}</Box>
            <Box sx={{gridColumn: '2', gridRow: '2'}}>{branch.description}</Box>
            <Box sx={{gridColumn: '2', gridRow: '3'}}>{branch.count}个文件</Box>
            <Box sx={{gridColumn: '2', gridRow: '4'}}>{humanize.filesize(branch.size)}</Box>
        </Box>
    )
};
