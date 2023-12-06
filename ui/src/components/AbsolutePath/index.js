import { Breadcrumbs, Stack } from "@mui/material";
import useResourceManager from 'hox/resourceManager';
import PathElement from "./PathElement";
import './index.scss';

export default () => {
    const [resourceManager, setResourceManager] = useResourceManager();
    const dirPath = resourceManager.dirPath || [];
    return (
        <Stack className='filePath'
            direction="row"
            justifyContent="flex-start"
            alignItems="center"
            spacing={1}
        >
            <Breadcrumbs separator=">" maxItems={5}>
                <PathElement name="我的云盘" />
                <PathElement name={resourceManager.driver?.name} dirPath={[]} />
                {dirPath.map((name, i) =>
                    <PathElement key={i} name={name} dirPath={dirPath.slice(0, i + 1)} />
                )}
            </Breadcrumbs>
        </Stack>
    )
};
