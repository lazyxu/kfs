import useResourceManager from 'hox/resourceManager';
import './index.scss';
import PathElement from "./PathElement";
import {Breadcrumbs, Stack} from "@mui/material";

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
                <PathElement type="driver" name="我的云盘"/>
                <PathElement name={resourceManager.driverName} filePath={[]}/>
                {resourceManager.filePath.map((elemName, i) =>
                    <PathElement
                        type={i === resourceManager.filePath.length-1 && !resourceManager.dirItems ? "file" : "dir"}
                        key={i} name={elemName}
                        filePath={resourceManager.filePath.slice(0, i + 1)}/>
                )}
            </Breadcrumbs>
        </Stack>
    )
};
