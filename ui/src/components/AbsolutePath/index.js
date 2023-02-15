import useResourceManager from 'hox/resourceManager';
import './index.scss';
import PathElement from "./PathElement";
import {Breadcrumbs, Link, Stack, Typography} from "@mui/material";
import {open} from "../../api/api";
import SvgIcon from "../Icon/SvgIcon";

export default () => {
    const [resourceManager, setResourceManager] = useResourceManager();
    return (
        <Stack className='filePath'
               direction="row"
               justifyContent="flex-start"
               alignItems="center"
               spacing={1}
        >
            <Link onClick={() => {
                open(setResourceManager, resourceManager.branchName, []);
            }} size="small"><SvgIcon icon="wangpan"/></Link>
            <Typography> ></Typography>
            <Breadcrumbs separator="/" maxItems={5}>
                <PathElement name={resourceManager.branchName} filePath={[]} icon="git"/>
                {resourceManager.filePath.map((elemName, i) =>
                    <PathElement key={i} name={elemName} filePath={resourceManager.filePath.slice(0, i + 1)}/>
                )}
            </Breadcrumbs>
        </Stack>
    )
};
