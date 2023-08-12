import {list} from "api/fs";
import useResourceManager from 'hox/resourceManager';
import SvgIcon from "components/Icon/SvgIcon";
import {Box, Button, Card, CardActions, CardContent, Link, Stack, Typography} from "@mui/material";
import useContextMenu from "hox/contextMenu";
import humanize from 'humanize';
import {deleteDriver} from "../../api/driver";
import DeleteIcon from '@mui/icons-material/Delete';
import DriveFileRenameOutlineIcon from '@mui/icons-material/DriveFileRenameOutline';

export default ({driversElm, driver}) => {
    const [resourceManager, setResourceManager] = useResourceManager();
    const [contextMenu, setContextMenu] = useContextMenu();
    return (
        <Card sx={{width: 275}} variant="outlined">
            <CardContent>
                <Link underline="hover" onClick={() => list(setResourceManager, driver.name, [])}>
                    <Stack
                        direction="row"
                        justifyContent="space-between"
                        alignItems="center"
                        spacing={2}
                    >
                        <SvgIcon icon="wangpan" fontSize="inherit"/>
                        <Box sx={{flex: 1}}>{driver.name}</Box>
                    </Stack>
                </Link>
                <Box variant="body">
                    {/* <Box>文件总数：{driver.count}</Box> */}
                    {/* <Box>总大小：{humanize.filesize(driver.size)}</Box> */}
                    {/* <Typography>可修改该云盘的设备：any</Typography> */}
                    {/* <Typography>可读取该云盘的设备：any</Typography> */}
                </Box>
                <Box color="text.secondary">
                    描述：{driver.description}
                </Box>
            </CardContent>
            <CardActions>
                <Button size="small" color="error" startIcon={<DeleteIcon/>} variant="outlined"
                        onClick={() => deleteDriver(setResourceManager, driver.name)}>删除</Button>
                <Button size="small" startIcon={<DriveFileRenameOutlineIcon/>} variant="outlined" disabled>重命名</Button>
                <Button size="small" variant="outlined" disabled>重置</Button>
            </CardActions>
        </Card>
    )
};
