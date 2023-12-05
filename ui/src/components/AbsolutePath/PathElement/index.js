import { Link, Stack } from "@mui/material";
import { listDriver } from "api/driver";
import { openDir } from "api/fs";
import useResourceManager from 'hox/resourceManager';
import useWindows, { APP_TEXT_VIEWER, newWindow } from "hox/windows";
import SvgIcon from "../../Icon/SvgIcon";

export default ({ type, name, icon, filePath }) => {
    const [resourceManager, setResourceManager] = useResourceManager();
    const { driver } = resourceManager;
    const [windows, setWindows] = useWindows();
    return (
        <Link color="inherit" underline="hover" onClick={(() => {
            if (filePath) {
                if (type === "file") {
                    newWindow(setWindows, APP_TEXT_VIEWER, { driver, filePath });
                } else {
                    openDir(setResourceManager, driver, filePath);
                }
            } else {
                listDriver(setResourceManager);
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
