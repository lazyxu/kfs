import useResourceManager from 'hox/resourceManager';
import './index.scss';
import PathElement from "./PathElement";
import Icon from "../Icon/Icon";
import {Breadcrumbs, Link, Stack} from "@mui/material";
import {open} from "../../api/api";
import useSysConfig from "../../hox/sysConfig";

export default () => {
    const [resourceManager, setResourceManager] = useResourceManager();
    const {sysConfig} = useSysConfig();
    return (
        <Stack className='filePath'
               direction="row"
               justifyContent="flex-start"
               alignItems="center"
               spacing={1}
        >
            <Link onClick={() => {
                open(sysConfig, setResourceManager, resourceManager.branchName, []);
            }} size="small"><Icon icon="wangpan"/></Link>
            <div style={{color: "rgba(255, 255, 255, 0.7)"}}> ></div>
            <Breadcrumbs separator="/" maxItems={5}>
                <PathElement name={resourceManager.branchName} filePath={[]} icon="git"/>
                {resourceManager.filePath.map((elemName, i) =>
                    <PathElement key={i} name={elemName} filePath={resourceManager.filePath.slice(0, i + 1)}/>
                )}
            </Breadcrumbs>
        </Stack>
    )
};
