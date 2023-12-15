import useResourceManager, { openDir, openDrivers } from '@kfs/common/hox/resourceManager';
import SvgIcon from "@kfs/mui/components/Icon/SvgIcon";
import { Breadcrumbs, Link, Stack } from "@mui/material";

function PathElement({ name, icon, dirPath }) {
    const [resourceManager, setResourceManager] = useResourceManager();
    const { driver } = resourceManager;
    return (
        <Link color="inherit" underline="hover" onClick={(() => {
            if (dirPath) {
                openDir(setResourceManager, driver, dirPath);
            } else {
                openDrivers(setResourceManager);
            }
        })}>
            <Stack direction="row"
                justifyContent="flex-start"
                alignItems="center"
            >
                {icon && <SvgIcon icon={icon} />}
                {name}
            </Stack>
        </Link>
    )
}

export default () => {
    const [resourceManager, setResourceManager] = useResourceManager();
    const dirPath = resourceManager.dirPath || [];
    return (
        <Stack sx={{ padding: "0 1em", margin: "0.3em" }}
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
