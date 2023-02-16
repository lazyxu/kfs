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
            <Breadcrumbs separator=">" maxItems={5}>
                <PathElement name="我的文件" icon="wangpan"/>
                <PathElement name={resourceManager.branchName} filePath={[]}/>
                {resourceManager.filePath.map((elemName, i) =>
                    <PathElement key={i} name={elemName} filePath={resourceManager.filePath.slice(0, i + 1)}/>
                )}
            </Breadcrumbs>
        </Stack>
    )
};
