import { Link, Stack } from "@mui/material";
import { listDriver } from "api/driver";
import { openDir, openFile } from "api/fs";
import useResourceManager from 'hox/resourceManager';
import SvgIcon from "../../Icon/SvgIcon";

export default ({type, name, icon, filePath}) => {
    const [resourceManager, setResourceManager] = useResourceManager();
    return (
        <Link color="inherit" underline="hover" onClick={(() => {
            if (filePath) {
                if (type === "file") {
                    openFile(setResourceManager, resourceManager.driverId, filePath, resourceManager.file);
                } else {
                    openDir(setResourceManager, resourceManager.driver, filePath);
                }
            } else {
                listDriver(setResourceManager);
            }
        })}>
            <Stack direction="row"
                   justifyContent="flex-start"
                   alignItems="center"
            >
                {icon && <SvgIcon icon={icon}/>}
                {name}
            </Stack>
        </Link>
    )
};
