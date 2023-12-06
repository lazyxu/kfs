import { Link, Stack } from "@mui/material";
import useResourceManager, { openDir, openDrivers } from 'hox/resourceManager';
import useWindows from "hox/windows";
import SvgIcon from "../../Icon/SvgIcon";

export default ({ type, name, icon, filePath }) => {
    const [resourceManager, setResourceManager] = useResourceManager();
    const { driver } = resourceManager;
    const [windows, setWindows] = useWindows();
    return (
        <Link color="inherit" underline="hover" onClick={(() => {
            if (filePath) {
                openDir(setResourceManager, driver, filePath);
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
};
