import { Link, Stack } from "@mui/material";
import useResourceManager, { openDir, openDrivers } from 'hox/resourceManager';
import SvgIcon from "../../Icon/SvgIcon";

export default ({ name, icon, dirPath }) => {
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
};
