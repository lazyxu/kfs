import {list} from "api/fs";
import useResourceManager from 'hox/resourceManager';
import SvgIcon from "components/Icon/SvgIcon";
import {Button, Card, CardActions, CardContent, Link, Stack, Typography} from "@mui/material";
import useContextMenu from "hox/contextMenu";
import humanize from 'humanize';
import {deleteBranch} from "../../api/branch";
import DeleteIcon from '@mui/icons-material/Delete';
import DriveFileRenameOutlineIcon from '@mui/icons-material/DriveFileRenameOutline';

export default ({branchesElm, branch}) => {
    const [resourceManager, setResourceManager] = useResourceManager();
    const [contextMenu, setContextMenu] = useContextMenu();
    return (
        <Card sx={{width: 275}}>
            <CardContent>
                <Link underline="hover" onClick={() => list(setResourceManager, branch.name, [])}>
                    <Stack
                        direction="row"
                        justifyContent="space-between"
                        alignItems="center"
                        spacing={2}
                    >
                        <SvgIcon icon="wangpan" fontSize="inherit"/>
                        <Typography sx={{flex: 1}}>{branch.name}</Typography>
                    </Stack>
                </Link>
                <Typography variant="body2">
                    <Typography>文件总数：{branch.count}</Typography>
                    <Typography>总大小：{humanize.filesize(branch.size)}</Typography>
                    {/* <Typography>可修改该云盘的设备：any</Typography> */}
                    {/* <Typography>可读取该云盘的设备：any</Typography> */}
                </Typography>
                <Typography color="text.secondary">
                    描述：{branch.description}
                </Typography>
            </CardContent>
            <CardActions>
                <Button size="small" color="error" startIcon={<DeleteIcon/>} variant="outlined"
                        onClick={() => deleteBranch(setResourceManager, branch.name)}>删除</Button>
                <Button size="small" startIcon={<DriveFileRenameOutlineIcon/>} variant="outlined">重命名</Button>
                <Button size="small" variant="outlined">重置</Button>
            </CardActions>
        </Card>
    )
};
